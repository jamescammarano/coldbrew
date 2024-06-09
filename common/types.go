package common

type Variable struct {
	Util  string   `yaml:"util,omitempty"`
	Value string   `yaml:"value,omitempty"`
	Vars  []string `yaml:"vars,omitempty"`
}
