package main

import (
	"flag"
	"nyx/extractor"
	"nyx/importer"
	"nyx/parser"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
)

var (
	configPath = flag.String("config_file", "./config/payslip.conf.json", "Path to the config file.")
	javaPath   = flag.String("java_path", "/usr/bin/java", "Path to the java binary.")
	tabulaPath = flag.String("tabula_path", "/Users/sumeets/bin/tabula-0.9.0.jar", "Path to the Tabula jar.")
	// MongoImporter
	mongoAddress    = flag.String("mongo_address", "localhost:27017", "mongodb server address.")
	mongoDatabase   = flag.String("mongo_db", "nyx", "mongodb database name.")
	mongoCollection = flag.String("mongo_collection", "payslips", "mongodb collection name.")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		glog.Exitf("missing required argument filename - %v", args)
	}

	config, err := extractor.NewExtractorConfig(*configPath)
	if err != nil {
		glog.Exit(err)
	}

	e := extractor.NewTabulaExtractor(config, args[0], *javaPath, *tabulaPath)

	i, err := importer.NewMongoImporter(*mongoAddress, *mongoDatabase, *mongoCollection)
	if err != nil {
		glog.Exit(err)
	}
	defer i.Close()

	for _, filename := range args {
		p, err := parser.NewPayslipParser(filename, e)
		if err != nil {
			glog.Error(err)
			continue
		}
		payslip, err := p.Parse()
		if err != nil {
			glog.Error(err)
			continue
		}
		i.Import(payslip, true, bson.M{"document_id": payslip.DocumentID})
	}
}
