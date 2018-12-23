package infra

import (
	"golang.org/x/oauth2"
)

const (
	slackAuthURL = "https://slack.com/oauth/authorize"
)

// GetOauth2Config return google oauth2 config
func GoogleOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     Env.GoogleClientID,
		ClientSecret: Env.GoogleClientSecret,
		Scopes:       Env.GoogleScopes,
		RedirectURL:  Env.GoogleRedirectedURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  Env.GoogleAuthURL,
			TokenURL: Env.GoogleTokenURL,
		},
	}
}

// SlackOauth2Config reutrn slack oauth2 oonfig
func SlackOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     Env.SlackClientID,
		ClientSecret: Env.SlackClientSecret,
		Scopes:       Env.SlackScope,
		RedirectURL:  Env.SlackRedirectedURL,
		Endpoint: oauth2.Endpoint{
			AuthURL: slackAuthURL,
		},
	}
}
