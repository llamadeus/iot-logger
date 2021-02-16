// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package types

import (
	"time"
)

type Message struct {
	// The id of the message.
	ID string `json:"id"`
	// The message content.
	Text string `json:"text"`
	// Timestamp when the message was recorded.
	Timestamp time.Time `json:"timestamp"`
}
