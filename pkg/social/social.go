package social

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/torkelo/grafana-pro/pkg/log"
	"github.com/torkelo/grafana-pro/pkg/models"
	"github.com/torkelo/grafana-pro/pkg/setting"

	"github.com/golang/oauth2"
)

type BasicUserInfo struct {
	Identity string
	Name     string
	Email    string
	Login    string
	Company  string
}

type SocialConnector interface {
	Type() int
	UserInfo(transport *oauth2.Transport) (*BasicUserInfo, error)

	AuthCodeURL(state, accessType, prompt string) string
	NewTransportFromCode(code string) (*oauth2.Transport, error)
}

var (
	SocialBaseUrl = "/login/"
	SocialMap     = make(map[string]SocialConnector)
)

func NewOAuthService() {
	if !setting.Cfg.MustBool("oauth", "enabled") {
		return
	}

	setting.OAuthService = &setting.OAuther{}
	setting.OAuthService.OAuthInfos = make(map[string]*setting.OAuthInfo)

	allOauthes := []string{"github", "google", "twitter"}

	// Load all OAuth config data.
	for _, name := range allOauthes {
		info := &setting.OAuthInfo{
			ClientId:     setting.Cfg.MustValue("oauth."+name, "client_id"),
			ClientSecret: setting.Cfg.MustValue("oauth."+name, "client_secret"),
			Scopes:       setting.Cfg.MustValueArray("oauth."+name, "scopes", " "),
			AuthUrl:      setting.Cfg.MustValue("oauth."+name, "auth_url"),
			TokenUrl:     setting.Cfg.MustValue("oauth."+name, "token_url"),
			Enabled:      setting.Cfg.MustBool("oauth."+name, "enabled"),
		}

		if !info.Enabled {
			continue
		}

		setting.OAuthService.OAuthInfos[name] = info
		options, err := oauth2.New(
			oauth2.Client(info.ClientId, info.ClientSecret),
			oauth2.Scope(info.Scopes...),
			oauth2.Endpoint(info.AuthUrl, info.TokenUrl),
			oauth2.RedirectURL(strings.TrimSuffix(setting.AppUrl, "/")+SocialBaseUrl+name),
		)

		if err != nil {
			log.Error(3, "Failed to init oauth service", err)
			continue
		}

		// GitHub.
		if name == "github" {
			setting.OAuthService.GitHub = true
			SocialMap["github"] = &SocialGithub{Options: options}
		}

		// Google.
		if name == "google" {
			setting.OAuthService.Google = true
			SocialMap["google"] = &SocialGoogle{Options: options}
		}
	}
}

type SocialGithub struct {
	*oauth2.Options
}

func (s *SocialGithub) Type() int {
	return int(models.GITHUB)
}

func (s *SocialGithub) UserInfo(transport *oauth2.Transport) (*BasicUserInfo, error) {
	var data struct {
		Id    int    `json:"id"`
		Name  string `json:"login"`
		Email string `json:"email"`
	}

	var err error
	client := http.Client{Transport: transport}
	r, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &BasicUserInfo{
		Identity: strconv.Itoa(data.Id),
		Name:     data.Name,
		Email:    data.Email,
	}, nil
}

//   ________                     .__
//  /  _____/  ____   ____   ____ |  |   ____
// /   \  ___ /  _ \ /  _ \ / ___\|  | _/ __ \
// \    \_\  (  <_> |  <_> ) /_/  >  |_\  ___/
//  \______  /\____/ \____/\___  /|____/\___  >
//         \/             /_____/           \/

type SocialGoogle struct {
	*oauth2.Options
}

func (s *SocialGoogle) Type() int {
	return int(models.GOOGLE)
}

func (s *SocialGoogle) UserInfo(transport *oauth2.Transport) (*BasicUserInfo, error) {
	var data struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var err error

	reqUrl := "https://www.googleapis.com/oauth2/v1/userinfo"
	client := http.Client{Transport: transport}
	r, err := client.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &BasicUserInfo{
		Identity: data.Id,
		Name:     data.Name,
		Email:    data.Email,
	}, nil
}
