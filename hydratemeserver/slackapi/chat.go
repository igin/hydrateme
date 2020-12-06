package slackapi

type chatPostMessagePayload struct {
	UserID  string `json:"channel"`
	Message string `json:"text"`
}

// ChatPostMessageResponse is the response returned by chat.postMessage
type ChatPostMessageResponse struct {
	Ok bool `json:"ok"`
}

// SendPrivateMessage sends a private message to the specified user
func (sc SlackClient) SendPrivateMessage(userID string, message string) (ChatPostMessageResponse, error) {
	privateMessage := chatPostMessagePayload{
		UserID:  userID,
		Message: message,
	}
	resp := ChatPostMessageResponse{}
	err := sc.Post(ChatPostMessage, privateMessage, &resp)
	return resp, err
}
