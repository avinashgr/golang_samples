package redirect_flow

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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

var configs config.Configurations

var client_id, client_secret, loginurl,
	redirect_uri, scope, oidc_info_url, token_url,
	grant_type, userinfo_url, logout_url, app_logout_url string

func init() {
	configs = config.ReadConfig()
	client_id = configs.OIDCConfigs.Client_id
	client_secret = configs.OIDCConfigs.Client_secret
	loginurl = configs.OIDCConfigs.Authserver_domain + "" + configs.OIDCConfigs.Authorize_url
	redirect_uri = configs.AppConfigs.AppUrl + ":" + configs.AppConfigs.Port + "/callback"
	scope = configs.OIDCConfigs.Scope
	oidc_info_url = configs.OIDCConfigs.Authserver_domain + "" + configs.OIDCConfigs.OIDC_info
	token_url = configs.OIDCConfigs.Authserver_domain + "" + configs.OIDCConfigs.Token_url
	grant_type = configs.OIDCConfigs.Grant_type
	userinfo_url = configs.OIDCConfigs.Authserver_domain + "" + configs.APIEndPointConfigs.Userinfo_url
	logout_url = configs.OIDCConfigs.Authserver_domain + "" + configs.OIDCConfigs.Logout_url
	app_logout_url = configs.AppConfigs.AppUrl + ":" + configs.AppConfigs.Port + configs.AppConfigs.Logout_url
}

func LoginRedirect(c *gin.Context) {
	g := gin.H{
		"loginurl":      loginurl,
		"response_type": "code",
		"client_id":     client_id,
		"redirect_uri":  redirect_uri,
		"scope":         scope,
		"oidc_info_url": oidc_info_url,
	}
	c.HTML(http.StatusOK, "landing.html", g)
}

// get the code from redirect url and display page
func RedirectRetrieve(c *gin.Context) {
	//get the authorization_code from url
	auth_code := c.Query("code")
	//get the state param from url
	state := c.Query("state")

	fmt.Printf("authorization_code: %s and state:%s \n", auth_code, state)

	c.HTML(http.StatusOK, "code.html", gin.H{
		"code":         auth_code,
		"state":        state,
		"tokenurl":     "/gettokens",
		"redirect_uri": redirect_uri,
		"token_url":    token_url,
	})

}

// retrieve the tokens for api calls
func RetrieveTokens(c *gin.Context) {

	auth_code := c.PostForm("code")
	var redirect_uri = redirect_uri
	var token_url = token_url
	fmt.Printf("the authcode from the call is %s \n", auth_code)

	formData := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"scope":         {scope},
		"code":          {auth_code},
		"redirect_uri":  {redirect_uri},
		"grant_type":    {grant_type},
	}

	fmt.Printf("the url for the post %s \n", token_url)

	resp, err := http.PostForm(token_url, formData)
	if err != nil {
		fmt.Println("error posting to the url \n", err)
		c.HTML(http.StatusBadRequest, "", gin.H{})
		return
	}

	defer resp.Body.Close()
	var result TokenResponse
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		fmt.Println("error reading the body", err)
		return
	}

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Cannot unmarshal JSON")
	}
	if result.IDToken == "" {
		c.Redirect(http.StatusFound, "/landing")
		return
	}
	claims := strings.Split(result.IDToken, ".")[1]

	decodedClaims, _ := b64.StdEncoding.DecodeString(claims)
	fmt.Printf("the decoded claims from the token are %s \n", string(decodedClaims))

	c.HTML(http.StatusOK, "tokens.html", gin.H{
		"id_token":     result.IDToken,
		"access_token": result.AccessToken,
		"data":         formData,
		"token_url":    token_url,
		"info_from_ID": string(decodedClaims),
		"userinfourl":  "/userinfo",
	})

}

// retrieve user info from access token
func RetrieveUserInfo(c *gin.Context) {

	access_token := c.PostForm("access_token")
	id_token := c.PostForm("id_token")
	userinfo_url := userinfo_url
	logout_url := logout_url
	app_logout_url := app_logout_url

	fmt.Printf("the id_token is %s \n", id_token)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", userinfo_url, nil)
	req.Header.Set("Authorization", "Bearer "+access_token)

	fmt.Printf("the request url is %v \n", userinfo_url)
	resp, err := client.Do(req)
	if err != nil {
		c.Redirect(http.StatusFound, "/landing")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Redirect(http.StatusFound, "/landing")
		return
	}
	// fmt.Printf("the user info from /userinfo call %s \n", string(body))

	c.HTML(http.StatusOK, "userinfo.html", gin.H{
		"userinfo":      string(body),
		"logouturl":     logout_url,
		"id_token_hint": id_token,
		"callback_url":  app_logout_url,
	})

}
