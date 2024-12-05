package zstackcloud

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://monvip:8080"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Inventory AuthResponseInventory `json:"inventory"`
}

type AuthResponseInventory struct {
	Token       string `json:"uuid"`
	AccountUuid string `json:"accountUuid"`
	UserUuid    string `json:"userUuid"`
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.Login()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Inventory.Token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// 这部分逻辑在需要使用token的地方需要加请求头
	// token := c.Token
	// req.Header.Set("Authorization", "OAuth"+token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
