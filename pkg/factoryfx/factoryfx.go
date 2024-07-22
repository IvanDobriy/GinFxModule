package factoryfx

import (
	"context"

	"github.com/IvanDobriy/GinFxModule/pkg/factory"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ConfigResult struct {
	fx.Out
	Config *factory.ConfigAliases
}

func newConfig(paths ...string) func() (ConfigResult, error) {
	return func() (ConfigResult, error) {
		config, err := factory.NewConfig(paths...)
		if err != nil {
			return ConfigResult{}, err
		}
		return ConfigResult{
			Config: config,
		}, nil
	}
}

type FactoryParatertes struct {
	fx.In
	Config *factory.ConfigAliases
}

type FactoryResult struct {
	fx.Out
	Factory *factory.Factory
}

func newFactory(params FactoryParatertes) (FactoryResult, error) {
	f, err := factory.NewFactory(params.Config)
	if err != nil {
		return FactoryResult{}, err
	}
	return FactoryResult{
		Factory: f,
	}, nil
}


func NewServerInwoker(alias factory.Alias, handlerConfigurators ...func(engine *gin.Engine) error) func(lc fx.Lifecycle, f *factory.Factory) error {
	return func(lc fx.Lifecycle, f *factory.Factory) error {
		engine, err := f.NewGinEngine(alias)
		if err != nil {
			return err
		}
		for _, configure := range handlerConfigurators {
			if err = configure(engine); err != nil {
				return err
			}
		}
		s, err := f.NewServer(alias, engine)
		if err != nil {
			return err
		}
		lc.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					go s.ListenAndServe()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return s.Shutdown(ctx)
				},
			},
		)
		return nil
	}
}

func NewModule(configPaths ...string) func() fx.Option {
	moduleName := "gin server module"
	return func() fx.Option {
		module := fx.Module(moduleName,
			fx.Provide(
				newConfig(configPaths...),
				newFactory,
			),
		)
		return module
	}
}
