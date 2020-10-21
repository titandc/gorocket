package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CacheChannelId struct {
	Name string
	Id   string
}

type SendMessage struct {
	MessageInfo *SendMessageInfo `json:"message"`
}

type SendMessageInfo struct {
	Id     string `json:"_id"`
	Text   string `json:"msg"`
	RoomId string `json:"rid"`
	TmId   string `json:"tmid"`
	TShow  bool   `json:"tshow"`
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

var cacheChannelId []CacheChannelId

/* Prefer function ThreadReply() */
func (c *Client) ThreadReplyByOriginalChannelMessage(channelName, originalText, replyText string) error {
	// Get channel ID
	channelName = strings.TrimPrefix(channelName, "#")
	groupId, err := c.GetGroupsNameInfo(channelName)
	if err != nil {
		return err
	}

	// List threads on channel
	threadsList, err := c.ListThreadsInGroup(groupId)
	if err != nil {
		return err
	}

	// Search in thread the original message to want reply.
	for _, thread := range threadsList.Threads {
		if originalText == thread.Msg {
			return c.ThreadReply(groupId, thread.ThreadId, replyText)
		}
	}

	// Original thread message not found.
	return fmt.Errorf("Thread not found: The group no contains message with original text.")
}

func (c *Client) ListThreadsInGroup(groupId string) (*ChannelThreadsList, error) {
	request, err := http.NewRequest("GET", c.getUrl()+"/api/v1/chat.getThreadsList?rid="+groupId, nil)
	_ = err

	response := new(ChannelThreadsList)
	if err := c.doRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

/* Note groupId is ChannelId */
func (c *Client) ThreadReply(channelName, threadId, msg string) error {
	channelName = strings.TrimPrefix(channelName, "#")
	groupId, err := c.GetGroupsNameInfo(channelName)
	if err != nil {
		return fmt.Errorf("Failed to get channel ID for name <%s>.\n", channelName)
	}

	replyMsg := &SendMessage{
		MessageInfo: &SendMessageInfo{
			Id:     c.GetRandomId(),
			Text:   msg,
			TShow:  true,
			RoomId: groupId,  // channel ID
			TmId:   threadId, // Thread/message ID
		},
	}

	body, err := json.Marshal(replyMsg)
	if err != nil {
		return fmt.Errorf("Failed to reply, json error: %s\n", err.Error())
	}
	request, err := http.NewRequest("POST", c.getUrl()+"/api/v1/chat.sendMessage", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	response := new(messageResponse)
	return c.doRequest(request, response)
}

func (c *Client) GetGroupsNameInfo(channelName string) (string, error) {
	// Check if Id already search
	for _, channel := range cacheChannelId {
		if channel.Name == channelName && channel.Id != "" {
			return channel.Id, nil
		}
	}

	// Is not in cache, search ID and add in cache
	request, err := http.NewRequest("GET", c.getUrl()+"/api/v1/groups.info?roomName="+channelName, nil)
	_ = err
	response := new(Groups)
	if err := c.doRequest(request, response); err != nil {
		return "", fmt.Errorf("Failed to get groups name <%s> infos: %s\n", channelName)
	}

	cacheChannelId = append(cacheChannelId, CacheChannelId{
		Name: channelName,
		Id:   response.GroupInfoReply.Id})
	
	return response.GroupInfoReply.Id, nil
}
