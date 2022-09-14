package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/qunv/visaurus/adapter/repositories"
	"github.com/qunv/visaurus/core"
	"github.com/qunv/visaurus/public/app"
	"github.com/qunv/visaurus/public/controllers"
	"github.com/qunv/visaurus/public/datasource"
	"github.com/qunv/visaurus/public/models"
	"github.com/qunv/visaurus/public/router"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		app.AppOpt(),
		core.PropertiesOpt(),

		// init datasource
		datasource.DatasourceOpt(),


		//init repository
		fx.Provide(repositories.NewWordRepositories),

		//init gin
		fx.Provide(gin.New),

		//init controller
		fx.Provide(controllers.NewSymnonymController),

		fx.Invoke(models.InitializeDatabase),
		fx.Invoke(router.RegisterHandler),
		fx.Invoke(router.RegisterGinRouters),
	)
}
