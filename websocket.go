package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func establishWebsocket() {
	wsURLWithToken := fmt.Sprintf(wsURL, authToken)
	conn, _, err := websocket.DefaultDialer.Dial(wsURLWithToken, nil)
	if err != nil {
		fmt.Printf("Error establishing WebSocket connection: %v\n", err)
		return
	}
	defer conn.Close()

	done = make(chan struct{})

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading WebSocket message: %v\n", err)
				close(done)
				return
			}
			handleMessage(message)
		}
	}()
	<-done
}

func handleMessage(message []byte) {
	var msg struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Printf("Error parsing WebSocket message: %v\n", err)
		return
	}

	switch msg.Type {
	case "userState":
		handleUserState(message)
	case "listUpdate":
		handleListUpdate(message)
		registerCourses()
	default:
		fmt.Printf("Received unknown message type: %s\n", msg.Type)
	}
}

func handleUserState(message []byte) {
	var parsedMessage UserStateMessage
	err := json.Unmarshal(message, &parsedMessage)
	if err != nil {
		return
	}

	if len(favoriteCourses) == 0 {
		favoriteCourses = parsedMessage.Message.Favorites
	}

	registrationTime = parsedMessage.Message.RegistrationTime
	fmt.Printf("favoriteCourses: %v\n", favoriteCourses)
}

func handleListUpdate(message []byte) {
	var parsedMessage ListUpdateMessage
	if err := json.Unmarshal(message, &parsedMessage); err != nil {
		fmt.Printf("Error parsing listUpdate message: %v\n", err)
		return
	}

	coursesData = parsedMessage.Message
	fmt.Println("Updated list of courses")

	timeRemaining := time.Until(time.Unix(int64(registrationTime)/1000, 0))
	if timeRemaining >= 0 {
		fmt.Print("\rRegistration will start in ", formatDuration(timeRemaining))
	}

	if onTimeRegistration {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			timeRemaining := time.Until(time.Unix(int64(registrationTime)/1000, 0))
			fmt.Print("\rRegistration will start in ", formatDuration(timeRemaining))
			if timeRemaining <= 0 {
				break
			}
		}
	}
	fmt.Println()
}