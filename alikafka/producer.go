package alikafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"errors"
	"github.com/Shopify/sarama"
)

// configs : 可选；(默认配置)
func (s *AliKafka) NewProducer(cfg *KafkaConfig) (err error) {
	// 消息队列配置
	mqConfig, err := s.initProducerConfig(cfg)
	if err != nil {
		msg := err.Error()
		fmt.Println(msg)
		s.producerConnectWithPanic(msg) // panic
		return err
	}
	// 消息写入者
	s.Producer, err = sarama.NewSyncProducer(cfg.Servers, mqConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		fmt.Println(msg)
		s.producerConnectWithPanic(msg) // panic
		return errors.New(msg)
	}
	return err
}

// 消息队列配置
func (s *AliKafka) initProducerConfig(cfg *KafkaConfig) (mqConfig *sarama.Config, err error) {
	// sarama.config
	mqConfig = sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = cfg.Ak
	mqConfig.Net.SASL.Password = cfg.Password
	mqConfig.Net.SASL.Handshake = true
	// 根证书
	certBytes, err := ioutil.ReadFile(cfg.CertFile)
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	// encoded fail with panic
	if !ok {
		msg := "kafka producer failed to parse root certificate"
		return mqConfig, errors.New(msg)
	}
	// tls
	mqConfig.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	mqConfig.Net.TLS.Enable = true
	mqConfig.Producer.Return.Successes = true
	// 验证配置
	if err = mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", *cfg, err)
		return mqConfig, errors.New(msg)
	}
	return mqConfig, err
}

// 消息数组
func (s *AliKafka) MessageSlice() (r []*sarama.ProducerMessage) {
	return r
}

// 生成消息
func (s *AliKafka) GenerateMessage(topic string, key string, content string) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}
}

// 写入消息
func (s *AliKafka) SendMessage(producer sarama.SyncProducer, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return producer.SendMessage(msg)
}

// 写入消息
func (s *AliKafka) SendMessages(producer sarama.SyncProducer, msg []*sarama.ProducerMessage) error {
	return producer.SendMessages(msg)
}
