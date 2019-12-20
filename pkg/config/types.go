package config

// Method is the type of check
type Handler string

// Target is a single configured entity which is monitored
type Target struct {
	ID       string  `yaml:"id"`
	Fill     string  `yaml:"fill"`
	Handler  Handler `yaml:"method,omitempty"`
	EndPoint string  `yaml:"endpoint,omitempty"`
}

// Config is the input for the main program how to monitor the targets
type Config struct {
	SvgPath string   `yaml:"svgpath"`
	Targets []Target `yaml:"targets"`
}
