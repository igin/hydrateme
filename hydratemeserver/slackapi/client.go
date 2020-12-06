package slackapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/schema"
)

// SlackClient enables easy connection to the slack API
type SlackClient struct {
	Token string
}

const slackAPIBaseURL = "https://slack.com/api/"

// NewSlackClient creates a new slack client with the specified token
func NewSlackClient(token string) *SlackClient {
	return &SlackClient{Token: token}
}

// SlackMethod describes the method executed on the slack API
type SlackMethod string

// Available methods on the slack API
const (
	ChatPostMessage SlackMethod = "chat.postMessage"
	UserGetPresence             = "users.getPresence"
)

// Post sends a post request with the specified body to the specified endpoint
// and parses the json response into the response object
func (sc SlackClient) Post(method SlackMethod, payload interface{}, response interface{}) error {
	payloadBuf := new(bytes.Buffer)
	payloadEncoder := json.NewEncoder(payloadBuf)
	payloadEncoder.Encode(payload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", slackAPIBaseURL, method), payloadBuf)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.Token))

	dump, _ := httputil.DumpRequest(req, true)
	log.Printf("--> %s %s", req.RequestURI, dump)

	resp, err := client.Do(req)
	dump, _ = httputil.DumpResponse(resp, true)
	log.Printf("<-- %s %s", resp.Status, dump)

	if err != nil {

		return err
	}

	responseDecoder := json.NewDecoder(resp.Body)
	return responseDecoder.Decode(response)
}

var queryEncoder = schema.NewEncoder()

// Get gets response for a method and uses the payload as query parameters
func (sc SlackClient) Get(method SlackMethod, payload interface{}, response interface{}) error {
	values := url.Values{}
	queryEncoder.Encode(payload, values)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?%s", slackAPIBaseURL, method, values.Encode()), nil)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.Token))

	dump, _ := httputil.DumpRequest(req, true)
	log.Printf("--> %s %s", req.RequestURI, dump)

	resp, err := client.Do(req)
	dump, _ = httputil.DumpResponse(resp, true)
	log.Printf("<-- %s %s", resp.Status, dump)

	if err != nil {

		return err
	}

	responseDecoder := json.NewDecoder(resp.Body)
	return responseDecoder.Decode(response)
}
