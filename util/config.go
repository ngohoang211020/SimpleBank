package util

import (
	"github.com/spf13/viper"
	"time"
)

// Configuration Config stores all configuration of the application
// The values are read by viber from a config file of environment variable
type Configuration struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	Port                string        `mapstructure:"PORT"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	MigrationURL        string        `mapstructure:"MIGRATION_URL"`
}

var Config Configuration

// LoadConfig reads configuration from file or environment variables
func (config *Configuration) LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(config)
	return
}
