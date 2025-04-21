package config

type KafkaConf struct {
	Broker              string
	Username            string
	Password            string
	ConsumerGroupPrefix string
	SchemaRegistry      string
}
