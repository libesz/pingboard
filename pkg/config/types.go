package config

type Target struct {
	SvgId    string `yaml:"id"`
	Fill     string `yaml:"fill"`
	Method   string `yaml:"method,omitempty"`
	EndPoint string `yaml:"endpoint,omitempty"`
}

type Config struct {
	SvgPath string   `yaml:"svgpath"`
	Targets []Target `yaml:"targets"`
}
