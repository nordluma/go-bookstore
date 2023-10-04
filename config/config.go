package config

import (
	"time"

	"github.com/spf13/viper"
)

// Load configuration from file
var InitConfig = initConfig

func initConfig(filename string, additionalDirs []string) error {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Read configs from file
	for _, dir := range additionalDirs {
		viper.AddConfigPath(dir)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Build config
	viper.ConfigFileUsed()
	viper.WatchConfig()

	return nil
}

func getConfigString(key string) string {
	return viper.GetString(key)
}

func getConfigInt(key string) int {
	return viper.GetInt(key)
}

func getConfigDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
