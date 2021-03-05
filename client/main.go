package main

import (
	"fmt"
	"io"
	"log"
	"main/pkg/transport/clienttcp"
	"os"
)

func main() {
	c, err := clienttcp.NewClient(nil)
	if err != nil {
		log.Println("Exiting:", err)
		os.Exit(1)
	}

	r, errSend := c.Send([]byte("xxx"))
	if errSend != io.EOF {
		log.Println("Error with send:", errSend, r)
		os.Exit(2)
	}

	fmt.Println("Server response:", r)
}
