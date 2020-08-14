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

	fmt.Println("-> Conectando BRZO", time.Now())

	conn := &Connection{accessKey, privateKey, false, nil}
	auth, err := conn.autenticate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err, string(debug.Stack()))
		<-time.After(time.Second * 5)
		go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)
		return nil, err
	}

	go func() {
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
			for {
				_, message, err := c.ReadMessage()
				fmt.Println("-> Message received BRZO", time.Now())
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
			// ping
			for {
				<-time.After(time.Second * 20)
				data, _ := json.Marshal(map[string]interface{}{
					"ping": true,
				})
				fmt.Println("-> BRZO send ping", time.Now())
				err := c.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("read:", err)
					closeConnection <- true
					return
				}
			}
		}()

		lt, err := conn.start()
		if err != nil {
			fmt.Println("-> Erro ao conectar BRZO", lt, err, time.Now())
			<-time.After(time.Second * 5)
			go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)
		}

		fmt.Println("-> BRZO Conectado", lt, err)

		select {
		case <-closeConnection:
		case <-interrupt:
		}

		conn.stop(auth)

		fmt.Println("-> Reconectando BRZO 5 segundos ->", lt, err, time.Now())

		<-time.After(time.Second * 5)
		if !conn.dispose {
			go NewConnect(accessKey, privateKey, handlerMessage, handlerAck, updateConn)
		}
	}()

	return conn, nil
}
