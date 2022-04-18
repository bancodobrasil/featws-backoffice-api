package config

import (
	"os"

	"github.com/spf13/viper"
)

//Config ...
type Config struct {
	Port                string `mapstructure:"PORT"`
	MongoURI            string `mapstructure:"FEATWS_API_MONGO_URI"`
	MongoDB             string `mapstructure:"FEATWS_API_MONGO_DB"`
	GitlabToken         string `mapstructure:"FEATWS_API_GITLAB_TOKEN"`
	GitlabURL           string `mapstructure:"FEATWS_API_GITLAB_URL"`
	GitlabNamespace     string `mapstructure:"FEATWS_API_GITLAB_NAMESPACE"`
	GitlabPrefix        string `mapstructure:"FEATWS_API_GITLAB_PREFIX"`
	GitlabDefaultBranch string `mapstructure:"FEATWS_API_GITLAB_DEFAULT_BRANCH"`
	GitlabCIScript      string `mapstructure:"FEATWS_API_GITLAB_CI_SCRIPT"`
}

var config = &Config{}

//LoadConfig ...
func LoadConfig() (err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", "9007")
	viper.SetDefault("FEATWS_API_MONGO_URI", "mongodb://localhost:27017/")
	viper.SetDefault("FEATWS_API_MONGO_DB", "resolverBridge")
	viper.SetDefault("FEATWS_API_GITLAB_TOKEN", "")
	viper.SetDefault("FEATWS_API_GITLAB_URL", "")
	viper.SetDefault("FEATWS_API_GITLAB_NAMESPACE", "")
	viper.SetDefault("FEATWS_API_GITLAB_PREFIX", "")
	viper.SetDefault("FEATWS_API_GITLAB_DEFAULT_BRANCH", "main")
	viper.SetDefault("FEATWS_API_GITLAB_CI_SCRIPT", "")

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
