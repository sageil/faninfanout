package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sageil/concurrency-patterns/efficient"
	"github.com/sageil/concurrency-patterns/naive"
)

func main() {
	log.Println("g'day")
	patternArg := os.Args[1]
	if patternArg == "" {
		fmt.Println("Usage: go run ./main <pattern>")
		os.Exit(1)
	}
	if patternArg != "efficient" && patternArg != "naive" {
		fmt.Println("Usage: go run ./main <pattern>")
		os.Exit(1)
	}
	if patternArg == "efficient" {
		go efficient.Run_Efficient()
	}
	if patternArg == "naive" {
		go naive.Run_Naive()
	}
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Adios!")
}
