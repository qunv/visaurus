package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qunv/visaurus/public/app"
	"github.com/qunv/visaurus/public/controllers"
	"github.com/qunv/visaurus/public/middlewares"
	"go.uber.org/fx"
)

type RegisterRouterIn struct {
	fx.In
	App                *app.App
	Engine             *gin.Engine
	SysnonymController *controllers.SymnonymController
}

func RegisterHandler(p RegisterRouterIn) {
	p.Engine.Use(middlewares.EnableCors())
}

func RegisterGinRouters(p RegisterRouterIn) {
	group := p.Engine.Group(p.App.Path())

	v1 := group.Group("v1")

	v1.GET("/symnonyms", p.SysnonymController.List)
}
