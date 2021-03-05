package main

import (
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

	var serverStopping bool

	for !serverStopping {
		select {
		case msg := <-comms:
			{
				log.Println("message:", string(msg))
			}
		case <-stop:
			{
				log.Println("should exit now")
				serverStopping = true
				break
			}
		}
	}

	log.Println("exiting...")
}
