package config

import (
	_ "embed"
	"github.com/r4g3baby/mcserver/pkg/server"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

//go:embed config.yaml
var defaultConfig []byte

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
}

func WriteDefaultConfig() error {
	path, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(path, "config.yaml"))
	if err != nil {
		return err
	}
	if _, err := file.Write(defaultConfig); err != nil {
		return err
	}
	return file.Close()
}
