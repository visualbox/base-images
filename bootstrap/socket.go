package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

const (
	// WSTypeTick ...
	WSTypeTick = "TICK"
	// WSTypeInit ...
	WSTypeInit = "INIT"
	// WSTypeTerminate ...
	WSTypeTerminate = "TERMINATE"
	// WSTypeRestart ...
	WSTypeRestart = "RESTART"
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
	Meta   string `json:"meta,omitempty"`
}

var (
	socket gowebsocket.Socket
)

func sendMessage(message wsMessage) {
	b, err := json.Marshal(message)
	if err != nil {
		log.Println("Unable to marshal message:", err)
		return
	}

	socket.SendText(string(b))
}

// Status ...
func Status(statusType string, data string) {
	sendMessage(wsMessage{
		Action: "message",
		Room:   EnvToken,
		I:      EnvI,
		Type:   statusType,
		Data:   data,
	})
}

// Output ...
func Output(data string) {
	sendMessage(wsMessage{
		Action: "message",
		Room:   EnvToken,
		I:      EnvI,
		Type:   WSTypeOutput,
		Data:   data,
	})
}

func onConnected(socket gowebsocket.Socket) {
	// Join
	sendMessage(wsMessage{
		Action: "join",
		Room:   EnvToken,
		Meta:   EnvI,
	})

	// Send INIT
	sendMessage(wsMessage{
		Action: "message",
		Room:   EnvToken,
		I:      EnvI,
		Type:   WSTypeInit,
	})

	wg.Done()
}

func onError(err error, socket gowebsocket.Socket) {
	log.Println("Recieved connect error", err)
	Terminate(true)
}

func onTextMessage(text string, socket gowebsocket.Socket) {
	b := []byte(text)
	var message wsMessage
	err := json.Unmarshal(b, &message)
	if err != nil {
		log.Println(err)
		return
	}

	if message.Type == "" {
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
			Terminate(true)
		}
	case WSTypeRestart:
		// Not for us
		if message.I != EnvI {
			return
		}

		EnvModel = message.Data

		// Terminate / start integration again
		go StartIntegration()

	default:
		log.Println("Unknown message type:", message.Type)
	}
}

// InitSocket ...
func InitSocket() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket = gowebsocket.New("wss://fmgqmvup1i.execute-api.eu-west-1.amazonaws.com/prod")

	socket.OnConnected = onConnected
	socket.OnConnectError = onError
	socket.OnDisconnected = onError
	socket.OnTextMessage = onTextMessage
	// socket.OnBinaryMessage = func(data [] byte, socket gowebsocket.Socket)

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("Interrupt")
			socket.Close()
			Terminate(true)
		}
	}
}
