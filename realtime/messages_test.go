package realtime

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/titandc/gorocket/api"
)

func TestClient_SubscribeToMessageStream(t *testing.T) {

	c := getLoggedInClient(t)

	general := api.Channel{Id: "GENERAL"}
	textToSend := "RealtimeTest"

	messageChannel, err := c.SubscribeToMessageStream(&general)

	assert.Nil(t, err, "Function returned error")
	assert.NotNil(t, messageChannel, "Function didn't returned general")

	sendAndAssertNoError(t, c, &general, textToSend)
	sendAndAssertNoError(t, c, &general, textToSend)
	sendAndAssertNoError(t, c, &general, textToSend)

	receivedMessage1 := <-messageChannel
	receivedMessage2 := <-messageChannel
	receivedMessage3 := <-messageChannel

	assertMessage(t, receivedMessage1)
	assertMessage(t, receivedMessage2)
	assertMessage(t, receivedMessage3)
}
func assertMessage(t *testing.T, message api.Message) {
	assert.NotNil(t, message.Id, "Id was not set")
	assert.Equal(t, "GENERAL", message.ChannelId,"Wrong channel id")
	assert.NotNil(t, message.Timestamp, "Timestamp was not set")
	assert.NotNil(t, message.User.Id, "UserId was not set")
	assert.NotNil(t, message.User.UserName, "Username was not set")
}
func sendAndAssertNoError(t *testing.T, c *Client, channel *api.Channel, text string) {
	m, err := c.SendMessage(channel, text)
	assert.Nil(t, err, "Error while sending message")
	assert.NotNil(t, m, "SendMessage should return a Message object")
}

func TestClient_SubscribeToMessageStream_UnknownChannel(t *testing.T) {

	c := getLoggedInClient(t)
	channel := api.Channel{Id: "unknown"}

	messageChannel, err := c.SubscribeToMessageStream(&channel)

	assert.NotNil(t, err, "Function didn't return error")
	assert.Nil(t, messageChannel, "Function returned channel, but shouldn't")
}
