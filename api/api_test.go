package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	*testing.T
}

func TestUnkownRouteReturns404(t *testing.T) {
	recorder := requestAPI(t, "GET", "/", nil)
	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
}

func TestHydrateEndpointReturns200(t *testing.T) {
	recorder := requestAPI(t, "GET", "/hydrate", nil)
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}

func TestHydrateEndpointReturnsMessage(t *testing.T) {
	recorder := requestAPI(t, "GET", "/hydrate", nil)
	assert.Equal(t, "Let's hydrate!", recorder.Body.String())
}

func TestHydrateEndpointPostReturns200WithUrlParams(t *testing.T) {
	recorder := requestAPI(t, "POST", "/hydrate", strings.NewReader("response_url=something"))
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}

type ReturnMessage struct {
	Text string `json:"text"`
}

func TestHydrateEndpointSendsMessageToReturnUrl(t *testing.T) {
	apiCalled := false
	var message ReturnMessage

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalled = true

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&message)
		if err != nil {
			t.Fatalf("Got an error while decoding the message sent back to slack")
		}

		fmt.Fprintf(w, "Thanks for your response")
	}))
	defer ts.Close()
	recorder := requestAPI(t, "POST", fmt.Sprintf("/hydrate"), strings.NewReader(fmt.Sprintf("response_url=%s", ts.URL)))
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	assert.True(t, apiCalled)
	assert.Equal(t, ReturnMessage{Text: "something"}, message)

}

func requestAPI(t *testing.T, method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	router := BuildRouter()
	router.ServeHTTP(rr, req)
	return rr
}
