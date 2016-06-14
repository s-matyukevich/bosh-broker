package config

type Config struct {
	Plans           map[string]*ServicePlan
	BrokerId        string
	BoshTarget      string
	BoshUser        string
	BoshPassword    string
	ServiceUser     string
	ServicePassword string
}

type ServicePlan struct {
	Name             string
	Description      string
	Release          string
	Stemcell         string
	ManifestTemplate string
	BindTemplate     string
	UnbindTemplate   string
	Params           []Param
}

type Param struct {
	Name    string
	Default interface{}
	Random  bool
}
