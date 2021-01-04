package app_config

import (
	"bytes"
	"text/tabwriter"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AppTimeoutSeconds int    `envconfig:"APP_TIMEOUT_SECONDS" required:"true" default:"30"`
	RedisAddr         string `envconfig:"REDIS_ADDR" required:"true" default:"localhost:6379"`
	TemplateDir       string `envconfig:"TEMPLATE_DIR" required:"true" default:"./web/template"`
	AssetsDir         string `envconfig:"ASSETS_DIR" required:"true" default:"./web/public/assets"`
}

func (config AppConfig) CalcAppTimeout() time.Duration {
	return time.Duration(config.AppTimeoutSeconds) * time.Second
}

func NewAppConfig() (*AppConfig, error) {
	var appConfig AppConfig
	err := envconfig.Process("", &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

func ShowConfigHelp() (string, error) {

	writer := bytes.NewBuffer(nil)
	tabs := tabwriter.NewWriter(writer, 1, 0, 4, ' ', 0)
	config := AppConfig{}
	if err := envconfig.Usagef("", &config, tabs, envconfig.DefaultTableFormat); err != nil {
		return "", err
	}
	if err := tabs.Flush(); err != nil {
		return "", err
	}
	return writer.String(), nil
}
