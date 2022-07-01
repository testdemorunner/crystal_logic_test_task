package main

import (
	"fmt"
	"time"
)

// Message keeps message data
type Message struct {
	Timestamp time.Time
	Message   string
}

// NewMessage is a constructor for Message struct
func NewMessage(timestamp time.Time) *Message {
	return &Message{
		Timestamp: timestamp,
		Message:   "Message",
	}
}

func (m Message) String() string {
	return fmt.Sprintf("Timestamp: %v Message: %s", m.Timestamp, m.Message)
}
