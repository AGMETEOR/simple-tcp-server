package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	HOST               = "localhost"
	PORT               = ":8083"
	TRANSPORT_PROTOCOL = "tcp"
)

func handleConnection(c net.Conn) {
	fmt.Println("Handling connection")
	buffer := make([]byte, 1024)
	_, err := c.Read(buffer)
	if err != nil {
		fmt.Println("An error in handleConnection: " + err.Error())
	}

	defer c.Close()

	c.Write([]byte(fmt.Sprintf("THE RESPONSE IS HERE %s", string(buffer[:]))))
}

func main() {
	var wg sync.WaitGroup

	quit := make(chan interface{}, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ln, err := net.Listen(TRANSPORT_PROTOCOL, PORT)
	if err != nil {
		fmt.Println("An error: " + err.Error())
	}

	wg.Add(1)
	go func() {
		<-signals
		fmt.Println("Received termination signal. Shutting down gracefully...")

		close(quit)

		if err := ln.Close(); err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	for {

		fmt.Println("Accepting connection")
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-quit:
				fmt.Println("Server has been shut down.")
				wg.Wait()
				os.Exit(0)
			default:
				fmt.Println("An error: " + err.Error())
				log.Fatal(err)
			}
		}
		go handleConnection(conn)
	}

}
