package svgmanip

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
)

func CheckAndChange(root *etree.Element, target config.Target) error {
	elem := root.FindElement(".//*[@id='" + target.SvgId + "']")
	if elem == nil {
		return errors.New("Could not find ID: " + target.SvgId)
	}
	re := regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
	style := elem.SelectAttr("style")
	if style == nil {
		return errors.New("Style does not exists for ID: " + string(target.SvgId))
	}
	style.Value = re.ReplaceAllString(elem.SelectAttr("style").Value, `fill:`+target.Fill)
	fmt.Printf(" ID: %s, updated style: %s\n", target.SvgId, elem.SelectAttr("style").Value)
	return nil
}

func UpdateDoc(doc *etree.Document, config config.Config) error {
	root := doc.SelectElement("svg")
	fmt.Println("ROOT element:", root.Tag)
	for _, target := range config.Targets {
		if err := CheckAndChange(root, target); err != nil {
			return err
		}
	}
	return nil
}
