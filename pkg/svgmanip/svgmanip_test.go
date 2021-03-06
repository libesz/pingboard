package svgmanip

import (
	"strings"
	"testing"

	"github.com/beevik/etree"
)

var goodConfig = []Target{{ID: "path10", Fill: "#00ff00"}}
var badConfig = []Target{{ID: "path11", Fill: "#00ff00"}}

func TestSingleChangeEmbedded(t *testing.T) {
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

	CheckAndChange(root, goodConfig[0])
	path := doc.SelectElement("svg").SelectElement("g").SelectElement("path")
	if !strings.Contains(path.SelectAttr("style").Value, "fill:#00ff00") {
		t.Errorf("Style mismatch %v", path.SelectAttr("style").Value)
	}
	if !strings.Contains(path.SelectAttr("style").Value, "bla=bla") {
		t.Errorf("Style mismatch2 %v", path.SelectAttr("style").Value)
	}
}

func TestSingleChange(t *testing.T) {
	var testXML = `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg>
	  <g>
		<path fill="#000000" id="path10" />
	  </g>
	</svg>
	`
	doc := etree.NewDocument()
	doc.ReadFromString(testXML)
	root := doc.SelectElement("svg")

	CheckAndChange(root, goodConfig[0])
	path := doc.SelectElement("svg").SelectElement("g").SelectElement("path")
	if !strings.Contains(path.SelectAttr("fill").Value, "#00ff00") {
		t.Errorf("Fill mismatch %v", path.SelectAttr("fill").Value)
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

	err := CheckAndChange(root, goodConfig[0])
	if err != nil {
		t.Errorf("Style should be added if missing")
	}
	err = CheckAndChange(root, badConfig[0])
	if err == nil {
		t.Errorf("Path should be missing")
	}
}

func TestDocUpdateEmbedded(t *testing.T) {
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
		t.Errorf("CheckDoc should fail due to missing path")
	}
	err = CheckDoc(doc, goodConfig)
	if err != nil {
		t.Errorf("CheckDoc should pass")
	}
}
