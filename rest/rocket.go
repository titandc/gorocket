// This package provides a RocketChat rest client.
package rest

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/titandc/gorocket/api"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	Protocol string
	Host     string
	Port     string

	// Use this switch to see all network communication.
	Debug bool

	auth *authInfo
}

type authInfo struct {
	token string
	id    string
}

// The base for the most of the json responses
type statusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RocketchatAuth struct {
	Host   string
	Port   string
	Login  string
	Passwd string
	TLS    bool
}

func NewClient(auth *RocketchatAuth, debug bool) (*Client, error) {
	var protocol string

	if auth.TLS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	c := &Client{Host: auth.Host, Port: auth.Port, Protocol: protocol, Debug: debug}
	if err := c.Login(api.UserCredentials{Email: auth.Login, Password: auth.Passwd}); err != nil {
		log.Println("Error while login: ", err)
		return nil, err
	}
	return c, nil
}

func (c *Client) getUrl() string {
	return fmt.Sprintf("%v://%v:%v", c.Protocol, c.Host, c.Port)
}

func (c *Client) doRequest(request *http.Request, responseBody interface{}) error {

	if c.auth != nil {
		request.Header.Set("X-Auth-Token", c.auth.token)
		request.Header.Set("X-User-Id", c.auth.id)
	}

	if c.Debug {
		log.Println(request)
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)

	if c.Debug {
		log.Println(string(bodyBytes))
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Request error: " + response.Status)
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, responseBody)
}

// Return random number
func (c *Client) GetRandomId() string {
	b := make([]byte, 17)
	rand.Read(b)
	return hex.EncodeToString(b)
}
