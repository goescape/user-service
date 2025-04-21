package config

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Broker              string
	Username            string
	Password            string
	ConsumerGroupPrefix string
	SchemaRegistry      string
}

func InitKafkaProducer(kafkaCfg KafkaConfig) (*sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5

	if kafkaCfg.Username != "" && kafkaCfg.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kafkaCfg.Username
		config.Net.SASL.Password = kafkaCfg.Password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Handshake = true
		config.Net.TLS.Enable = false
	}

	producer, err := sarama.NewSyncProducer([]string{kafkaCfg.Broker}, config)
	if err != nil {
		panic(err)
	}

	log.Printf("[Success] - Connected to kafka broker")
	return &producer, nil
}
