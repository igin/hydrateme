package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnkownRouteReturns404(t *testing.T) {
	recorder := requestAPI(t, "GET", "/", nil)
	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
}

func TestHydrateEndpointReturns200(t *testing.T) {
	recorder := requestAPI(t, "GET", "/hydrate", nil)
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}

func TestHydrateEndpointReturnsEmptyArray(t *testing.T) {
	recorder := requestAPI(t, "GET", "/hydrate", nil)
	assert.Equal(t, "[]\n", recorder.Body.String())
}

func TestHydrateEndpointReturnsFilledArrayIfTasksExist(t *testing.T) {
	randomUserID := randSeq(10)
	recorder := requestAPI(t, "GET", fmt.Sprintf("/hydrate?userID=%s", randomUserID), nil)
	assert.Equal(t, "[]\n", recorder.Body.String())
}

func TestHydrateEndpointPostReturns200WithUrlParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	randomUserID := randSeq(10)
	recorder := requestAPI(t, "POST", "/hydrate", strings.NewReader(fmt.Sprintf("response_url=%s&user_id=%s", url.QueryEscape(ts.URL), randomUserID)))
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}

type TaskResponse struct {
	SlackUserID string `json:"SlackUserID"`
}

func TestReturnsHydrationTasksForUserAfterTheyAreCreated(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	randomUserID := randSeq(10)
	requestAPI(t, "POST", "/hydrate", strings.NewReader(fmt.Sprintf("response_url=%s&user_id=%s", url.QueryEscape(ts.URL), randomUserID)))
	recorder := requestAPI(t, "GET", fmt.Sprintf("/hydrate?userID=%s", randomUserID), nil)
	var tasks []TaskResponse
	json.Unmarshal(recorder.Body.Bytes(), &tasks)
	assert.Equal(t, 1, len(tasks))
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
	recorder := requestAPI(t, "POST", fmt.Sprintf("/hydrate"), strings.NewReader(fmt.Sprintf("response_url=%s&user_id=someuser", ts.URL)))
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	assert.True(t, apiCalled)
	assert.Equal(t, ReturnMessage{Text: "You are going to be hydrated soon."}, message)
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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	resetDatastore()
}

func resetDatastore() {
	datastoreHost := os.Getenv("DATASTORE_EMULATOR_HOST")
	resetEndpoint := fmt.Sprintf("http://%s/reset", datastoreHost)
	client := &http.Client{}
	_, err := client.Post(resetEndpoint, "application/json", nil)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
