package alikafka

import (
	"github.com/chrisho/mosquito/helper"
	"strings"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

// 配置
type KafkaConfig struct {
	Servers    []string // xxx:8080
	Topics     []string // Topic
	Ak         string   // AccessKey
	Password   string   // SecretKey的后10位
	ConsumerId string   // ConsumerID
	CertFile   string   // 根证书路径
	ConfigPath string   // 配置文件目录
}

type AliKafka struct {
	Producer sarama.SyncProducer
	Consumer *cluster.Consumer
}

var configPath = "config/"
var producerConnectPanic bool
var consumerConnectPanic bool

func init() {
	// connect with panic : default false
	producerConnectPanic = strings.ToLower(helper.GetEnv("KafkaProducerConnectPanic")) == "true"
	consumerConnectPanic = strings.ToLower(helper.GetEnv("KafkaConsumerConnectPanic")) == "true"
}

// NewAliKafka
func NewAliKafka() *AliKafka {
	return new(AliKafka)
}

// 初始化配置
//kafkaConfig = &KafkaConfig{
//	Servers:    strings.Split(helper.GetEnv("KafkaServers"), ","),
//	Topics:     strings.Split(helper.GetEnv("KafkaTopics"), ","),
//	Ak:         helper.GetEnv("KafkaAccessKey"),
//	Password:   helper.GetEnv("KafkaPassword"),
//	ConsumerId: helper.GetEnv("KafkaConsumerId"),
//	CertFile:   configPath + helper.GetEnv("KafkaCertFile"),
//}
func NewConfig() *KafkaConfig {
	return &KafkaConfig{
		ConfigPath: configPath,
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
