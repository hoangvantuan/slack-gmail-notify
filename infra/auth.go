package infra

import (
	"golang.org/x/oauth2"
)

var (
	slackAuthURL   = "https://slack.com/oauth/authorize"
	slackScopes    = []string{"commands", "bot"}
	googleScopes   = []string{"https://www.googleapis.com/auth/gmail.readonly", "https://www.googleapis.com/auth/gmail.modify"}
	googleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL = "https://www.googleapis.com/oauth2/v3/token"
)

// GoogleOauth2Config return google oauth2 config
func GoogleOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     Env.GoogleClientID,
		ClientSecret: Env.GoogleClientSecret,
		Scopes:       googleScopes,
		RedirectURL:  Env.APIHost + "/" + Env.GoogleRedirectedPath,
		Endpoint: oauth2.Endpoint{
			AuthURL:  googleAuthURL,
			TokenURL: googleTokenURL,
		},
	}
}

// SlackOauth2Config return slack oauth2 config
func SlackOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     Env.SlackClientID,
		ClientSecret: Env.SlackClientSecret,
		Scopes:       slackScopes,
		RedirectURL:  Env.APIHost + "/" + Env.SlackRedirectedPath,
		Endpoint: oauth2.Endpoint{
			AuthURL: slackAuthURL,
		},
	}
}
