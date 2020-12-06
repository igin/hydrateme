package slackapi_test

import (
	"encoding/json"
	"fmt"
	"hydratemeserver/slackapi"
	"hydratemeserver/testspyserver"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUserAvailableCallsCorrectURL(t *testing.T) {
	spy := testspyserver.NewSpyServer()
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	sc.IsUserAvailable("testuserid")

	assert.True(t, spy.WasCalled)
	assert.Equal(t, "/users.getPresence", spy.LastRequest.URL.Path)
}

func TestIsUserAvailableAddsCorrectQueryParameter(t *testing.T) {
	spy := testspyserver.NewSpyServer()
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	sc.IsUserAvailable("testuserid")

	assert.True(t, spy.WasCalled)
	assert.Equal(t, "user=testuserid", spy.LastRequest.URL.Query().Encode())
}

func TestIsUserAvailableParsesPositiveResponseCorrectly(t *testing.T) {
	spy := testspyserver.NewSpyServerReturning(func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(struct {
			Ok       bool   `json:"ok"`
			Presence string `json:"presence"`
		}{
			Ok:       true,
			Presence: "active",
		})
	})
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	isAvailable, err := sc.IsUserAvailable("testuserid")

	assert.True(t, isAvailable)
	assert.Nil(t, err)
}

func TestIsUserAvailableParsesNegativeResponseCorrectly(t *testing.T) {
	spy := testspyserver.NewSpyServerReturning(func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(struct {
			Ok       bool   `json:"ok"`
			Presence string `json:"presence"`
		}{
			Ok:       true,
			Presence: "away",
		})
	})
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	isAvailable, err := sc.IsUserAvailable("testuserid")

	assert.False(t, isAvailable)
	assert.Nil(t, err)
}

func TestIsUserReturnsErrorIfResponseIsNotOk(t *testing.T) {
	spy := testspyserver.NewSpyServerReturning(func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(struct {
			Ok       bool   `json:"ok"`
			Presence string `json:"presence"`
		}{
			Ok:       false,
			Presence: "away",
		})
	})
	defer spy.Close()

	sc := slackapi.NewSlackClient("testtoken")
	sc.BaseURL = fmt.Sprintf("%s/", spy.URL)

	_, err := sc.IsUserAvailable("testuserid")

	assert.NotNil(t, err)
}
