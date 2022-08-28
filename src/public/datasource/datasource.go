package datasource

import (
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatasourceOpt() fx.Option {
	return fx.Options(
		fx.Provide(NewDataSourceProperties),
		fx.Provide(NewDatasource),
	)
}

type DatasourceOut struct {
	fx.Out
	Connection *gorm.DB
}

func NewDatasource(props *DatasourceProperties) (DatasourceOut, error) {
	var connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		props.Username, props.Password, props.Host, props.Port, props.Database, props.Params)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	out := DatasourceOut{}
	out.Connection = db
	return out, err
}
