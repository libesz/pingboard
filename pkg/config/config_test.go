package config

import (
	"log"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestYaml(t *testing.T) {
	var data = `
svgpath: "a.svg"
targets:
- id: path10
  fill: "#00ff00"
  method: ping
  endpoint: example.com
`

	config := Config{}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		t.Errorf("a")
		log.Fatalf("error unmarshal: %v", err)
	}
	if config.Targets[0].Fill != "#00ff00" {
		t.Errorf("error %v", config)
	}
	if config.SvgPath != "a.svg" {
		t.Errorf("error %v", config)
	}
	if config.Targets[0].Method != "ping" {
		t.Errorf("error %v", config)
	}
	if config.Targets[0].EndPoint != "example.com" {
		t.Errorf("error %v", config)
	}
}
