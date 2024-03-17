package steam_api_token

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	baseUrl            *url.URL
	authorized         bool
	authorizationToken string
	httpClient         *http.Client
}

func NewClient(baseUrl *url.URL, authorizationToken string) *Client {
	return &Client{
		baseUrl:            baseUrl,
		authorized:         authorizationToken != "",
		authorizationToken: fmt.Sprintf("Bearer %s", authorizationToken),
		httpClient:         &http.Client{},
	}
}

func (c *Client) GetSteamApiToken(ctx context.Context, appId int, description string) (string, error) {
	url := *c.baseUrl
	url.Path = path.Join(url.Path, "token", strconv.Itoa(appId), description)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", err
	}
	if c.authorizationToken != "" {
		req.Header.Add("Authorization", c.authorizationToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var token string
	if _, err = fmt.Fscan(resp.Body, &token); err != nil {
		return "", err
	}

	return token, nil
}
