package kafka

import (
	"log"
	"user-svc/config"

	"github.com/IBM/sarama"
)

// KafkaProducer adalah struct yang menyimpan instance Kafka producer
type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// NewKafkaProducer menginisiasi Kafka producer berdasarkan config dan mengembalikan koneksi Kafka
func NewKafkaProducer(cfg config.KafkaConfig) (*sarama.SyncProducer, error) {
	// Set konfigurasi untuk Kafka
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Retry.Max = 5

	// Cek jika menggunakan SASL authentication
	if cfg.Username != "" && cfg.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = cfg.Username
		saramaConfig.Net.SASL.Password = cfg.Password
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		saramaConfig.Net.SASL.Handshake = true
		saramaConfig.Net.TLS.Enable = false
	}

	// Buat Kafka producer
	producer, err := sarama.NewSyncProducer([]string{cfg.Broker}, saramaConfig)
	if err != nil {
		log.Printf("Gagal konek ke Kafka broker: %v", err)
		return nil, err
	}

	log.Printf("Kafka producer berhasil terkoneksi ke: %s", cfg.Broker)

	// Mengembalikan pointer ke SyncProducer yang terhubung
	return &producer, nil
}
