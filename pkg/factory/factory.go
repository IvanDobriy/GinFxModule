package factory

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IvanDobriy/tag/pkg/tag"
	"github.com/IvanDobriy/tag/pkg/tag/cast"
	"github.com/gin-gonic/gin"
)

type Alias string

type Config struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

func GetConfigFromMap(raw map[string]any) (*Config, error) {
	address, err := tag.GetRequired("addr", raw, func(value any) (*string, error) {
		r, e := cast.ToStringE(value)
		return &r, e
	})
	if err != nil {
		return nil, err
	}
	readTimeout, err := tag.GetRequired("readTimeout", raw, func(value any) (*time.Duration, error) {
		r, e := cast.ToDurationE(value)
		return &r, e
	})
	if err != nil {
		return nil, err
	}
	writeTimeout, err := tag.GetRequired("writeTimeout", raw, func(value any) (*time.Duration, error) {
		r, e := cast.ToDurationE(value)
		return &r, e
	})
	if err != nil {
		return nil, err
	}
	maxHeaderBytes, err := tag.GetRequired("maxHeaderBytes", raw, func(value any) (*int, error) {
		r, e := cast.ToIntE(value)
		return &r, e
	})
	if err != nil {
		return nil, err
	}
	return &Config{Addr: *address, ReadTimeout: *readTimeout, WriteTimeout: *writeTimeout, MaxHeaderBytes: *maxHeaderBytes}, nil
}

type ConfigAliases struct {
	Aliases *map[Alias]*Config
}

func GetConfigAlasesFromMap(raw map[string]any) (*ConfigAliases, error) {
	rawConfig, err := tag.GetRequired("gin", raw, func(value any) (*map[string]any, error) {
		r, e := cast.ToStringMapE(value)
		return &r, e
	})
	if err != nil {
		return nil, err
	}
	aliases := map[Alias]*Config{}
	for name, el := range *rawConfig {
		rawAlias, err := cast.ToStringMapE(el)
		if err != nil {
			return nil, err
		}
		aliases[Alias(name)], err = GetConfigFromMap(rawAlias)
		if err != nil {
			return nil, err
		}
	}
	return &ConfigAliases{Aliases: &aliases}, nil
}

type Factory struct {
	ConfigAliases *ConfigAliases
}

func NewConfig(paths ...string) (ca *ConfigAliases, err error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("configuration file paths is empty")
	}
	var conf *os.File = nil
	for _, path := range paths {
		file, e := os.OpenFile(path, os.O_RDONLY, 0640)
		if e == nil {
			conf = file
			break
		}
	}
	if conf == nil {
		return nil, fmt.Errorf("no config files found, paths: %v", paths)
	}
	defer func() {
		if e := conf.Close(); e != nil {
			err = e
		}
	}()

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
