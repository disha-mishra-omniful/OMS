package kafkaa

import (
	"github.com/omniful/go_commons/kafka"
)

type Producer struct {
	*kafka.ProducerClient
}

var ClientInstance *Producer

func Get() *Producer {
	return ClientInstance
}

func Set(client *kafka.ProducerClient) {
	ClientInstance = &Producer{client}
}
