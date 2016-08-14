package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"nyx/extractor"
	"nyx/models"
	"nyx/util"
	"time"

	"github.com/golang/glog"
)

// PayslipParser parsers a given file using an extractor, into a payslip.
// The parser is responsible for converting csv type data to the corresponding
// data model. The extractor does the actual file parsing and extraction of csv
// data.
type PayslipParser struct {
	filename  string
	extractor extractor.Extractor
	offset    float64
	payslip   *models.Payslip
	sections  []string
}

// NewPayslipParser returns a new initialized PayslipParser instance.
// The parser is specific to a given file, instantiate a new one for each file.
func NewPayslipParser(filename string, currency string, e extractor.Extractor) (*PayslipParser, error) {
	var err error
	parser := &PayslipParser{
		filename:  filename,
		extractor: e,
		payslip:   &models.Payslip{Currency: currency},
	}
	// Get the page offset based on the reference page's height.
	parser.offset, err = e.GetOffset(filename)
	if err != nil {
		return nil, fmt.Errorf("error getting page offset - %s", err.Error())
	}
	return parser, nil
}

// extractSection extracts a named section from the given filename.
// It also parses the named section to csv rows.
func (p *PayslipParser) extractSection(name string) ([][]string, error) {
	output, err := p.extractor.ExtractSection(name, p.filename, p.offset)
	if err != nil {
		return nil, fmt.Errorf("error extracting %s section - %s", name, err.Error())
	}
	return p.parseSection(name, output)
}

// parseSection parses a named section from the given data into csv formatted rows.
// It is generic and agnostic to the content of the rows.
func (p *PayslipParser) parseSection(name string, data []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing %s section - %s", name, err.Error())
	}
	return records, nil
}

// parsePeriod parses the pay period section.
func (p *PayslipParser) parsePeriod() error {
	sectionRows, err := p.extractSection("pay_period")
	if err != nil {
		return err
	}
	p.payslip.StartDate, err = time.Parse("01/02/2006", sectionRows[0][1])
	if err != nil {
		return err
	}
	p.payslip.EndDate, err = time.Parse("01/02/2006", sectionRows[1][1])
	if err != nil {
		return err
	}
	p.payslip.Date, err = time.Parse("01/02/2006", sectionRows[2][1])
	if err != nil {
		return err
	}
	p.payslip.DocumentID = sectionRows[3][1]
	p.payslip.NetPay, err = util.ToAmount(sectionRows[4][1])
	p.payslip.Date, err = time.Parse("01/02/2006", sectionRows[2][1])
	if err != nil {
		return err
	}
	return nil
}

// parseSummary parses the pay summary section.
func (p *PayslipParser) parseSummary() error {
	sectionRows, err := p.extractSection("pay_summary")
	if err != nil {
		return err
	}
	p.payslip.Gross, err = util.ToAmount(sectionRows[1][1])
	if err != nil {
		return err
	}
	p.payslip.Taxable, err = util.ToAmount(sectionRows[1][2])
	if err != nil {
		return err
	}
	p.payslip.Taxes, err = util.ToAmount(sectionRows[1][3])
	if err != nil {
		return err
	}
	p.payslip.Deductions, err = util.ToAmount(sectionRows[1][4])
	if err != nil {
		return err
	}
	p.payslip.Net, err = util.ToAmount(sectionRows[1][5])
	if err != nil {
		return err
	}
	return nil
}

// parseHead parses the current and ytd values for a given payslip head
// from the given csv data row.
func (p *PayslipParser) parseHead(row []string, currentIndex, ytdIndex int) (*models.PayslipHead, error) {
	var (
		head = &models.PayslipHead{Name: row[0]}
		err  error
	)
	if row[currentIndex] == "" && row[ytdIndex] == "" {
		return nil, fmt.Errorf("error parsing either current or ytd")
	}
	if head.Current, err = util.ToAmount(row[currentIndex]); err != nil {
		return nil, fmt.Errorf("error parsing current %s - %s", row[0], err)
	}
	if head.YTD, err = util.ToAmount(row[ytdIndex]); err != nil {
		return nil, fmt.Errorf("error parsing ytd %s - %s", row[0], err)
	}
	return head, nil
}

// parseEarnings parses the earnings section.
func (p *PayslipParser) parseEarnings() error {
	sectionRows, err := p.extractSection("earnings")
	if err != nil {
		return err
	}
	for _, row := range sectionRows[1 : len(sectionRows)-1] {
		head, err := p.parseHead(row, 3, 4)
		if err != nil {
			glog.Error(err)
			continue
		}
		p.payslip.IncomeHeads = append(p.payslip.IncomeHeads, head)
	}
	return nil
}

// parseDeductions parses the deductions section.
func (p *PayslipParser) parseDeductions() error {
	sectionRows, err := p.extractSection("deductions")
	if err != nil {
		return err
	}
	for _, row := range sectionRows[1:] {
		head, err := p.parseHead(row, 1, 2)
		if err != nil {
			glog.Error(err)
			continue
		}
		p.payslip.DeductionHeads = append(p.payslip.DeductionHeads, head)
	}
	return nil
}

// parseTaxes parses the taxes section.
func (p *PayslipParser) parseTaxes() error {
	sectionRows, err := p.extractSection("taxes")
	if err != nil {
		return err
	}
	for _, row := range sectionRows[1:] {
		head, err := p.parseHead(row, 1, 2)
		if err != nil {
			glog.Error(err)
			continue
		}
		p.payslip.TaxHeads = append(p.payslip.TaxHeads, head)
	}
	return nil
}

// Parse is the main entry point for a parser that parses
// the given file using the given extractor and returns
// a populated Payslip model.
func (p *PayslipParser) Parse() (*models.Payslip, error) {
	var err error
	if err = p.parsePeriod(); err != nil {
		return nil, err
	}
	if err = p.parseSummary(); err != nil {
		return nil, err
	}
	if err = p.parseEarnings(); err != nil {
		return nil, err
	}
	if err = p.parseDeductions(); err != nil {
		return nil, err
	}
	if err = p.parseTaxes(); err != nil {
		return nil, err
	}
	return p.payslip, nil
}
