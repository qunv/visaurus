package datasource

import "github.com/qunv/visaurus/core/config"

func NewDataSourceProperties(loader config.Loader) (*DatasourceProperties, error) {
	props := DatasourceProperties{}
	err := loader.Bind(props)
	return &props, err
}

// DatasourceProperties ..
type DatasourceProperties struct {
	Host     string `validate:"required" default:"local"`
	Port     int    `validate:"required" default:"3306"`
	Database string
	Username string
	Password string
	Params   string
}

func (d DatasourceProperties) Prefix() string {
	return "app.datasource"
}
