package config

type Config struct {
	Plans           map[string]*ServicePlan
	BrokerId        string `yaml:"broker_id"`
	BoshTarget      string `yaml:"bosh_target"`
	BoshUser        string `yaml:"bosh_user"`
	BoshPassword    string `yaml:"bosh_password"`
	ServiceUser     string `yaml:"service_user"`
	ServicePassword string `yaml:"service_password"`
}

type ServicePlan struct {
	Name             string
	Description      string
	Release          string
	Stemcell         string
	ManifestTemplate string `yaml:"manifest_template"`
	BindTemplate     string `yaml:"bind_template"`
	UnbindTemplate   string `yaml:"unbind_template"`
	Params           []Param
}

type Param struct {
	Name     string
	Default  interface{}
	Random   bool
	Optional bool
}
