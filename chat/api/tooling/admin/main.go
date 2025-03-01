// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var url = "ws://localhost:3000/v1/connect"

func main() {
	if err := hack(); err != nil {
		log.Fatalf("Error happened: %v", err)
	}
}

func hack() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", url)

	c, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		log.Fatal("error dialing ws: %w", err)
	}
	defer c.Close()

	_, msg, err := c.ReadMessage()
	if err != nil {
		return fmt.Errorf("Error Reading Message: %v", err)
	}

	if string(msg) != "Hello" {
		return fmt.Errorf("Hello message not correct: exp- Hello, got- %v", string(msg))
	}

	user := struct {
		ID   uuid.UUID
		Name string
	}{
		ID:   uuid.New(),
		Name: "Natnael",
	}
	usr, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}
	if err = c.WriteMessage(websocket.TextMessage, usr); err != nil {
		fmt.Errorf("Error Writing Message: %v", err)
	}

	_, msg, err = c.ReadMessage()
	if err != nil {
		return fmt.Errorf("Error Reading Message: %v", err)
	}

	fmt.Println(string(msg))

	return nil
}
