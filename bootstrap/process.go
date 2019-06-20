package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

var (
	// Proc ...
	Proc *exec.Cmd
)

func pipeToStream(
	useStdOut bool,
	onStdio func(string),
	wgPipeReady *sync.WaitGroup,
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

	wgPipeReady.Done()
	wgProcStarted.Wait()

	scanner := bufio.NewScanner(stdio)
	for scanner.Scan() {
		onStdio(scanner.Text())
	}
}

func cmdStream(command []string, onStdout func(string), onStderr func(string)) error {
	Proc = exec.Command(command[0], command[1:]...)
	Proc.Env = os.Environ()
	Proc.Env = append(Proc.Env, fmt.Sprintf("MODEL=%s", EnvModel))

	var wgPipeReady sync.WaitGroup
	var wgProcStarted sync.WaitGroup
	wgPipeReady.Add(2)
	wgProcStarted.Add(1)

	go pipeToStream(true, onStdout, &wgPipeReady, &wgProcStarted, Proc)
	go pipeToStream(false, onStderr, &wgPipeReady, &wgProcStarted, Proc)

	logErr := func(err error) {
		go Status(WSTypeError, err.Error())
		log.Println(err)
	}

	wgPipeReady.Wait()

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

	// Terminate any previous process
	Terminate(false)

	// Run integration CMD
	err := run()
	if err != nil {
		Status(WSTypeError, "Process terminated")
	}
}
