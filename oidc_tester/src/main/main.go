package main

import (
	"fmt"
	"net/http"

	redirect_flow "avinash.com/testgo/pkg"
	"avinash.com/testgo/utils/config"
	"github.com/gin-gonic/gin"
)

// Token response from the auth server
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// store the configs in a global
var configuration config.Configurations

// store redirect uri
var logout_url, widget_logout_url string

func main() {
	router := gin.Default()

	initConfig()

	fmt.Printf("the configuration is loaded as %v \n", configuration)

	router.LoadHTMLGlob("./templates/*.html")
	router.Static("/assets", "./assets/")

	//login widget
	router.GET("/login", LoginCustom)
	//landing url
	router.GET("/landing", redirect_flow.LoginRedirect)
	//callback url
	router.GET("/callback", redirect_flow.RedirectRetrieve)
	//post code, client id and secret
	router.POST("/gettokens", redirect_flow.RetrieveTokens)
	//get user info
	router.POST("/userinfo", redirect_flow.RetrieveUserInfo)
	//app landing page
	router.GET("/applanding", RedirectSession)
	router.Run(configuration.AppConfigs.Host + ":" + configuration.AppConfigs.Port)
}

func initConfig() {
	configuration = config.ReadConfig()
	logout_url = configuration.OIDCConfigs.Authserver_domain + "" + configuration.OIDCConfigs.Logout_url
	widget_logout_url = configuration.AppConfigs.AppUrl + ":" + configuration.AppConfigs.Port + configuration.CustomWidgetConfigs.Post_logout_url
}
func LoginCustom(c *gin.Context) {
	g := gin.H{
		"widgetversion": configuration.CustomWidgetConfigs.Widgetversion,
		"domain":        configuration.OIDCConfigs.Authserver_domain,
		"client_id":     configuration.CustomWidgetConfigs.Client_id,
		"redirect_url":  configuration.CustomWidgetConfigs.Redirect_url,
	}
	c.HTML(http.StatusOK, "login.html", g)
}

// get the code from redirect url and display page
func RedirectSession(c *gin.Context) {
	//get the authorization_code from url
	interaction_code := c.Query("interaction_code")
	//get the state param from url
	state := c.Query("state")

	fmt.Printf("interaction_code: %s and state:%s \n", interaction_code, state)

	c.HTML(http.StatusOK, "applanding.html", gin.H{
		"interaction_code": interaction_code,
		"state":            state,
		"domain":           configuration.OIDCConfigs.Authserver_domain,
		"logouturl":        logout_url,
		"post_logout_url":  widget_logout_url,
		"client_id":        configuration.CustomWidgetConfigs.Client_id,
	})

}
