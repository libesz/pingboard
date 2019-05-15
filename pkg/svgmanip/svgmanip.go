package svgmanip

import (
	"fmt"
	"regexp"

	"github.com/beevik/etree"
)

type Target struct {
	SvgId    string `yaml:"id"`
	Fill     string `yaml:"fill"`
	Method   string `yaml:"method,omitempty"`
	EndPoint string `yaml:"endpoint,omitempty"`
}

type Config struct {
	Targets []Target `yaml:"targets"`
}

func CheckAndChange(elem *etree.Element, config Config) {
	for _, v := range config.Targets {
		if elem.SelectAttr("id").Value == v.SvgId {
			var re = regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
			elem.SelectAttr("style").Value = re.ReplaceAllString(elem.SelectAttr("style").Value, `fill:`+v.Fill)
			fmt.Printf("  ATTR: %s=%s\n", "style", elem.SelectAttr("style").Value)
		}
	}
}

func UpdateDoc(doc *etree.Document, config Config) {
	root := doc.SelectElement("svg")
	fmt.Println("ROOT element:", root.Tag)
	for _, g := range root.SelectElements("g") {
		fmt.Println(" CHILD element:", g.Tag)
		for _, path := range g.SelectElements("path") {
			fmt.Println(" CHILD element:", path.Tag)
			CheckAndChange(path, config)
		}
		for _, ellipse := range g.SelectElements("ellipse") {
			fmt.Println(" CHILD element:", ellipse.Tag)
			CheckAndChange(ellipse, config)
		}
		for _, rect := range g.SelectElements("rect") {
			fmt.Println(" CHILD element:", rect.Tag)
			CheckAndChange(rect, config)
		}
	}
}
