package config

import (
	"os"

	"github.com/spf13/viper"
)

//Config ...
type Config struct {
	Port     string `mapstructure:"PORT"`
	MongoURI string `mapstructure:"FEATWS_RESOLVER_BRIDGE_MONGO_URI"`
	MongoDB  string `mapstructure:"FEATWS_RESOLVER_BRIDGE_MONGO_DB"`
}

var config = &Config{}

//LoadConfig ...
func LoadConfig() (err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", "9007")
	viper.SetDefault("FEATWS_RESOLVER_BRIDGE_MONGO_URI", "mongodb://localhost:27017/")
	viper.SetDefault("FEATWS_RESOLVER_BRIDGE_MONGO_DB", "resolverBridge")

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			return
		}
	}

	err = viper.Unmarshal(config)

	return
}

// GetConfig ...
func GetConfig() *Config {
	return config
}
