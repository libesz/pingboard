package svgmanip

import (
	"errors"
	"log"
	"regexp"

	"github.com/beevik/etree"
)

func findByID(root *etree.Element, target Target) (*etree.Element, error) {
	elem := root.FindElement(".//*[@id='" + target.ID + "']")
	if elem == nil {
		return nil, errors.New("[svgmanip] Could not find object with ID: " + target.ID)
	}
	return elem, nil
}

func getStyleAttr(elem *etree.Element, target Target) (*etree.Attr, bool, error) {
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

func getTitleElem(elem *etree.Element, target Target) (*etree.Element, error) {
	title := elem.SelectElement("title")
	if title != nil {
		return title, nil
	}
	log.Println("[svgmanip] Title does not exists for ID: " + string(target.ID) + ". Creating new element.")
	title = elem.CreateElement("title")
	return title, nil
}

func changeFill(style *etree.Attr, embedded bool, target Target) {
	if embedded {
		re := regexp.MustCompile(`fill:#[0-9a-zA-Z]{6}`)
		style.Value = re.ReplaceAllString(style.Value, `fill:`+target.Fill)
	} else {
		style.Value = target.Fill
	}
	log.Printf("[svgmanip] SVG ID: %s, updated fill property: %s\n", target.ID, target.Fill)
}

func changeTitle(title *etree.Element, target Target) {
	title.CreateText(target.LastChange)
	log.Printf("[svgmanip] SVG ID: %s, updated title element: %s\n", target.ID, target.Fill)
}

func CheckAndChange(root *etree.Element, target Target) error {
	elem, err := findByID(root, target)
	if err != nil {
		return err
	}
	fill, embedded, err := getStyleAttr(elem, target)
	if err != nil {
		return err
	}
	changeFill(fill, embedded, target)
	title, err := getTitleElem(elem, target)
	if err != nil {
		return err
	}
	changeTitle(title, target)
	return nil
}

func UpdateDoc(doc *etree.Document, targets []Target) error {
	root := doc.SelectElement("svg")
	//fmt.Println("ROOT element:", root.Tag)
	for _, target := range targets {
		if err := CheckAndChange(root, target); err != nil {
			return err
		}
	}
	return nil
}

func CheckDoc(doc *etree.Document, targets []Target) error {
	root := doc.SelectElement("svg")
	//fmt.Println("ROOT element:", root.Tag)
	for _, target := range targets {
		if _, err := findByID(root, target); err != nil {
			return err
		}
	}
	return nil
}
