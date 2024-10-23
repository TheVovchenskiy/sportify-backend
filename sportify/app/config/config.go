package config

import (
	"fmt"
	"sync"
	"time"

	// consulapi "github.com/hashicorp/consul/api"
	// vaultapi "github.com/hashicorp/vault/api"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote" // for consul
)

const (
	consulUpdateInterval = 5 * time.Second

	consulConfigPath = "config/sportify"
	// vaultConfigPath = "config/vault"
)

var (
	// This is a global variable that contains the configuration for the Sportify application.
	// It is initialized in the init function and can oly be accessed through the GetGlobalConfig function.
	// The key idea here is that each time config dynamically changes we update pointers to the new config.
	// But if some request is in progress we still use old config, that was stored in context in some middleware.
	globalConfig *Config //nolint:gochecknoglobals

	configMutex sync.RWMutex //nolint:gochecknoglobals
)

func newConfig() (*Config, error) {
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal viper config: %w", err)
	}

	return &config, nil
}

// UpdateGlobalConfig updates the global configuration for the Sportify application.
func UpdateGlobalConfig() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	updateConfig, err := newConfig()
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}
	globalConfig = updateConfig
	return nil
}

// GetGlobalConfig returns the pointer to the global configuration for the Sportify application.
func GetGlobalConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return globalConfig
}

// TODO: reconfigure config structure

// Config is a struct that contains the configuration for the Sportify application.
type Config struct {
	App struct {
		Port          string `mapstructure:"port"`
		APIPrefix     string `mapstructure:"api_prefix"`
		PathPhotos    string `mapstructure:"path_photos"`
		IAMToken      string `mapstructure:"iam_token"`
		FolderID      string `mapstructure:"folder_id"`
		URLPrefixFile string `mapstructure:"url_prefix_file"`
	}

	Logger struct {
		ProductionMode  bool     `mapstructure:"production_mode"`
		LoggerOutput    []string `mapstructure:"logger_output"`
		LoggerErrOutput []string `mapstructure:"logger_err_output"`
	}

	Postgres struct {
		URL      string `mapstructure:"url"`
		DB       string `mapstructure:"db"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	Bot struct {
		APIURL string `mapstructure:"api_url"`
		Port   string `mapstructure:"port"`
		Token  string `mapstructure:"token"`
	}

	Consul struct {
		Address string `mapstructure:"address"`
	}
}

func InitConfig(configFilePaths []string) error {
	initDefaults()
	initEnvironment()

	configFileErr := InitConfigFile(configFilePaths)

	consulErr := initConsul()

	if configFileErr != nil && consulErr != nil {
		return fmt.Errorf("init config file and consul: %w, %w", configFileErr, consulErr)
	}

	err := UpdateGlobalConfig()
	if err != nil {
		return fmt.Errorf("update global config: %w", err)
	}

	return nil
}

func initDefaults() {
	viper.SetDefault("logger.production_mode", false)
	viper.SetDefault("logger.logger_output", []string{"stdout"})
	viper.SetDefault("logger.logger_err_output", []string{"stderr"})

	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.api_prefix", "/api/v1/")
	viper.SetDefault("app.path_photos", "./photos")

	viper.SetDefault("bot.port", "8090")

	viper.SetDefault("consul.address", "localhost:8500")
}

func initEnvironment() {
	viper.MustBindEnv("bot.token", "BOT_TOKEN")

	viper.MustBindEnv("postgres.db", "POSTGRES_DB")
	viper.MustBindEnv("postgres.user", "POSTGRES_USER")
	viper.MustBindEnv("postgres.password", "POSTGRES_PASSWORD")
}

// InitConfigFile searches for the config file in the given paths and initializes the viper config with it.
func InitConfigFile(configPaths []string) error {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

func initConsul() error {
	consulAddress := viper.GetString("consul.address")

	err := viper.AddRemoteProvider("consul", consulAddress, consulConfigPath)
	if err != nil {
		return fmt.Errorf("add consul provider: %w", err)
	}
	viper.SetConfigType("json")
	err = viper.ReadRemoteConfig()
	if err != nil {
		return fmt.Errorf("read consul config: %w", err)
	}
	return nil
}

// WatchRemoteConfig runs go routines that watch for changes in the remote config and updates the viper config.
func WatchRemoteConfig(logger *mylogger.MyLogger) {
	// TODO: watch only when remote is set
	logger.Info("watching remote config")
	go watchConsul(logger)
}

func watchConsul(logger *mylogger.MyLogger) {
	logger.Infof("watching consul config every %s", consulUpdateInterval)

	ticker := time.Tick(consulUpdateInterval)
	for range ticker {
		err := viper.WatchRemoteConfig()
		if err != nil {
			logger.Errorf("error watching consul config: %v\n", err)
			continue
		}

		logger.With("config", viper.AllSettings()).Debug("consul config updated")

		err = UpdateGlobalConfig()
		if err != nil {
			logger.Errorf("error updating global config: %v\n", err)
		}
	}
}
