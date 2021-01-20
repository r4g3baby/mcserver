package internal

import (
	"github.com/r4g3baby/mcserver/pkg/server"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Debug  bool
		Logger Logger
		Server server.Config
	}

	Logger struct {
		Enabled    bool
		Filename   string
		MaxSize    int
		MaxAge     int
		MaxBackups int
		LocalTime  bool
		Compress   bool
	}
)

func init() {
	viper.SetDefault("debug", false)
	viper.SetDefault("logger.enabled", false)
	viper.SetDefault("logger.filename", "logs/latest.log")
	viper.SetDefault("logger.maxsize", 2)
	viper.SetDefault("logger.maxage", 7)
	viper.SetDefault("logger.maxbackups", 5)
	viper.SetDefault("logger.localtime", false)
	viper.SetDefault("logger.compress", true)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 25565)
}
