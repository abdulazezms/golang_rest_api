package util

import (
	"github.com/spf13/viper"
)

//Config struct holds all configuration for the app.
//Values are read by Viper from a config file or env variables.
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAaddress string `mapstructure:"SERVER_ADDRESS"`
}

//LoadConfig will load configuration from `path` or override their values from env variables, if provided.
func LoadConfig(path, configType, configName string) (config Config, err error) {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType) //json, xml, etc.
	viper.AddConfigPath(path)
	viper.AutomaticEnv() //Automatically override values read from config file with the values of the corresponding env variables, if they exist.
	err = viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		return config, err
	}
	
	err = viper.Unmarshal(&config)
	return config, err
}