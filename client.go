package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	urlFlag := flag.String("url", "", "WebSocket server URL")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Please provide WebSocket server URL using -url flag")
		return
	}

	url := *urlFlag

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	conn.SetPingHandler(func(data string) error {
		fmt.Println(data)
		err := conn.WriteMessage(websocket.PongMessage, []byte("PONG"))
		if err != nil {
			fmt.Println("Error sending pong:", err)
			return err
		}
		return nil
	})

	done := make(chan struct{})

	go func() {
		defer conn.Close()
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			fmt.Println("Received message from server:", string(message))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection.")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error sending close message:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
