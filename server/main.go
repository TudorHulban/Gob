package main

import (
	"fmt"
	"log"
	"main/pkg/transport/servertcp"
	"os"
)

func main() {
	s, err := servertcp.NewServer(nil)
	if err != nil {
		log.Println("Exiting:", err)
		os.Exit(1)
	}

	comms, stop, errServe := s.Serve()
	if errServe != nil {
		log.Println("Exiting:", errServe)
		os.Exit(2)
	}

	for {
		select {
		case msg := <-comms:
			{
				fmt.Println("message:", msg)
			}
		case <-stop:
			break
		}
	}
}
