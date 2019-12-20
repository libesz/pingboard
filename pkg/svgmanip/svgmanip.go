package svgmanip

import (
	"errors"
	"log"
	"regexp"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
)

func findByID(root *etree.Element, target config.Target) (*etree.Element, error) {
	elem := root.FindElement(".//*[@id='" + target.ID + "']")
	if elem == nil {
		return nil, errors.New("[svgmanip] Could not find object with ID: " + target.ID)
	}
	return elem, nil
}

func findAndGetStyleAttr(root *etree.Element, target config.Target) (*etree.Attr, bool, error) {
	elem, err := findByID(root, target)
	if err != nil {
		return nil, false, err
	}
	style := elem.SelectAttr("style")
	if style != nil {
		return style, true, nil
	}
	log.Println("[svgmanip] Style does not exists for ID: " + string(target.ID) + ". Falling back to search fill property.")
	style = elem.SelectAttr("fill")
	if style != nil {
		return style, false, nil
	}
	log.Println("[svgmanip] Fill does not exists for ID: " + string(target.ID) + ". Creating new attribute.")
	style = elem.CreateAttr("fill", "none")
	return style, false, nil
}

func change(style *etree.Attr, embedded bool, target config.Target) {
	if embedded {
		re := regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
		style.Value = re.ReplaceAllString(style.Value, `fill:`+target.Fill)
	} else {
		style.Value = target.Fill
	}
	log.Printf("[svgmanip] SVG ID: %s, updated fill property: %s\n", target.ID, target.Fill)
}

func CheckAndChange(root *etree.Element, config config.Target) error {
	attr, embedded, err := findAndGetStyleAttr(root, config)
	if err != nil {
		return err
	}
	change(attr, embedded, config)
	return nil
}

func UpdateDoc(doc *etree.Document, targets []config.Target) error {
	root := doc.SelectElement("svg")
	//fmt.Println("ROOT element:", root.Tag)
	for _, target := range targets {
		if err := CheckAndChange(root, target); err != nil {
			return err
		}
	}
	return nil
}

func CheckDoc(doc *etree.Document, config config.Config) error {
	root := doc.SelectElement("svg")
	//fmt.Println("ROOT element:", root.Tag)
	for _, target := range config.Targets {
		if _, err := findByID(root, target); err != nil {
			return err
		}
	}
	return nil
}
