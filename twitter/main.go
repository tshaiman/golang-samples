package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"ts-tweets/tweeter"
)

func main() {
	tweeter.Start()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("Stopping Stream...")

}
