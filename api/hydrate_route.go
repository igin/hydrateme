package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

// StartHydrationRequest is a request to start hydration alerts
type StartHydrationRequest struct {
	User string `json:"user"`
}

type slackCommandFormValues struct {
	TeamID         string `schema:"team_id"`
	TeamDomain     string `schema:"team_domain"`
	EnterpriseID   string `schema:"enterprise_id"`
	EnterpriseName string `schema:"enterprise_name"`
	ChannelID      string `schema:"channel_id"`
	ChannelName    string `schema:"channel_name"`
	UserID         string `schema:"user_id"`
	UserName       string `schema:"user_name"`
	Command        string `schema:"command"`
	Text           string `schema:"text"`
	ResponseURL    string `schema:"response_url,required"`
	TriggerID      string `schema:"trigger_id"`
	APIAppID       string `schema:"api_app_id"`
}

type slackResponseMessage struct {
	Text string `json:"text"`
}

// GetHydrateRoute handles creation of hydration alerts
func GetHydrateRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Let's hydrate!")
}

var decoder = schema.NewDecoder()

// CreateHydrateRoute handles creation of hydration alerts
func CreateHydrateRoute(w http.ResponseWriter, r *http.Request) {
	var command slackCommandFormValues
	err := decoder.Decode(&command, r.URL.Query())
	if err != nil {
		log.Printf("Failed to parse form in hydrate route with error %s", err.Error())
		http.Error(w, "Failed to parse form in hydrate route with error %s", http.StatusBadRequest)
		return
	}

	returnURL := command.ResponseURL

	response := slackResponseMessage{
		Text: "something",
	}

	payloadBuf := new(bytes.Buffer)
	messageEncoder := json.NewEncoder(payloadBuf)
	messageEncoder.Encode(response)

	http.Post(returnURL, "application/json", payloadBuf)
	fmt.Fprint(w, "Let's hydrate!")
}
