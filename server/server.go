package server

import (
	"broker/model"

	"github.com/segmentio/ksuid"
)

type IBroker interface {
	Get() model.QueueData
	Enqueue(payload model.EnqueuePayload) error
	Poll(queueName string) (*model.ActiveMessage, error)
	CompletePoll(queueId ksuid.KSUID) error
}
