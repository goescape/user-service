package broker

import (
	"fmt"
	"log"
	"net/http"
	"user-svc/helpers/fault"
	"user-svc/model"

	"github.com/IBM/sarama"
)

type producer struct {
	kafka sarama.SyncProducer
}

func NewProducer(kafka sarama.SyncProducer) *producer {
	return &producer{
		kafka: kafka,
	}
}

type KafkaProducer interface {
	SendMessage(payload model.KafkaPublish) error
}

func (p *producer) SendMessage(payload model.KafkaPublish) error {
	partition, offset, err := p.kafka.SendMessage(&sarama.ProducerMessage{
		Topic: payload.Topic,
		Key:   sarama.StringEncoder(payload.Key),
		Value: sarama.ByteEncoder(payload.Value),
	})
	if err != nil {
		return fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed sending message to kafka: %v", err.Error()),
		)
	}

	log.Printf("Success send message, topic [%s], partition [%d], offset [%d]", payload.Topic, partition, offset)
	return nil
}
