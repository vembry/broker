package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type ActiveMessage struct {
	Id         ksuid.KSUID `json:"id"`
	QueueName  string      `json:"queue_name"`
	PollExpiry time.Time   `json:"poll_expiry"`
	Queue      *Message    `json:"queue"`
}

type IdleQueue struct {
	Messages []*Message `json:"items"`
	// add other info
	// ...
}

type EnqueuePayload struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

type Message struct {
	Payload string
	// add other info
	// ...
}

type QueueData struct {
	ActiveMessageCount int64 `json:"active_task_count"`
	IdleQueueCount     int64 `json:"idle_queue_count"`
}
