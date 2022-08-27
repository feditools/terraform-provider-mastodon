package provider

import (
	"context"
	"github.com/mattn/go-mastodon"
)

func (p *mastodonProvider) server() string {
	return  p.schema + "://" + p.domain
}

func (p *mastodonProvider) getAccessToken() string {
	p.accessTokenLock.RLock()
	defer p.accessTokenLock.RUnlock()

	return p.accessToken
}

func (p *mastodonProvider) newAuthenticatedClient(ctx context.Context, clientID, clientSecret, accessToken string) (*mastodon.Client, error) {
	// use given access token
	if accessToken != "" {
		return mastodon.NewClient(&mastodon.Config{
			Server:       p.server(),
			ClientID:     clientID,
			ClientSecret: clientSecret,
			AccessToken:  accessToken,
		}), nil
	}

	// check cache for access token
	if cachedAccessToken := p.getAccessToken(); cachedAccessToken != "" {
		return mastodon.NewClient(&mastodon.Config{
			Server:       p.server(),
			ClientID:     clientID,
			ClientSecret: clientSecret,
			AccessToken:  cachedAccessToken,
		}), nil
	}

	p.accessTokenLock.Lock()
	defer p.accessTokenLock.Unlock()

	client := mastodon.NewClient(&mastodon.Config{
		Server:       p.server(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	err := client.AuthenticateApp(ctx)
	if err != nil {
		return nil, err
	}

	p.accessToken = client.Config.AccessToken

	return client, nil
}

func (p *mastodonProvider) newUnauthenticatedClient() *mastodon.Client {
	return mastodon.NewClient(&mastodon.Config{
		Server: p.server(),
	})
}
