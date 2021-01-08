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
	viper.SetDefault("Debug", false)
	viper.SetDefault("Logger.Enabled", false)
	viper.SetDefault("Logger.Filename", "logs/latest.log")
	viper.SetDefault("Logger.MaxSize", 10)
	viper.SetDefault("Logger.MaxAge", 7)
	viper.SetDefault("Logger.MaxBackups", 3)
	viper.SetDefault("Logger.LocalTime", false)
	viper.SetDefault("Logger.Compress", true)
	viper.SetDefault("Server.Host", "0.0.0.0")
	viper.SetDefault("Server.Port", 25565)
}
