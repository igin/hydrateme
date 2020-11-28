package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"cloud.google.com/go/datastore"
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
	UserID         string `schema:"user_id,required"`
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
	x, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("Failed to dump request with error %s", err.Error())
		http.Error(w, "Failed to parse form in hydrate route with error %s", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form in hydrate route with error %s", err.Error())
		http.Error(w, "Failed to parse form in hydrate route with error %s", http.StatusBadRequest)
		return
	}

	log.Println(fmt.Sprintf("%q", x))

	var command slackCommandFormValues
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(&command, r.PostForm)
	if err != nil {
		log.Printf("Failed to parse form in hydrate route with error %s", err.Error())
		http.Error(w, "Failed to parse form in hydrate route with error %s", http.StatusBadRequest)
		return
	}

	_, err = createHydrationTask(&command)
	if err != nil {
		log.Printf("Failed to create hydration task with error %s", err.Error())
		http.Error(w, "Failed to create hydration task with error %s", http.StatusBadRequest)
		return
	}

	respondToCommand(&command, "something")

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

// HydrationTask represents a notification configuration for a specific slack user
type HydrationTask struct {
	DateCreated           time.Time
	SlackUserID           string
	SlackWorkspaceID      string
	AlertFrequencyMinutes int
	StartTime             time.Time
	EndTime               time.Time
}

func createHydrationTask(command *slackCommandFormValues) (*datastore.Key, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, datastore.DetectProjectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	task := &HydrationTask{
		SlackUserID: command.UserID,
	}
	key := datastore.IncompleteKey("HydrationTask", nil)
	return client.Put(ctx, key, task)
}

type respondable struct {
	ResponseURL string
}

func respondToCommand(command *slackCommandFormValues, message string) {
	returnURL := command.ResponseURL

	response := slackResponseMessage{
		Text: message,
	}

	payloadBuf := new(bytes.Buffer)
	messageEncoder := json.NewEncoder(payloadBuf)
	messageEncoder.Encode(response)

	http.Post(returnURL, "application/json", payloadBuf)
}
