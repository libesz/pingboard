package svgmanip

import (
	"testing"
	"github.com/beevik/etree"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

/*func TestHello(t *testing.T) {
	var testXML = `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg>
	  <g>
		<path style="fill:#000000" id="path10" />
	  </g>
	</svg>
	`
	doc := etree.NewDocument()
	doc.ReadFromString(testXML)
	path := doc.SelectElement("svg").SelectElement("g").SelectElement("path")

	CheckAndChange(path)
	if path.SelectAttr("style").Value != "fill:#00ff00" {
        t.Errorf("Style mismatch")
    }
}*/

func TestYaml(t *testing.T) {
	var data = `
changemap:
- id: path10
  fill: "#00ff00"
`

	configs := Configs{}
    
	err := yaml.Unmarshal([]byte(data), &configs)
	if err != nil {
		t.Errorf("a")
		log.Fatalf("error unmarshal: %v", err)
	}
	//t.Errorf("error %v", q)
	if configs.Cfgs[0].Fill != "#00ff00" {
		t.Errorf("error %v", configs)
	}
}

func TestNewInterface(t *testing.T) {
	config := Configs{[]Config{{ID:"path10", Fill:"#00ff00"}}}
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