package main

import (
	"fmt"
	"hydratemeserver/slackapi"
	"log"
	"net/http"
	"os"
)

// SendAlertsRoute sends alerts to all registered users
func SendAlertsRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Let's send alerts!")
	sc := slackapi.NewSlackClient(os.Getenv("SLACK_TOKEN"))
	uid := "U013EN1K59T"
	available, err := sc.IsUserAvailable(uid)
	if err != nil {
		log.Fatalf("Failed to retrieve availability: %s", err.Error())
	}

	if available {
		resp, err := sc.SendPrivateMessage(uid, "This is awesome")
		if err != nil {
			log.Fatalf("Failed to send private message: %s", err.Error())
		}

		if !resp.Ok {
			log.Fatalf("Slack responded with ok=false indicating that the message was not sent")
		}
	}
}
