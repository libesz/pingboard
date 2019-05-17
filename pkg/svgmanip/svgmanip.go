package svgmanip

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
)

func findAndCheck(root *etree.Element, target config.Target) (*etree.Attr, error) {
	elem := root.FindElement(".//*[@id='" + target.SvgId + "']")
	if elem == nil {
		return nil, errors.New("Could not find ID: " + target.SvgId)
	}
	style := elem.SelectAttr("style")
	if style == nil {
		return nil, errors.New("Style does not exists for ID: " + string(target.SvgId))
	}
	return style, nil
}

func change(style *etree.Attr, target config.Target) {
	re := regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
	style.Value = re.ReplaceAllString(style.Value, `fill:`+target.Fill)
	fmt.Printf(" ID: %s, updated style: %s\n", target.SvgId, style.Value)
}

func CheckAndChange(root *etree.Element, config config.Target) error {
	attr, err := findAndCheck(root, config)
	if err != nil {
		return err
	}
	change(attr, config)
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

func CheckDoc(doc *etree.Document, config config.Config) error {
	root := doc.SelectElement("svg")
	fmt.Println("ROOT element:", root.Tag)
	for _, target := range config.Targets {
		if _, err := findAndCheck(root, target); err != nil {
			return err
		}
	}
	return nil
}
