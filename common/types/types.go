package types

type Variable struct {
	Util  string   `yaml:"util,omitempty"`
	Value string   `yaml:"value,omitempty"`
	Vars  []string `yaml:"vars,omitempty"`
}

type File struct {
	Name  string `yaml:"name"`
	Src   string `yaml:"src"`
	Dest  string `yaml:"dest"`
	Perms int32  `yaml:"perms"`
}

type Directory struct {
	Name  string `yaml:"name"`
	Dest  string `yaml:"dest"`
	Src   string `yaml:"src,omitempty"`
	Perms int32  `yaml:"perms,omitempty"`
}

type Addons []struct {
	Roast string `yaml:"roast"`
}

type Config struct {
	Tags        []string            `yaml:"tags"`
	Variables   map[string]Variable `yaml:"vars"`
	Directories []Directory         `yaml:"mkdirs,omitempty"`
	Files       []File              `yaml:"files,omitempty"`
	Restart     string              `yaml:"restart,omitempty"`
	Addons      Addons              `yaml:"addons,omitempty"`
	Scripts     []string            `yaml:"scripts,omitempty"`
}
