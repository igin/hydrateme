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
	err := handleHydrationRequest(r)
	if err != nil {
		dump, _ := httputil.DumpRequest(r, true)
		log.Fatalf("Failed to create hydration request with error %s. \nRequest:\n %s", err.Error(), dump)
		http.Error(w, fmt.Sprintf("Failed to create hydration task with error %s", err.Error()), http.StatusBadRequest)
	}

	fmt.Fprint(w, "Let's hydrate!")
}

func handleHydrationRequest(r *http.Request) error {
	command, err := parseCommand(r)
	if err != nil {
		return err
	}
	_, err = createHydrationTask(command)
	if err != nil {
		return err
	}

	return respondToCommand(command, "something")
}

func parseCommand(r *http.Request) (*slackCommandFormValues, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	command := &slackCommandFormValues{}
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(command, r.PostForm)
	return command, err
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

func respondToCommand(command *slackCommandFormValues, message string) error {
	returnURL := command.ResponseURL

	response := slackResponseMessage{
		Text: message,
	}

	payloadBuf := new(bytes.Buffer)
	messageEncoder := json.NewEncoder(payloadBuf)
	messageEncoder.Encode(response)

	_, err := http.Post(returnURL, "application/json", payloadBuf)
	return err
}
