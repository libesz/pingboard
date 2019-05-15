package svgmanip

import (
	"fmt"
	"regexp"

	"github.com/beevik/etree"
)

type Config struct {
	ID   string `yaml:"id"`
	Fill string `yaml:"fill"`
}

type Configs struct {
	Cfgs []Config `yaml:"changemap"`
}

func CheckAndChange(elem *etree.Element, configs Configs) {
	for _, v := range configs.Cfgs {
		if elem.SelectAttr("id").Value == v.ID {
			var re = regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
			elem.SelectAttr("style").Value = re.ReplaceAllString(elem.SelectAttr("style").Value, `fill:`+v.Fill)
			fmt.Printf("  ATTR: %s=%s\n", "style", elem.SelectAttr("style").Value)
		}
	}
}

func UpdateDoc(doc *etree.Document, config Configs) {
	root := doc.SelectElement("svg")
	fmt.Println("ROOT element:", root.Tag)
	for _, g := range root.SelectElements("g") {
		//g.FindElements
		fmt.Println(" CHILD element:", g.Tag)
		/*if title := g.SelectElement("title"); title != nil {
			lang := title.SelectAttrValue("lang", "unknown")
			fmt.Printf("  TITLE: %s (%s)\n", title.Text(), lang)
		}*/
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
