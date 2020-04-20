package config

import (
	"time"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-plugins/config/source/consul"
)

// default service configs
var (
	DefaultVersion          = "latest"
	DefaultRegisterTTL      = time.Second * 10
	DefaultRegisterInterval = time.Second * 5
)

func Init(path string) error {
	return config.Load(consul.NewSource(
		consul.WithPrefix(path),
		consul.StripPrefix(true),
	))
}
