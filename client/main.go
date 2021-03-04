package main

import (
	"fmt"
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
	fmt.Println(r, errSend)
}
