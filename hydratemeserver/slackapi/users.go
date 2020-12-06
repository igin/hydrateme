package slackapi

import (
	"fmt"
	"log"
)

type userGetPresencePayload struct {
	UserID string `schema:"user"`
}

type userGetPresenceResponse struct {
	Ok       bool   `json:"ok"`
	Presence string `json:"presence"`
}

// IsUserAvailable returns true if the user is available. False otherwise
func (sc *SlackClient) IsUserAvailable(userID string) (bool, error) {
	payload := userGetPresencePayload{
		UserID: userID,
	}
	response := userGetPresenceResponse{}
	err := sc.Get(UserGetPresence, payload, &response)
	if err != nil {
		return false, fmt.Errorf("Failed to retrieve presence of user: %s", err.Error())
	}
	if !response.Ok {
		return false, fmt.Errorf("Failed to retrieve presence of user: response is not ok")
	}

	log.Printf("presence: %s", response.Presence)
	return response.Presence == "active", nil
}
