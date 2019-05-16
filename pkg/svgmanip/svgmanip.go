package svgmanip

import (
	"fmt"
	"regexp"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
)

func CheckAndChange(root *etree.Element, target config.Target) {
	elem := root.FindElement(".//*[@id='" + target.SvgId + "']")
	if elem == nil {
		fmt.Printf(" Could not find ID: %s", target.SvgId)
		return
	}
	re := regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
	style := elem.SelectAttr("style")
	if style == nil {
		fmt.Printf(" Style does not exists for ID: %s\n", target.SvgId)
		return
	}
	style.Value = re.ReplaceAllString(elem.SelectAttr("style").Value, `fill:`+target.Fill)
	fmt.Printf(" ID: %s, updated style: %s\n", target.SvgId, elem.SelectAttr("style").Value)
}

func UpdateDoc(doc *etree.Document, config config.Config) {
	root := doc.SelectElement("svg")
	fmt.Println("ROOT element:", root.Tag)
	for _, target := range config.Targets {
		CheckAndChange(root, target)
	}
}
