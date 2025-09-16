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

// func main() {
// 	conn, err := net.Dial("tcp", "localhost:6052")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	done := make(chan struct{})
// 	go func() {
// 		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
// 		log.Println("done")
// 		done <- struct{}{} // signal the main goroutine
// 	}()
// 	mustCopy(conn, os.Stdin)
// 	conn.Close()
// 	<-done // wait for background goroutine to finish
// }

// func mustCopy(dst io.Writer, src io.Reader) {
// 	if _, err := io.Copy(dst, src); err != nil {
// 		log.Fatal(err)
// 	}
// }
