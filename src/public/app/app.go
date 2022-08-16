package app

import (
	"context"
	"github.com/qunv/visaurus/core"
	"github.com/qunv/visaurus/core/config"
	"go.uber.org/fx"
	"net/http"
)

func AppOpt() fx.Option {
	return fx.Options(
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		core.ProvideProps(config.NewAppProperties),
		fx.Provide(New),
	)
}

type App struct {
	props    *config.AppProperties
	context  context.Context
	handlers []func(next http.Handler) http.Handler
}

func New(context context.Context, props *config.AppProperties) *App {
	app := App{context: context, props: props}
	//TODO: add handler here
	return &app
}

func (a App) Name() string {
	return a.props.Name
}

func (a App) Port() int {
	return a.props.Port
}

func (a App) Path() string {
	return a.props.Path
}

func (a App) Context() context.Context {
	return a.context
}

func (a App) Handlers() []func(next http.Handler) http.Handler {
	return a.handlers
}

func (a *App) AddHandler(handlers ...func(next http.Handler) http.Handler) {
	a.handlers = append(a.handlers, handlers...)
}
