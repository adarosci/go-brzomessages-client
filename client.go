package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

// BrzoMessagesClient client
type BrzoMessagesClient struct {
	accessKey, privateKey string
}

// NewConnect create new connection client messages
func NewConnect(accessKey, privateKey string,
	handlerMessage func(MessageReceived) bool,
	handlerAck func(MessageAck) bool) error {

	conn := connection{accessKey, privateKey}
	auth, err := conn.autenticate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err, string(debug.Stack()))
		<-time.After(time.Millisecond * 250)
		go NewConnect(accessKey, privateKey, handlerMessage, handlerAck)
		return err
	}

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		closeConnection := make(chan bool)

		uri := fmt.Sprintf("%v/%v?token=%v&p=%v", SocketURL, accessKey, accessKey, auth)
		c, _, err := websocket.DefaultDialer.Dial(uri, nil)

		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		go func() {
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					closeConnection <- true
					return
				}

				messageData := MessageData{}
				if er := json.Unmarshal(message, &messageData); er != nil {
					continue
				}

				messageReceived := MessageReceived{}
				if er := json.Unmarshal([]byte(messageData.Body.Data), &messageReceived); er == nil && messageReceived.Type != "" {
					if handlerMessage(messageReceived) {
						conn.confirm(messageReceived.Data.Info.ID, messageReceived.Data.Info.RemoteJid)
					}
					continue
				}

				messageAck := MessageAck{}
				if er := json.Unmarshal([]byte(messageData.Body.Data), &messageAck); er == nil {
					if handlerAck(messageAck) {
						conn.confirm(messageAck.ID, messageAck.To)
					}
					continue
				}
			}
		}()

		go func() {
			// ping
			for {
				<-time.After(time.Second * 20)
				data, _ := json.Marshal(map[string]interface{}{
					"ping": true,
				})
				err := c.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("read:", err)
					closeConnection <- true
					return
				}
			}
		}()

		conn.start()

		select {
		case <-closeConnection:
		case <-interrupt:
		}

		conn.stop()

		<-time.After(time.Millisecond * 250)

		go NewConnect(accessKey, privateKey, handlerMessage, handlerAck)
	}()

	return nil
}
