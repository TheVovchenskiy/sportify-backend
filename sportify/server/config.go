package server

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once = sync.Once{}

type Config struct {
	ProductionMode  bool
	PortPublic      string
	PortTg          string
	APIPrefix       string
	PathPhotos      string
	FolderID        string
	IAMToken        string
	URLDatabase     string
	LoggerOutput    []string
	LoggerErrOutput []string
}

func NewConfig(configFile string) (*Config, error) {
	once.Do(func() {
		viper.SetConfigName(configFile)
		viper.SetConfigType("yaml")
		viper.AddConfigPath("../config")

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("read config: %w", err))
		}
	})

	return &Config{
		ProductionMode:  viper.GetBool("production_mode"),
		PortPublic:      viper.GetString("port_public"),
		PortTg:          viper.GetString("port_tg"),
		APIPrefix:       viper.GetString("api_prefix"),
		PathPhotos:      viper.GetString("path_photos"),
		FolderID:        viper.GetString("folder_id"),
		IAMToken:        viper.GetString("iam_token"),
		URLDatabase:     viper.GetString("url_database"),
		LoggerOutput:    strings.Split(viper.GetString("logger_output"), " "),
		LoggerErrOutput: strings.Split(viper.GetString("logger_err_output"), " "),
	}, nil
}
