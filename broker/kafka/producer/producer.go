package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// KafkaProducerInterface adalah interface yang mendefinisikan method yang harus ada pada Kafka producer
type KafkaProducerInterface interface {
	SendMessage(topic string, key string, value []byte) error
}

// KafkaProducer adalah struct yang menyimpan instance Kafka producer
type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// NewKafkaProducer menginisiasi Kafka producer berdasarkan producer yang sudah ada
func NewKafkaProducer(p sarama.SyncProducer) KafkaProducerInterface {
	return &KafkaProducer{
		Producer: p,
	}
}

// SendMessage mengirim pesan ke Kafka
func (k *KafkaProducer) SendMessage(topic string, key string, value []byte) error {
	// Membuat message yang akan dikirim ke Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	// Kirim pesan ke Kafka
	partition, offset, err := k.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("Gagal kirim pesan ke Kafka: %v", err)
		return err
	}

	// Log jika pengiriman pesan berhasil
	log.Printf("Pesan terkirim ke topic [%s], partition [%d], offset [%d]", topic, partition, offset)
	return nil
}
