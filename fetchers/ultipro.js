#!/usr/bin/env node

// Define command line arguments.
var argv = require('yargs')
  .default('loglevel', 'error')
  .default('path', '.')
  .default('config', '~/.cryptex.json')
  .default('show', false)
  .default('timeout', 50000)
  .default('query', 'latest')
  .default('useragent', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36')
  .default('zoom', '0.5')
  .argv;

// Load dependencies.
var cryptex = require('cryptex'),
    moment = require('moment'),
    path = require('path'),
    winston = require('winston'),
    untildify = require('untildify'),
    util = require('util');

// Secrets configuration.
var conf = require(untildify(argv.config))[process.env.NODE_ENV || 'default'],
    secrets = new cryptex.Cryptex({config: conf});

// Initialize download manager.
var Nightmare = require('nightmare');
require('nightmare-download-manager')(Nightmare);

var nightmare = Nightmare({
  show: argv.show,
  openDevTools: argv.show,
  waitTimeout: argv.timeout,
  gotoTimeout: argv.timeout,
  downloadTimeout: argv.timeout,
  downloadResponseWait: argv.timeout,
  webPreferences: {
    partition: 'nopersist'
  }
});

nightmare.on('download', function(state, downloadItem) {
  winston.info('download event:\n', state, downloadItem)
  if (state == 'started') {
    var match = downloadItem.url.match(/gennumber=(.*)!pagererid/);
    if (match.length >= 2) {
      filename = path.join(argv.path, match[1]+'-'+argv.zoom+'.pdf');
      nightmare.emit('download', filename, downloadItem);
    }
  }
});

/**
 * @typedef {{
 *    filename: (string|undefined),
 *    link: (string|undefined), 
 * }}
 */
var DownloadRef;

/*
 * parseDownloads compiles a list of files to download based on the query argument.
 * @param {!string} query - specifies how to filter the list of files to be downloaded.
 * @param {!Object} conf - is an associative array of configuration parameters.
 * @return {Array<DownloadRef>}
 */
var parseDownloads = function(query, zoom, conf) {
  // NOTE: This needs to be inlined since it runs in the browsers context.
  function parser(row) {
    match = row.href.match("'(pages/.*)', '(coid.*)', true")
    if (match.length == 3) {
      return [
        conf.base_url,        
        "/Customs/GOOG/",
        match[1], "?", conf.query_params, "!!", match[2],
        "!printtopdf=inline!PdfPageScaleToFitSinglePage=true!PdfOutputAreaWidth=8!PdfOutputAreaHeight=10.5!PdfZoomLevel=",
        zoom,
        "!PdfInvisibleElements=childDividerContainer;buttonbar"
      ].join('');
    }
    return null;
  };

  var rows = [];
  var nodelist = document.querySelectorAll('table tr td a[href*="EEPayrollPayCheckDetail.aspx"]');
  for(var i = nodelist.length; i--; rows.unshift(nodelist[i]));
  rows.reverse();
  var subset = rows;
  switch (query) {
    case 'all':
      // do nothing, subset is initialized to all rows.
      break;
    case (query.match(/last\d+/) || {}).input:
      count = parseInt(query.substr(4));
      // If parsed rows are more than desired count, recalculate the subset.
      // else use the default of all rows.
      if (rows.length > count)
        subset = rows.slice(-1 * count);
      break;
    case 'latest':
      // Fall-through to default.
    default:
      // Use the first row as subset.
      subset = rows.slice(-1);
  }
  return subset.map(parser);
};

/*
 * queueDownload queues an individual file for download.
 * The queue is processed sequencially since we need to nagivate to different
 * urls for each download.
 * @param {DownloadRef} download
 * @return {Promise}
 */
var queueDownload = function(download) {
  winston.info(download)
  return new Promise(function(resolve) {
    return nightmare.goto(download)
    .waitDownloadsComplete()
    .then(function(info) {
      winston.info('finished downloading - ', download);
      resolve(info);
    }, function(info) {
      winston.error('error downloading - ', info);
      // resolve anyways so that the next download continues.
      resolve(info);
    })
    .catch(function(error) {
      winston.error('error: ', error);
      // resolve anyways so that the next download continues.
      resolve(error);
    });
  });
};

/**
 * startSession initiates the browser session.
 * @param {!string} username
 * @param {!string} password
 * @param {!string} args - command line arguments.
 */
var startSession = function(username, password, args, conf) {
  nightmare
    .downloadManager()
    .useragent(args.useragent)
    .viewport(1440, 990)
    .goto(conf.oauth_url)
    .insert('input#Email', username)
    .click('input#next')
    .wait('input#Passwd')
    .insert('input#Passwd', password)
    .click('input#signIn')
    .wait('iframe#ContentFrame')
    .goto(conf.base_url + '/pages/VIEW/EePayrollPayCheckHistory.aspx?' + conf.query_params)
    .wait(2000)
    .evaluate(parseDownloads, args.query, args.zoom, conf)
    .then(function(files) {
      winston.info('parsed', files.length, 'files - ', files);
      return files.reduce(function(sequence, next) {
        winston.info('queued download - ', next);
        return sequence.then(function() {
          return queueDownload(next);
        });
      }, Promise.resolve('done'));
    })
    .then(function(done) {
      // Explicitly log out to clear server side session.
      return nightmare.goto(conf.base_url + "/logout.aspx")
        .wait('input[name="ctl00$Content$btnReturnLogin"]')
    })
    .then(function(){
      nightmare.end();
      // kill the Electron process explicitly to ensure no orphan child processes
      nightmare.proc.disconnect();
      nightmare.proc.kill();
      nightmare.ended = true;
      nightmare = null;
    })
    .then(function() {
      winston.info('done.');
    })
    .catch(function (error) {
      if (error.code != -3) {
        winston.error('download failed:', error);
      }
    });
};

// Main
winston.level = argv.loglevel;
secrets.getSecrets(['username', 'password']).then(function(values) {
  startSession(values.username, values.password, argv, conf);
});