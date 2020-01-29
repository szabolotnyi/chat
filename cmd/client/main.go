package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

var (
	msgChan    = make(chan string, 100)
	recivedMsg = make(chan string, 100)
)

func main() {
	host := getEnv("HOST", "localhost")
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+host+":3000/ws", nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	defer conn.Close()

	go func(c *websocket.Conn) {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			recivedMsg <- string(message)

		}
	}(conn)

	go func() {
		for msg := range msgChan {
			conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
		}
	}()

	writeFile()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)

		msgChan <- text
	}
}

func writeFile() {
	// open output file
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}

	go func() {
		// close fo on exit and check for its returned error
		defer fo.Close()

		// write a chunk
		for msg := range recivedMsg {
			if _, err := fo.Write([]byte(msg)); err != nil {
				panic(err)
			}
		}

	}()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
