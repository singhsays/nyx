package extractor

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

type PageAttrs struct {
	Height float64 `json:"height"`
}

// TabulaExtractor extracts tables from pdf files using tabula-java CLI.
type TabulaExtractor struct {
	config      *ExtractorConfig
	extractArgs []string
	detectArgs  []string
}

// NewTabulaExtractor returns an initialized instance of TabuleExtractor.
func NewTabulaExtractor(config *ExtractorConfig, filename, javaPath, tabulaPath string) *TabulaExtractor {
	return &TabulaExtractor{
		config:      config,
		extractArgs: []string{javaPath, "-Djava.awt.headless=true", "-jar", tabulaPath, "--no-spreadsheet", "-i"},
		detectArgs:  []string{"-fJSON", "--guess", "-d"},
	}
}

// getBounds calculates the bounds for a given section name.
func (e *TabulaExtractor) getBounds(section string, offset float64) (Bounds, error) {
	bounds := Bounds{}
	config, ok := e.config.SectionsByName[section]
	if !ok {
		return bounds, fmt.Errorf("Section %s not found.", section)
	}
	bounds.X1 = config.SectionBounds.X1
	bounds.Y1 = config.SectionBounds.Y1
	bounds.X2 = config.SectionBounds.X2
	bounds.Y2 = config.SectionBounds.Y2
	// NOTE: Assume that width is fixed and offset applies to height only.
	if strings.Contains(config.Adjustment, "t") {
		bounds.Y1 = bounds.Y1 - offset
	}
	if strings.Contains(config.Adjustment, "b") {
		bounds.Y2 = bounds.Y2 - offset
	}
	return bounds, nil
}

// GetOffset calculates a page's offset from basePageHeight.
func (e *TabulaExtractor) GetOffset(filename string) (float64, error) {
	args := append(e.extractArgs, e.detectArgs...)
	args = append(args, filename)
	glog.Infof("Executing %v", args)
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		glog.Errorf("error detecting page offset - %s", err.Error())
		return 0.0, err
	}
	glog.Infof("Detection Output - %s", output)
	pageConfig := []PageAttrs{PageAttrs{Height: 0.0}}
	err = json.Unmarshal(output, &pageConfig)
	if err != nil {
		glog.Errorf("error detecting page offset - %s", err.Error())
		return 0.0, err
	}
	return e.config.BasePageHeight - pageConfig[0].Height, nil
}

// ExtractSection extracts a single section from the given file.
func (e *TabulaExtractor) ExtractSection(section string, filename string, offset float64) ([]byte, error) {
	glog.Infof("Extracting section %s from %s at offset %f", section, filename, offset)
	bounds, err := e.getBounds(section, offset)
	if err != nil {
		return nil, err
	}
	boundsArgs := fmt.Sprintf("-a %f,%f,%f,%f", bounds.Y1, bounds.X1, bounds.Y2, bounds.X2)
	args := append(e.extractArgs, boundsArgs, filename)
	glog.Infof("Executing %v", args)
	return exec.Command(args[0], args[1:]...).Output()
}
