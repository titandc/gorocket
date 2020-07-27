package rest

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/titandc/gorocket/api"
	"net/http"
)

type PostMessageReturn struct {
	MessageId string
	ChannelId string
}

type SendMessage struct {
	MessageInfo *SendMessageInfo `json:"message"`
}

type SendMessageInfo struct {
	Id     string `json:"_id"`
	Text   string `json:"msg"`
	RoomId string `json:"rid"`
	TmId   string `json:"tmid"`
}

type Groups struct {
	GroupInfoReply *GroupsInfoReply `json:"group,omitempty"`
}

type GroupsInfoReply struct {
	Id string `json:"_id"`
}

type ChannelThreadsList struct {
	Threads []Threads `json:"threads"`
}

type Threads struct {
	ThreadId string `json:"_id"`
	Msg      string `json:"msg"`
}

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
func (c *Client) Send(channel *api.Channel, msg string) (*PostMessageReturn, error) {
	values := map[string]string{"channel": channel.Name, "text": msg}
	body, _ := json.Marshal(values)

	request, err := http.NewRequest("POST", c.getUrl()+"/api/v1/chat.postMessage", bytes.NewBuffer(body))
	_ =  err
	response := new(messageResponse)

	err = c.doRequest(request, response)
	if err != nil {
		return nil, err
	}
	ret := &PostMessageReturn{
		MessageId: response.Message.Id,
		ChannelId: response.Message.ChannelId,
	}
	return ret, nil
}

func (c *Client) randomId() string {
	b := make([]byte, 17)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *Client) GetGroupsNameInfo(channelName string) (string, error) {
	request, err := http.NewRequest("GET", c.getUrl()+"/api/v1/groups.info?roomName=" + channelName, nil)
	_ =  err
	response := new(Groups)
	if err := c.doRequest(request, response); err != nil {
		return "", err
	}
	return response.GroupInfoReply.Id, nil
}

func (c *Client) ListGroupThread(groupId string) (*ChannelThreadsList, error) {
	request, err := http.NewRequest("GET", c.getUrl()+"/api/v1/chat.getThreadsList?rid=" + groupId, nil)
	_ =  err

	response := new(ChannelThreadsList)
	if err := c.doRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) SendReplyMsg(groupId, threadId, msg string) error {
	replyMsg := &SendMessage{
		MessageInfo: &SendMessageInfo{
			Id:     c.randomId(),
			Text:   msg,
			RoomId: groupId, // channel ID
			TmId:   threadId, // Thread/message ID
		},
	}

	body, _ := json.Marshal(replyMsg)
	request, err := http.NewRequest("POST", c.getUrl()+"/api/v1/chat.sendMessage", bytes.NewBuffer(body))
	_ = err
	response := new(messageResponse)

	return c.doRequest(request, response)
}

func (c *Client) SendReply(channel *api.Channel, originalMsg, newMsg string) error {
	groupId, err := c.GetGroupsNameInfo(channel.Name)
	if err != nil {
		return err
	}

	threadsList, err := c.ListGroupThread(groupId)
	if err != nil {
		return err
	}

	for _, thread := range threadsList.Threads {
		if originalMsg == thread.Msg {
			return c.SendReplyMsg(groupId, threadsList.Threads[0].ThreadId, newMsg)
		}
	}
	return fmt.Errorf("Original thread message not found.")
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
