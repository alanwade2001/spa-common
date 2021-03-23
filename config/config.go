package services

import (
	"github.com/spf13/viper"
)

// ConfigService s
type ConfigService struct {
}

// Load f
func (cs ConfigService) Load(configPath string) error {

	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return nil
}
