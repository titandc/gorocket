package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/titandc/gorocket/api"
	"net/http"
)

type messagesResponse struct {
	statusResponse
	ChannelName string `json:"channel"`
	Messages    []api.Message `json:"messages"`
}

type messageResponse struct {
	statusResponse
	ChannelName string `json:"channel"`
	Message     api.Message `json:"message"`
}

type Page struct {
	Count int
}

// Sends a message to a channel. The name of the channel has to be not nil.
// The message will be html escaped.
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage
func (c *Client) Send(channel *api.Channel, msg string) error {
	values := map[string]string{"channel": channel.Name, "text": msg}
	body, _ := json.Marshal(values)

	request, err := http.NewRequest("POST", c.getUrl()+"/api/v1/chat.postMessage", bytes.NewBuffer(body))
	_ =  err
	response := new(messageResponse)

	return c.doRequest(request, response)
}

// Get messages from a channel. The channel id has to be not nil. Optionally a
// count can be specified to limit the size of the returned messages.
//
// https://rocket.chat/docs/developer-guides/rest-api/channels/history
func (c *Client) GetMessages(channel *api.Channel, page *Page) ([]api.Message, error) {
	u := fmt.Sprintf("%s/api/v1/channels.history?roomId=%s", c.getUrl(), channel.Id)

	if page != nil {
		u = fmt.Sprintf("%s&count=%d", u, page.Count)
	}

	request, _ := http.NewRequest("GET", u, nil)
	response := new(messagesResponse)

	if err := c.doRequest(request, response); err != nil {
		return nil, err
	}

	return response.Messages, nil
}
