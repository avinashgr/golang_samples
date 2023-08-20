package config

import (
	"fmt"

	viper "github.com/spf13/viper"
)

type Configurations struct {
	OIDCConfigs        OIDCConfigurations
	APIEndPointConfigs APIEndPointConfigs
	AppConfigs         AppConfigurations
	CustomWidgetConfigs CustomWidgetConfigs
}

type APIEndPointConfigs struct {
	Api_domain   string
	Userinfo_url string
}

type AppConfigurations struct {
	Port   		string
	AppUrl 		string
	Logout_url	string
	Host 		string
}
type OIDCConfigurations struct {
	Client_id         string
	Client_secret     string
	Authserver_domain string
	Scope             string
	Grant_type        string
	Authorize_url     string
	Token_url         string
	Logout_url        string
	OIDC_info		  string
}

type CustomWidgetConfigs struct {
	Widgetversion	  string
	Client_id         string
	Redirect_url	  string
	Post_logout_url	  string
}

func ReadConfig() Configurations {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	viper.AutomaticEnv()
	var configuration Configurations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	fmt.Printf("the configuration value %s", configuration.OIDCConfigs.Authserver_domain)
	fmt.Printf("the configuration value for userinfo url is  %s", configuration.APIEndPointConfigs.Userinfo_url)
	return configuration
}
