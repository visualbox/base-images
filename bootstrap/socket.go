package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

const (
	// WSTypeTick ...
	WSTypeTick = "TICK"
	// WSTypeInit ...
	WSTypeInit = "INIT"
	// WSTypeTerminate ...
	WSTypeTerminate = "TERMINATE"
	// WSTypeInfo ...
	WSTypeInfo = "INFO"
	// WSTypeOutput ...
	WSTypeOutput = "OUTPUT"
	// WSTypeWarning ...
	WSTypeWarning = "WARNING"
	// WSTypeError ...
	WSTypeError = "ERROR"
)

type wsMessage struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Room   string `json:"room,omitempty"`
	I      string `json:"i,omitempty"`
	Data   string `json:"data,omitempty"`
}

var (
	c *websocket.Conn
)

func wsSendMessage(message wsMessage) {
	b, err := json.Marshal(message)

	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

// Status ...
func Status(statusType string, data string) {
	wsSendMessage(wsMessage{
		Action: "message",
		Room:   EnvToken,
		I:      EnvI,
		Type:   statusType,
		Data:   data,
	})
}

// Output ...
func Output(data string) {
	wsSendMessage(wsMessage{
		Action: "message",
		Room:   EnvToken,
		I:      EnvI,
		Type:   WSTypeOutput,
		Data:   data,
	})
}

func onMessage(data []byte) {
	var message wsMessage
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Println(err)
		return
	}

	switch message.Type {
	case WSTypeTick:
		Tick()
		// go Status(StatusTypeTick, "") // Should be made own message type
	case WSTypeTerminate:
		// Kill integration process and container
		// if 'i' is not present or same as EnvI.
		if message.I == "" || message.I == EnvI {
			Terminate()
		}
	default:
		log.Println("Unknown message type:", message.Type)
	}
}

// InitSocket ...
func InitSocket() {

	c, _, err := websocket.DefaultDialer.Dial("wss://fmgqmvup1i.execute-api.eu-west-1.amazonaws.com/prod", nil)
	if err != nil {
		log.Println(err)
		Terminate()
	}
	defer c.Close()

	// Join
	wsSendMessage(wsMessage{
		Action: "join",
		Room:   EnvToken,
	})

	// Send INIT
	wsSendMessage(wsMessage{
		Action: "message",
		Room:   EnvToken,
		I:      EnvI,
		Type:   WSTypeInit,
	})

	// Receive loop
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			onMessage(message)
		}
	}()
}
