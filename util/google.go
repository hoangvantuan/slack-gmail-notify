package util

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mdshun/slack-gmail-notify/infra"
	"golang.org/x/oauth2"
	gm "google.golang.org/api/gmail/v1"
)

// GmailSrv return gmail service
func GmailSrv(token *oauth2.Token) (*gm.Service, error) {
	conf := &oauth2.Config{
		ClientID:     infra.Env.GoogleClientID,
		ClientSecret: infra.Env.GoogleClientSecret,
		Scopes:       infra.Env.GoogleScopes,
		RedirectURL:  infra.Env.GoogleRedirectedURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  infra.Env.GoogleAuthURL,
			TokenURL: infra.Env.GoogleTokenURL,
		},
	}

	accessToken, _ := Decrypt(token.AccessToken)
	refreshToken, _ := Decrypt(token.RefreshToken)

	// get gmail
	client := conf.Client(context.Background(), &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	})

	srv, err := gm.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "have error while creare gmail service")
	}

	return srv, nil
}
