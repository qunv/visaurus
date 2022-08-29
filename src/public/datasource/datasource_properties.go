package datasource

import (
	"github.com/qunv/visaurus/core/config"
)

func NewDataSourceProperties(loader config.Loader) (*DatasourceProperties, error) {
	props := DatasourceProperties{}
	err := loader.Bind(&props)
	return &props, err
}

// DatasourceProperties ..
type DatasourceProperties struct {
	Host     string `validate:"required" default:"localhost"`
	Port     int    `validate:"required" default:"3306"`
	Database string `validate:"required"`
	Username string `validate:"required"`
	Password string
	Params   string
}

func (p DatasourceProperties) Prefix() string {
	return "app.datasource"
}
