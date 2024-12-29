package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	doneChan := make(chan bool, 1)

	//Notify the channel on SIGINT and SIGTERM
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	//Simulate a long-running Goroutine
	go func() {
		for {
			select {
			case <-doneChan:
				fmt.Println("Cleaning up resources....")
				time.Sleep(2 * time.Second) // Simulate cleanup time
				fmt.Println("Cleanup complete.Existing")
				return
			default:
				fmt.Println("Working....")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	//wait for a signal
	sig := <-signalChan
	fmt.Printf("\nReceived signal: %s\n", sig)
	doneChan <- true            //Notify Goroutine to stop
	time.Sleep(3 * time.Second) //Ensure cleanup completes
}
