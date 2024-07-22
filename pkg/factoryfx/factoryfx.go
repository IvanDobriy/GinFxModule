package factoryfx

import (
	"context"

	"github.com/IvanDobriy/GinFxModule/pkg/factory"
	"github.com/ethereum/go-ethereum/beacon/engine"
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
	config *factory.ConfigAliases
}

type FactoryResult struct {
	fx.Out
	Factory *factory.Factory
}

func newFactory(params FactoryParatertes) (FactoryResult, error) {
	f, err := factory.NewFactory(params.config)
	if err != nil {
		return FactoryResult{}, err
	}
	return FactoryResult{
		Factory: f,
	}, nil
}

type InwokerConfig struct {
	alieas      factory.Alias
	conifgPaths []string
}

func NewServerInwoker(config InwokerConfig, handlerConfigurators ...func(engine *gin.Engine) error) func(lc fx.Lifecycle, f *factory.Factory) error {
	return func(lc fx.Lifecycle, f *factory.Factory) error {
		engine, err := f.NewGinEngine(config.alieas)
		if err != nil {
			return err
		}
		for _, configure := range handlerConfigurators {
			if err = configure(engine); err != nil {
				return err
			}
		}
		s, err := f.NewServer(config.alieas, engine)
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
