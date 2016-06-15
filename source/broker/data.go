package broker

import (
	"github.com/s-matyukevich/bosh-broker/source/config"
	"github.com/s-matyukevich/bosh-broker/source/tmpl"
)

type ServiceInstance struct {
	Config         *config.ServicePlan
	Templates      *Templates
	InstanceParams map[string]interface{}
	LastTaskId     string
}

type Templates struct {
	ReleaseTmpl  *tmpl.Template
	StemcellTmpl *tmpl.Template
	ManifestTmpl *tmpl.Template
	BindTmpl     *tmpl.Template
	UnbindTmpl   *tmpl.Template
}
