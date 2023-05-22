package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config type contains various configuration parameters for a Go application, including database
// and GitLab API credentials, authentication mode, and external host information.
//
// Property:
//   - AllowOrigins: This property specifies the allowed origins for CORS requests. It is a string value
//
// that can contain multiple origins separated by commas.
//   - Port: The port number on which the server will listen for incoming requests.
//   - MysqlURI: The URI for connecting to the MySQL database used by the API.
//   - Migrate: it's used to specify whether to run database migrations or not. If the value is set to "true", the application will run database migrations on startup. If the value is set to "false", the application will not run database migrations.
//   - GitlabToken: This's a token used for authentication with GitLab API. It allows the application to access GitLab resources on behalf of a user or a bot account.
//   - GitlabURL: The URL of the GitLab instance that the API will interact with.
//   - GitlabNamespace: The namespace or group name in GitLab where the project is located.
//   - GitlabPrefix: The GitlabPrefix property is a string that represents the prefix to be used for all GitLab API requests. It is used to specify the namespace or group where the project's located. For example, if the GitlabPrefix is set to "mygroup/myproject", all API requests will be made.
//   - GitlabDefaultBranch: This property represents the default branch name for a GitLab repository. When creating a new repository, GitLab will use this branch as the default branch.
//   - GitlabCIScript - GitlabCIScript is a property in the Config struct that represents the GitLab CI script that will be used for building and testing the project. It is specified in the configuration file using the key "FEATWS_API_GITLAB_CI_SCRIPT".
//   - ExternalHost - This property represents the external host name or IP address of the server where the application is running. It is used to configure the application to listen on a specific network interface or to generate URLs that can be accessed from outside the server.
//   - OpenAMURL: The URL of the OpenAM server used for authentication.
//   - AuthMode - This property specifies the authentication mode used by the API. It can have values like "jwt", "oauth2", "basic", etc.
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

// LoadConfig loads configuration settings from a file and sets default values if they are no present.
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

// GetConfig returns a pointer to a Config object.
func GetConfig() *Config {
	return config
}
