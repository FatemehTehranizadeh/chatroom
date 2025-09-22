package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:6052")
	if err != nil {
		log.Fatal("Error while dialing: ", err)
	}
	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, connection)
		log.Println("done")
		done <- struct{}{}
	}()
	
	io.Copy(connection, os.Stdin)
	connection.Close()
	<-done
}
