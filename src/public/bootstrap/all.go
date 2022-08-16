package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/qunv/visaurus/core"
	"github.com/qunv/visaurus/public/app"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		app.AppOpt(),
		core.PropertiesOpt(),

		//init gin
		fx.Provide(gin.New),
	)
}
