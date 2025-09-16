package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {

	listener, err := net.Listen("tcp", "localhost:6052")
	if err != nil {
		log.Fatal("Error while listening: ", err)
	}

	go broadcaster()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Error while connecting: ", err)
			continue
		}
		fmt.Println("Connection is accepted!")
		go handleConnection(connection)
	}

}

func handleConnection(conn net.Conn) {
	eachClientCh := make(chan string)
	go clientWriter(conn, eachClientCh)
	clientAddress := conn.RemoteAddr().String()
	eachClientCh <- "Welcome to the room " + clientAddress + "\n"
	messages <- clientAddress + " joined the room" + "\n"
	entering <- eachClientCh

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		messages <- clientAddress + ": " + scanner.Text()
	}

	leaving <- eachClientCh
	messages <- clientAddress + "left the room" + "\n"
	defer conn.Close()
}

func broadcaster() {
	clients := make(map[client]struct{})
	for {
		select {
		case newMember := <-entering:
			clients[newMember] = struct{}{}
		case leftMember := <-leaving:
			delete(clients, leftMember)
			close(leftMember)
		case newMessage := <-messages:
			for cli := range clients {
				cli <- newMessage
			}
		}
	}
}

func clientWriter(conn net.Conn, messagesCh <-chan string) {
	for newMessage := range messagesCh {
		conn.Write([]byte(newMessage))
		// fmt.Fprintln(conn, newMessage)
	}

}
