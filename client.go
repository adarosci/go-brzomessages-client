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
	handlerAck func(MessageAck) bool,
	updateConn func(conn *Connection)) (*Connection, error) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintln(os.Stderr, p, string(debug.Stack()))
			<-time.After(time.Millisecond * 250)

			go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)

			return
		}
	}()

	conn := &Connection{accessKey, privateKey, false, nil}
	auth, err := conn.autenticate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err, string(debug.Stack()))
		<-time.After(time.Millisecond * 250)
		go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)
		return nil, err
	}

	go func() {
		defer func() {
			if p := recover(); p != nil {
				fmt.Fprintln(os.Stderr, p, string(debug.Stack()))
				<-time.After(time.Millisecond * 250)

				go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)

				return
			}
		}()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		closeConnection := make(chan bool)
		conn.channelClose = closeConnection

		updateConn(conn)

		uri := fmt.Sprintf("%v/%v?token=%v&p=%v", SocketURL, accessKey, accessKey, auth)
		c, _, err := websocket.DefaultDialer.Dial(uri, nil)

		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					fmt.Fprintln(os.Stderr, p, string(debug.Stack()))
					<-time.After(time.Millisecond * 250)

					go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)

					return
				}
			}()
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
						content, _ := json.Marshal(map[string]string{
							"Id":        messageReceived.Data.Info.ID,
							"RemoteJid": messageReceived.Data.Info.RemoteJid,
							"Token":     accessKey,
						})
						c.WriteMessage(websocket.TextMessage, content)
						//conn.confirm(messageReceived.Data.Info.ID, messageReceived.Data.Info.RemoteJid)
					}
					continue
				}

				messageAck := MessageAck{}
				if er := json.Unmarshal([]byte(messageData.Body.Data), &messageAck); er == nil {
					if handlerAck(messageAck) {
						content, _ := json.Marshal(map[string]string{
							"Id":        messageAck.ID,
							"RemoteJid": messageAck.To,
							"Token":     accessKey,
						})
						c.WriteMessage(websocket.TextMessage, content)
						//conn.confirm(messageAck.ID, messageAck.To)
					}
					continue
				}
			}
		}()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					fmt.Fprintln(os.Stderr, p, string(debug.Stack()))
					<-time.After(time.Millisecond * 250)

					go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)

					return
				}
			}()
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

		conn.stop(auth)

		<-time.After(time.Millisecond * 250)
		if !conn.dispose {
			go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)
		}
	}()

	return conn, nil
}
