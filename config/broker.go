package config

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Address []string
}

func InitKafkaProducer(cfg KafkaConfig) (*sarama.SyncProducer, error) {
	log.Println("masuk")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Tunggu ACK dari semua broker
	config.Producer.Retry.Max = 5                    // Retry maksimal
	config.Producer.Return.Successes = true          // Perlu untuk SyncProducer

	producer, err := sarama.NewSyncProducer(cfg.Address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start kafka producer: %v", err)
	}

	return &producer, nil
}
