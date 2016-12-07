package extractor

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

type Bounds struct {
	X1 float64
	Y1 float64
	X2 float64
	Y2 float64
}

type SectionConfig struct {
	Name          string
	Page          int
	Height        float64
	Width         float64
	SectionBounds Bounds `json:"bounds"`
	Adjustment    string
}

type ExtractorConfig struct {
	BasePageHeight float64         `json:"base_page_height"`
	BasePageWidth  float64         `json:"base_page_width"`
	Sections       []SectionConfig `json:"sections"`
	SectionsByName map[string]SectionConfig
}

func NewExtractorConfig(filename string) (*ExtractorConfig, error) {
	config := &ExtractorConfig{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Errorf("error reading config %s - %s", filename, err.Error())
		return nil, err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		glog.Errorf("error loading config %s - %s", filename, err.Error())
	}
	config.SectionsByName = make(map[string]SectionConfig, len(config.Sections))
	for _, section := range config.Sections {
		config.SectionsByName[section.Name] = section
	}
	glog.Infof("loaded config - %#v", config)
	return config, nil
}

type Extractor interface {
	ExtractSection(string, string, float64, float64) ([]byte, error)
	GetOffset(string) (float64, float64, error)
}
