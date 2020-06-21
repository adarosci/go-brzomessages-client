package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SendMessageText send message text
func (c *Connection) SendMessageText(message MessageText) (int, bool) {
	return c.sendMessage(message)
}

// SendMessageTemplate send message text
func (c *Connection) SendMessageTemplate(message MessageTemplate) (int, bool) {
	return c.sendMessage(message)
}

// SendMessageAudio send message text
func (c *Connection) SendMessageAudio(message MessageAudio) (int, bool) {
	return c.sendMessage(message)
}

// SendMessageImage send message text
func (c *Connection) SendMessageImage(message MessageImage) (int, bool) {
	return c.sendMessage(message)
}

// SendMessageVideo send message text
func (c *Connection) SendMessageVideo(message MessageVideo) (int, bool) {
	return c.sendMessage(message)
}

// SendMessageDocument send message text
func (c *Connection) SendMessageDocument(message MessageDocument) (int, bool) {
	return c.sendMessage(message)
}

// SendMessageContact send message text
func (c *Connection) SendMessageContact(message MessageContact) (int, bool) {
	return c.sendMessage(message)
}

func (c *Connection) sendMessage(message interface{}) (int, bool) {
	url := BotURL + "?token=" + c.token + "&wait=true"
	method := "POST"

	data, _ := json.Marshal(message)
	payload := strings.NewReader(string(data))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.privateKey)

	res, err := client.Do(req)
	return res.StatusCode, res.StatusCode == 200
}
