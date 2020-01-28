package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	host := getEnv("HOST", "localhost")
	c, _, _ := websocket.DefaultDialer.Dial("ws://"+host+":3000/ws", nil)

	go func(c *websocket.Conn) {
		defer c.Close()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(fmt.Sprintf("Unsupported format:'%v'", message))

		}
	}(c)

	select {}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}