package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Get(filename string) (Config, error) {
	configSource, err := ioutil.ReadFile(filename)
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
