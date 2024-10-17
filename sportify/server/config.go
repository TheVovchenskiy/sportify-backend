package server

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

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

func NewConfig() (*Config, error) {
	filePath := flag.String(
		"configfile",
		"",
		"you should use file path to config file. NOT SECURE example: ./config.example.yaml",
	)
	flag.Parse()

	if filePath == nil {
		return nil, errors.New("flag configfile is nil")
	}

	viper.SetConfigName(*filePath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

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
