package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config ...
// The above code defines a struct type Config with various fields mapped to environment variables
// using mapstructure tags in Go.
// @property {string} AllowOrigins - This property specifies the allowed origins for CORS (Cross-Origin
// Resource Sharing) requests. It is a string value that can contain multiple origins separated by
// commas.
// @property {string} Port - The port number on which the server will listen for incoming requests.
// @property {string} MysqlURI - This property is used to specify the URI for connecting to a MySQL
// database. It is likely used by the application to store and retrieve data from the database.
// @property {string} Migrate - The Migrate property is used to specify whether or not to run database
// migrations when the application starts up. It is a string value that can be set to "true" or
// "false". If set to "true", the application will attempt to run any pending database migrations
// before starting up. If set
// @property {string} GitlabToken - GitlabToken is a string property that holds the access token for
// the Gitlab API. This token is used to authenticate requests made to the Gitlab API.
// @property {string} GitlabURL - The URL of the GitLab instance that the application will interact
// with.
// @property {string} GitlabNamespace - GitlabNamespace is a configuration property that specifies the
// namespace or group name in GitLab where the project is located. This property is used by the
// application to interact with the GitLab API and perform various operations on the project, such as
// creating a new branch, merging changes, or triggering a pipeline.
// @property {string} GitlabPrefix - GitlabPrefix is a string property that represents the prefix to be
// used for GitLab project names.
// @property {string} GitlabDefaultBranch - GitlabDefaultBranch is a configuration property that
// specifies the default branch name for GitLab repositories. This property is used in the context of
// the GitLab integration in the application.
// @property {string} GitlabCIScript - GitlabCIScript is a property in the Config struct that
// represents the GitLab CI script that will be used for building and testing the project. It is
// specified in the configuration file using the key "FEATWS_API_GITLAB_CI_SCRIPT".
// @property {string} ExternalHost - The external host name or IP address of the server where the
// application is running. This is used for generating URLs and links in the application.
// @property {string} OpenAMURL - OpenAMURL is a string property that represents the URL of the OpenAM
// authentication server.
// @property {string} AuthMode - This property specifies the authentication mode used by the API. It is
// defined in the configuration file and can be set to a specific value using the `mapstructure` tag.
type Config struct {
	AllowOrigins        string `mapstructure:"ALLOW_ORIGINS"`
	Port                string `mapstructure:"PORT"`
	MysqlURI            string `mapstructure:"FEATWS_API_MYSQL_URI"`
	Migrate             string `mapstructure:"MIGRATE"`
	GitlabToken         string `mapstructure:"FEATWS_API_GITLAB_TOKEN"`
	GitlabURL           string `mapstructure:"FEATWS_API_GITLAB_URL"`
	GitlabNamespace     string `mapstructure:"FEATWS_API_GITLAB_NAMESPACE"`
	GitlabPrefix        string `mapstructure:"FEATWS_API_GITLAB_PREFIX"`
	GitlabDefaultBranch string `mapstructure:"FEATWS_API_GITLAB_DEFAULT_BRANCH"`
	GitlabCIScript      string `mapstructure:"FEATWS_API_GITLAB_CI_SCRIPT"`
	ExternalHost        string `mapstructure:"EXTERNAL_HOST"`
	OpenAMURL           string `mapstructure:"OPENAM_URL"`
	AuthMode            string `mapstructure:"FEATWS_API_AUTH_MODE"`
}

var config = &Config{}

// LoadConfig ...
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
	viper.SetDefault("MIGRATE", "")
	viper.SetDefault("OPENAM_URL", "")
	viper.SetDefault("FEATWS_API_AUTH_MODE", "none")

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

// GetConfig returns a pointer to a Config object
func GetConfig() *Config {
	return config
}
