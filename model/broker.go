package model

type KafkaPublish struct {
	Topic string `json:"topic"`
	Key   string `json:"key"`
	Value []byte `json:"byte"`
}
