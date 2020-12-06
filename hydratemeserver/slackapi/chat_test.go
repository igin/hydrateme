package slackapi_test

import (
	"encoding/json"
	"fmt"
	"hydratemeserver/slackapi"
	"hydratemeserver/testspyserver"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendPrivateMessageCallsCorrectURL(t *testing.T) {
	spy := testspyserver.NewSpyServer()
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	sc.SendPrivateMessage("userid", "textmessage")

	assert.True(t, spy.WasCalled)
	assert.Equal(t, "/chat.postMessage", spy.LastRequest.URL.Path)
}

type postMessagePayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func TestSendPrivateMessageAddsCorrectMessageBody(t *testing.T) {
	spy := testspyserver.NewSpyServer()
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	sc.SendPrivateMessage("userid", "textmessage")

	assert.True(t, spy.WasCalled)

	decoder := json.NewDecoder(strings.NewReader(spy.LastBody))
	actualPayload := postMessagePayload{}
	decoder.Decode(&actualPayload)
	assert.Equal(t, postMessagePayload{
		Channel: "userid",
		Text:    "textmessage",
	}, actualPayload)
}
