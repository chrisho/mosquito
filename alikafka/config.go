package alikafka

import (
	"github.com/chrisho/mosquito/helper"
	"strings"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

// 配置
type KafkaConfig struct {
	Servers    []string `json:"servers"`    // xxx:8080
	Topics     []string `json:"topics"`     // Topic
	Ak         string   `json:"ak"`         // AccessKey
	Password   string   `json:"password"`   // SecretKey的后10位
	ConsumerId string   `json:"consumerId"` // ConsumerID
	CertFile   string   `json:"cert_file"`  // 根证书路径
}

type AliKafka struct {
	Topics   []string
	Producer sarama.SyncProducer
	Consumer *cluster.Consumer
}

var configPath = "config/"
var kafkaConfig *KafkaConfig
var producerConnectPanic bool
var consumerConnectPanic bool

func init() {
	kafkaConfig = &KafkaConfig{
		Servers:    strings.Split(helper.GetEnv("KafkaServers"), ","),
		Topics:     strings.Split(helper.GetEnv("KafkaTopics"), ","),
		Ak:         helper.GetEnv("KafkaAccessKey"),
		Password:   helper.GetEnv("KafkaPassword"),
		ConsumerId: helper.GetEnv("KafkaConsumerId"),
		CertFile:   configPath + helper.GetEnv("KafkaCertFile"),
	}
	// connect with panic : default false
	producerConnectPanic = strings.ToLower(helper.GetEnv("KafkaProducerConnectPanic")) == "true"
	consumerConnectPanic = strings.ToLower(helper.GetEnv("KafkaConsumerConnectPanic")) == "true"
}

// 获取ali-kafka
func GetAliKafka() *AliKafka {
	return &AliKafka{
		Topics: kafkaConfig.Topics,
	}
}

// producer 链接错误，引发panic
func (s *AliKafka) producerConnectWithPanic(msg string) {
	if producerConnectPanic {
		panic(msg)
	}
}

// consumer 链接错误，引发panic
func (s *AliKafka) consumerConnectWithPanic(msg string) {
	if consumerConnectPanic {
		panic(msg)
	}
}
