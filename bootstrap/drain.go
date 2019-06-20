package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

const (
	timeout     = 60
	timeoutTick = 5
)

var (
	lastCheck = int32(time.Now().Unix())
)

func killIntegration() error {
	if Proc == nil {
		return nil
	}
	pgid, err := syscall.Getpgid(Proc.Process.Pid)
	if err != nil {
		return err
	}
	if err = syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		return err
	}
	return nil
}

// Terminate ...
func Terminate(killAll bool) {
	// Kill integration process
	if err := killIntegration(); err != nil {
		log.Fatal("Failed to kill integration process")
	}

	// Kill container
	if killAll {
		os.Exit(0)
	}
}

// Tick - Update last checked timestamp.
func Tick() {
	lastCheck = int32(time.Now().Unix())
}

// Drain - each TIMEOUT_TICK interval
// and check that the diff of lastCheck and
// now isn't over TIMEOUT.
func Drain() {
	for {
		now := int32(time.Now().Unix())
		diff := now - lastCheck

		log.Printf("tick (terminating in %v seconds)\n", timeout-diff)

		if diff >= timeout {
			wg.Done()
		}

		time.Sleep(timeoutTick * time.Second)
	}
}
