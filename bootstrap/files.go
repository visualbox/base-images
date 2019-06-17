package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var (
	// Proc ...
	Proc *exec.Cmd
	tr   = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: tr}
)

func pipeToStream(
	useStdOut bool,
	onStdio func(string),
	wgStreamSetup *sync.WaitGroup,
	wgProcStarted *sync.WaitGroup,
	cmd *exec.Cmd,
) {

	var stdio io.ReadCloser
	var err error

	if useStdOut {
		stdio, err = cmd.StdoutPipe()
	} else {
		stdio, err = cmd.StderrPipe()
	}
	if err != nil {
		panic(err)
	}
	wgStreamSetup.Done()
	scanner := bufio.NewScanner(stdio)
	for scanner.Scan() {
		onStdio(scanner.Text())
	}
}

func cmdStream(command []string, onStdout func(string), onStderr func(string)) error {
	// command = append([]string{"-c"}, command...)
	Proc = exec.Command(command[0], command[1:]...)
	Proc.Env = os.Environ()
	Proc.Env = append(Proc.Env, fmt.Sprintf("MODEL=%s", EnvModel))

	var wgStreamSetup sync.WaitGroup
	var wgProcStarted sync.WaitGroup
	wgStreamSetup.Add(2)
	wgProcStarted.Add(1)

	go pipeToStream(true, onStdout, &wgStreamSetup, &wgProcStarted, Proc)
	go pipeToStream(false, onStderr, &wgStreamSetup, &wgProcStarted, Proc)

	logErr := func(err error) {
		go Status(WSTypeError, err.Error())
		log.Println(err)
	}

	wgStreamSetup.Wait()

	err := Proc.Start()
	if err != nil {
		logErr(err)
	}

	wgProcStarted.Done()

	err = Proc.Wait()
	if err != nil {
		logErr(err)
	}

	return err
}

func run() error {
	cmd := EnvArgs

	// Parse output string as JSON.
	// If it fails, send a status instead.
	onStdout := func(m string) {
		go Status(WSTypeInfo, m)
	}
	onStderr := func(m string) {
		go Status(WSTypeError, m)
	}

	err := cmdStream(cmd, onStdout, onStderr)
	return err
}

// StartIntegration ...
func StartIntegration() {
	go Status(WSTypeInfo, "Starting integration")

	// Run integration CMD
	err := run()
	if err != nil {
		Status(WSTypeError, "Unable to run integration command")
		log.Println(err)
		Terminate()
	}
}
