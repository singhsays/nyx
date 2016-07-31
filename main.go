package main

import (
	"flag"
	"nyx/extractor"
	"nyx/parser"

	"github.com/golang/glog"
)

var (
	configPath = flag.String("config_file", "./config/payslip.conf.json", "Path to the config file.")
	javaPath   = flag.String("java_path", "/usr/bin/java", "Path to the java binary.")
	tabulaPath = flag.String("tabula_path", "/Users/sumeets/bin/tabula-0.9.0.jar", "Path to the Tabula jar.")
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
	p, err := parser.NewPayslipParser(args[0], e)
	if err != nil {
		glog.Exit(err)
	}

	err = p.Parse()
}
