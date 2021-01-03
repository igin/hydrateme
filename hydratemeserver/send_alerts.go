package main

import (
	"errors"
	"hydratemeserver/database"
	"hydratemeserver/slackapi"
	"log"
	"net/http"
	"os"
)

// SendAlertsRoute sends alerts to all registered users
func SendAlertsRoute(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetOverdueHydrationTasks()
	if err != nil {
		log.Printf("Couldn't get hydration tasks: %s", err.Error())
		return
	}

	if len(tasks) == 0 {
		log.Println("No overdue hydration tasks")
		return
	}

	for _, task := range tasks {
		err := sendHydrationAlertToUser(task.SlackUserID)
		if err != nil {
			log.Printf("Couldn't send hydration request: %s", err.Error())
			continue
		}

		err = database.SetWasHydrated(task.ID)
		if err != nil {
			log.Printf("Couldn't update hydration task: %s", err.Error())
		}

	}

}

func sendHydrationAlertToUser(slackUserID string) error {
	log.Printf("Sending alert to slack user: %s", slackUserID)
	sc := slackapi.NewSlackClient(os.Getenv("SLACK_TOKEN"))
	available, err := sc.IsUserAvailable(slackUserID)
	if err != nil || !available {
		return errors.New("The user is currently not available")
	}

	resp, err := sc.SendPrivateMessage(slackUserID, "Are you hydrated?")
	if err != nil {
		log.Printf("Failed to send private message: %s", err.Error())
		return err
	}

	if !resp.Ok {
		log.Printf("Slack responded with ok=false indicating that the message was not sent")
		return errors.New("Slack responded with ok=false indicating that the message was not sent")
	}

	return nil
}
