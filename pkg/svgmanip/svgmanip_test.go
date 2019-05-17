package svgmanip

import (
	"log"
	"strings"
	"testing"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"gopkg.in/yaml.v2"
)

func TestYaml(t *testing.T) {
	var data = `
svgpath: "a.svg"
targets:
- id: path10
  fill: "#00ff00"
`

	config := config.Config{}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		t.Errorf("a")
		log.Fatalf("error unmarshal: %v", err)
	}
	if config.Targets[0].Fill != "#00ff00" {
		t.Errorf("error %v", config)
	}
	if config.SvgPath != "a.svg" {
		t.Errorf("error %v", config)
	}
}

var goodConfig = config.Config{
	Targets: []config.Target{{SvgId: "path10", Fill: "#00ff00"}}}
var badConfig = config.Config{
	Targets: []config.Target{{SvgId: "path11", Fill: "#00ff00"}}}

func TestSingleChange(t *testing.T) {
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
	root := doc.SelectElement("svg")

	CheckAndChange(root, goodConfig.Targets[0])
	path := doc.SelectElement("svg").SelectElement("g").SelectElement("path")
	if !strings.Contains(path.SelectAttr("style").Value, "fill:#00ff00") {
		t.Errorf("Style mismatch %v", path.SelectAttr("style").Value)
	}
	if !strings.Contains(path.SelectAttr("style").Value, "bla=bla") {
		t.Errorf("Style mismatch2 %v", path.SelectAttr("style").Value)
	}
}

func TestChangeErrors(t *testing.T) {
	var testXML = `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg>
	  <g>
		<path id="path10" />
	  </g>
	</svg>
	`
	doc := etree.NewDocument()
	doc.ReadFromString(testXML)
	root := doc.SelectElement("svg")

	err := CheckAndChange(root, goodConfig.Targets[0])
	if err == nil {
		t.Errorf("Style should be missing")
	}
	err = CheckAndChange(root, badConfig.Targets[0])
	if err == nil {
		t.Errorf("Path should be missing")
	}
}

func TestDocUpdate(t *testing.T) {
	var testXML = `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg>
	  <g><g>
		<path style="fill:#000000, bla=bla" id="path10" />
	  </g></g>
	</svg>
	`
	doc := etree.NewDocument()
	doc.ReadFromString(testXML)
	err := UpdateDoc(doc, goodConfig)
	if err != nil {
		t.Errorf("UpdateDoc should pass here")
	}
	path := doc.SelectElement("svg").SelectElement("g").SelectElement("g").SelectElement("path")
	if !strings.Contains(path.SelectAttr("style").Value, "fill:#00ff00") {
		t.Errorf("Style mismatch %v", path.SelectAttr("style").Value)
	}
	if !strings.Contains(path.SelectAttr("style").Value, "bla=bla") {
		t.Errorf("Style mismatch2 %v", path.SelectAttr("style").Value)
	}
	err = UpdateDoc(doc, badConfig)
	if err == nil {
		t.Errorf("UpdateDoc should fail here")
	}
}

func TestDocCheck(t *testing.T) {
	var testXML = `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg>
	  <g><g>
		<path style="fill:#000000, bla=bla" id="path10" />
	  </g></g>
	</svg>
	`
	doc := etree.NewDocument()
	doc.ReadFromString(testXML)
	err := CheckDoc(doc, badConfig)
	if err == nil {
		t.Errorf("CheckDoc should fail")
	}
	err = CheckDoc(doc, goodConfig)
	if err != nil {
		t.Errorf("CheckDoc should pass")
	}
}
