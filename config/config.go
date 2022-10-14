package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//Config ...
type Config struct {
	AllowOrigins        string `mapstructure:"ALLOW_ORIGINS"`
	Port                string `mapstructure:"PORT"`
	MysqlURI            string `mapstructure:"FEATWS_API_MYSQL_URI"`
	GitlabToken         string `mapstructure:"FEATWS_API_GITLAB_TOKEN"`
	GitlabURL           string `mapstructure:"FEATWS_API_GITLAB_URL"`
	GitlabNamespace     string `mapstructure:"FEATWS_API_GITLAB_NAMESPACE"`
	GitlabPrefix        string `mapstructure:"FEATWS_API_GITLAB_PREFIX"`
	GitlabDefaultBranch string `mapstructure:"FEATWS_API_GITLAB_DEFAULT_BRANCH"`
	GitlabCIScript      string `mapstructure:"FEATWS_API_GITLAB_CI_SCRIPT"`
	ExternalHost        string `mapstructure:"EXTERNAL_HOST"`
	OpenAMURL           string `mapstructure:"OPENAM_URL"`
}

var config = &Config{}

//LoadConfig ...
func LoadConfig() (err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("ALLOW_ORIGINS", "http://localhost")
	viper.SetDefault("PORT", "9007")
	viper.SetDefault("FEATWS_API_MYSQL_URI", "api:api@tcp(localhost:3306)/api")
	viper.SetDefault("FEATWS_API_GITLAB_TOKEN", "")
	viper.SetDefault("FEATWS_API_GITLAB_URL", "")
	viper.SetDefault("FEATWS_API_GITLAB_NAMESPACE", "")
	viper.SetDefault("FEATWS_API_GITLAB_PREFIX", "")
	viper.SetDefault("FEATWS_API_GITLAB_DEFAULT_BRANCH", "main")
	viper.SetDefault("FEATWS_API_GITLAB_CI_SCRIPT", "")
	viper.SetDefault("EXTERNAL_HOST", "localhost:9007")
	viper.SetDefault("OPENAM_URL", "")

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			log.Errorf("Error on Load Config: %v", err)
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
