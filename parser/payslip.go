package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"nyx/extractor"

	"github.com/golang/glog"
)

type PayslipParser struct {
	sections []string
}

func NewPayslipParser() *PayslipParser {
	return &PayslipParser{
		sections: []string{"pay_period", "pay_summary", "earnings", "deductions", "taxes"},
	}
}

func (p *PayslipParser) parseSection(name string, data []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	fmt.Println(records)
	return records, nil
}

func (p *PayslipParser) Parse(e extractor.Extractor, filename string) {
	offset, err := e.GetOffset(filename)
	if err != nil {
		glog.Exitf("error getting page offset - %s", err.Error())
	}

	for _, section := range p.sections {
		fmt.Println(section)
		output, err := e.ExtractSection(section, filename, offset)
		if err != nil {
			glog.Errorf("error extracting section %s - %s", section, err.Error())
		}
		p.parseSection(section, output)
	}
}
