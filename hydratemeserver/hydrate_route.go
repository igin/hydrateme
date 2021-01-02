package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hydratemeserver/database"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

func returnError(w http.ResponseWriter, errorMessage string, status int) {
	log.Println(errorMessage)
	http.Error(w, errorMessage, status)
}

// GetHydrateRoute handles creation of hydration alerts
func GetHydrateRoute(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	tasks, err := database.GetHydrationTasksOfUser(userID)
	if err != nil {
		returnError(w, fmt.Sprintf("Couldn't get tasks. Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	messageEncoder := json.NewEncoder(w)
	messageEncoder.Encode(tasks)
}

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

var decoder = schema.NewDecoder()

// CreateHydrateRoute handles creation of hydration alerts
func CreateHydrateRoute(w http.ResponseWriter, r *http.Request) {
	err := handleHydrationRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create hydration task with error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Thanks for you hydration request.")
}

func handleHydrationRequest(r *http.Request) error {
	command, err := parseCommand(r)
	if err != nil {
		return err
	}
	_, err = database.CreateHydrationTask(command.UserID)
	if err != nil {
		return err
	}

	err = respondToCommand(command, "You are going to be hydrated soon.")
	return nil
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
