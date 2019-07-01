package main

import (
	"log"
	"os"
	"sync"
)

var (
	// EnvI - Integration dashboard index.
	EnvI = os.Getenv("I")

	// EnvModel - Initial integration configuration model.
	EnvModel = os.Getenv("MODEL")

	// EnvToken - Instance token.
	EnvToken = os.Getenv("TOKEN")

	// EnvArgs - Entrypoint args.
	EnvArgs = os.Args[1:]

	// EnvRestAPIID - AWS API GW ID
	EnvRestAPIID = os.Getenv("REST_API_ID")

	// EnvWsAPIID - AWS API GW WebSocket Endpoint
	EnvWsAPIID = os.Getenv("WS_API_ID")

	wg = &sync.WaitGroup{}
)

func main() {
	log.Println("Init Unix socket server")
	go InitUnixSocket()

	log.Println("Init WebSocket client")
	wg.Add(1)
	go InitSocket()
	wg.Wait()

	log.Println("Start integration and drain")
	wg.Add(1)
	go StartIntegration()
	go Drain()
	wg.Wait()

	Terminate(true)
}
