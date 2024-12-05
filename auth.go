package zstackcloud

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SignIn - Get a new token for user
func (c *Client) Login() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	// 对密码进行 SHA-512 加密
	hash := sha512.New()
	hash.Write([]byte(c.Auth.Password))
	hashedPassword := hash.Sum(nil)

	hex.EncodeToString(hashedPassword)
	hashed_password := hex.EncodeToString(hashedPassword)

	// 构造请求体数据
	requestData := map[string]interface{}{
		"logInByAccount": map[string]interface{}{
			"accountName": c.Auth.Username,
			"password":    hashed_password,
		},
	}

	rb, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/zstack/v1/accounts/login", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// SignOut - Revoke the token for a user
// 说明: 系统已预设登录Session的阈值，默认为500。若在短时间内大量调用登陆API且不退出登陆 , 达到阈值后会导致无法登陆新的Session。
func (c *Client) Logout() error {
	fmt.Println("正在释放Session")
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/zstack/v1/accounts/sessions/%s", c.HostURL, c.Token), strings.NewReader(string("")))
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	// fmt.Println(string(body))
	if string(body) != "{}" {
		return errors.New(string(body))
	}
	fmt.Println("释放Session成功")
	return nil
}
