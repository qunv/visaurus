package log

import (
	"github.com/qunv/visaurus/core/config"
)

func NewProperties(loader config.Loader) (*Properties, error) {
	props := Properties{}
	err := loader.Bind(&props)
	return &props, err
}

type Properties struct {
	Development    bool `default:"false"`
	JsonOutputMode bool `default:"true"`
	CallerSkip     int  `default:"2"`
}

func (l Properties) Prefix() string {
	return "app.logging"
}
