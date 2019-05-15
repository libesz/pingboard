package svgmanip

import (
	"testing"
	"github.com/beevik/etree"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

func TestYaml(t *testing.T) {
	var data = `
targets:
- id: path10
  fill: "#00ff00"
`

	config := Config{}
    
	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		t.Errorf("a")
		log.Fatalf("error unmarshal: %v", err)
	}
	if config.Targets[0].Fill != "#00ff00" {
		t.Errorf("error %v", config)
	}
}

func TestSingleChange(t *testing.T) {
	config := Config{[]Target{{SvgId:"path10", Fill:"#00ff00"}}}
	var testXML = `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg>
	  <g>
		<path style="fill:#000000, bla=bla" id="path10" />
	  </g>
	</svg>
	`
	doc := etree.NewDocument()
	doc.ReadFromString(testXML)
	path := doc.SelectElement("svg").SelectElement("g").SelectElement("path")

	CheckAndChange(path, config)
	if !strings.Contains(path.SelectAttr("style").Value, "fill:#00ff00") {
        t.Errorf("Style mismatch %v", path.SelectAttr("style").Value)
  }
}