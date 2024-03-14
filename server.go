package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	conn.SetPongHandler(func(data string) error {
		fmt.Println(data)
		return nil
	})

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				break
			}

			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Println("SEND PING")
				err := conn.WriteMessage(websocket.PingMessage, []byte("PING"))
				if err != nil {
					fmt.Println("Error sending ping:", err)
					return
				}
			}
		}
	}()

	select {}
}

func main() {
	http.HandleFunc("/", echoHandler)
	fmt.Println("Server started on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
