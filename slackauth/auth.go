package slackauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultSlackBaseURL = "https://slack.com"
)

var scopes = [...]string{"chat:write", "files:read", "files:write", "commands"}

//go:generate mockgen -source ./slackauth.go -destination ./mock/mock_slackauth.go
type Auth interface {
	AccessURL() string
	InvokeURL() string
	CallbackHTML(TeamID string) string
	ExchangeSlackAuthCodeForToken(code string) (authResp Response, err error)
}

type auth struct {
	AppID        string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	APIURL       string
}

type Response struct {
	Ok          bool       `json:"ok,omitempty"`
	AppID       string     `json:"app_id,omitempty"`
	AuthedUser  AuthedUser `json:"authed_user,omitempty"`
	Scope       string     `json:"scope,omitempty"`
	TokenType   string     `json:"token_type,omitempty"`
	AccessToken string     `json:"access_token,omitempty"`
	BotUserID   string     `json:"bot_user_id,omitempty"`
	Team        TeamInfo   `json:"team,omitempty"`
	Enterprise  string     `json:"enterprise,omitempty"`
	Error       string     `json:"error,omitempty"`
}

type AuthedUser struct {
	ID string `json:"id,omitempty"`
}

type TeamInfo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (a *auth) AccessURL() string {
	return fmt.Sprintf("%s/api/oauth.v2.access", a.APIURL)
}

func (a *auth) InvokeURL() string {
	return fmt.Sprintf("https://slack.com/oauth/v2/authorize?client_id=%s&redirect_uri=%s&scope=%s",
		a.ClientID, a.RedirectURL, strings.Join(scopes[:], ","))
}

func (a *auth) CallbackHTML(TeamID string) string {
	return fmt.Sprintf(
		"<html><head><meta http-equiv=\"refresh\" content=\"0;URL=https://slack.com/app_redirect?app=%s&team=%s\"></head></html>", a.AppID, TeamID)
}

func NewAuth(clientID, clientSecret, redirectURL, APIURL string) Auth {
	return &auth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		APIURL:       APIURL,
	}
}

func (a *auth) ExchangeSlackAuthCodeForToken(code string) (authResp Response, err error) {

	v := url.Values{}
	v.Set("client_id", a.ClientID)
	v.Set("client_secret", a.ClientSecret)
	v.Set("code", code)
	v.Set("redirect_uri", a.RedirectURL)

	body := strings.NewReader(v.Encode())

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/oauth.v2.access", a.APIURL), body)
	if err != nil {
		return authResp, fmt.Errorf("error creating slack access token request err=%v", err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.ClientID, a.ClientSecret)))))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return authResp, fmt.Errorf("failed to get slack Auth response %v", err)
	}

	tokenBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return authResp, fmt.Errorf("fail to get slack access token [%s]: %s", resp.Status, tokenBody)
	}

	err = json.Unmarshal(tokenBody, &authResp)
	if err != nil {
		return authResp, fmt.Errorf("error unmarshal slack Auth response %v", err)
	}

	if !authResp.Ok {
		return authResp, fmt.Errorf(authResp.Error)
	}

	return authResp, nil
}
