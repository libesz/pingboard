package config

import (
	"io"
	"io/ioutil"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"gopkg.in/yaml.v2"
)

func Get(reader io.Reader) (Config, error) {
	configSource, err := ioutil.ReadAll(reader)
	if err != nil {
		return Config{}, err
	}
	config := Config{}

	err = yaml.Unmarshal([]byte(configSource), &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func Validate(configData Config, svg *etree.Document) error {
	var tmpUpdateRules []svgmanip.Target
	for _, v := range configData.Targets {
		tmpUpdateRules = append(tmpUpdateRules, svgmanip.Target{ID: v.ID})
	}
	if err := svgmanip.CheckDoc(svg, tmpUpdateRules); err != nil {
		return err
	}
	return nil
}
