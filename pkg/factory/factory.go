package factory

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Alias string

type Config struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type ConfigAliases struct {
	Aliases *map[Alias]*Config
}

type Factory struct {
	ConfigAliases *ConfigAliases
}

func NewConfig(paths ...string) (*ConfigAliases, error) {
	return &ConfigAliases{}, nil
}

func NewFactory(configAliases *ConfigAliases) (*Factory, error) {
	factory := &Factory{
		ConfigAliases: configAliases,
	}
	return factory, nil
}

func (f *Factory) NewGinEngine(alias Alias) (*gin.Engine, error) {
	return gin.Default(), nil
}

func (f *Factory) NewServer(alias Alias, ginEngine *gin.Engine) (*http.Server, error) {
	config := (*f.ConfigAliases.Aliases)[alias]
	if config == nil {
		return nil, fmt.Errorf("alias: %s not found", alias)
	}
	server := &http.Server{
		Addr:           config.Addr,
		Handler:        ginEngine,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	return server, nil
}
