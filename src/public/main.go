package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qunv/visaurus/core/log"
	"github.com/qunv/visaurus/public/app"
	"github.com/qunv/visaurus/public/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.All(), StartOpt()).Run()
}

func StartOpt() fx.Option {
	return fx.Invoke(func(lc fx.Lifecycle, app *app.App, engine *gin.Engine) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Infof("Application will be served at %d", app.Port())
				go func() {
					if err := engine.Run(fmt.Sprintf(":%d", app.Port())); err != nil {
						log.Fatalf("Start app got an error [%v]", err)
					}
				}()
				return nil
			},
		})
	})
}
