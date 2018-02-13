package alikafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"errors"
	"time"
)

// configs : 可选；(默认配置)
func (s *AliKafka) NewConsumer(configs ...*KafkaConfig) (err error) {
	// cluster.config
	cfg := kafkaConfig
	if len(configs) > 0 {
		cfg = configs[0]
	}
	clusterCfg, err := s.initConsumerConfig(cfg)
	if err != nil {
		msg := err.Error()
		fmt.Println(msg)
		s.consumerConnectWithPanic(msg) // panic
		return err
	}
	// 消息消费者
	s.Consumer, err = cluster.NewConsumer(cfg.Servers, cfg.ConsumerId, cfg.Topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		fmt.Println(msg)
		s.consumerConnectWithPanic(msg) // panic
		return err
	}
	return err
}

// 消息队列配置
func (s *AliKafka) initConsumerConfig(cfg *KafkaConfig) (clusterCfg *cluster.Config, err error) {
	// cluster.config
	clusterCfg = cluster.NewConfig()
	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = cfg.Ak
	clusterCfg.Net.SASL.Password = cfg.Password
	clusterCfg.Net.SASL.Handshake = true
	// 根证书
	certBytes, err := ioutil.ReadFile(cfg.CertFile)
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		msg := "kafka producer failed to parse root certificate"
		return clusterCfg, errors.New(msg)
	}
	// tls
	clusterCfg.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	clusterCfg.Net.TLS.Enable = true
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	clusterCfg.Group.Return.Notifications = true
	clusterCfg.Version = sarama.V0_10_0_0 // 版本
	// 验证配置
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		return clusterCfg, errors.New(msg)
	}
	return clusterCfg, err
}

// configs : 可选
func (s *AliKafka) ProcessMessage(msgFunc []func(msgKey, msgValue []byte), configs ...*KafkaConfig) (err error) {
	if s.Consumer == nil {
		return errors.New("consumer is nil pointer, please run NewConsumer")
	}
	if len(msgFunc) == 0 {
		return err
	}
	fmt.Println("consumer is start...")
	// 开始消费
	for {
		select {
		case msg, ok := <-s.Consumer.Messages():
			if ok {
				for _, f := range msgFunc {
					f(msg.Key, msg.Value)
				}
				s.Consumer.MarkOffset(msg, "") // mark message as processed
			}
		case err, ok := <-s.Consumer.Errors():
			if ok {
				fmt.Printf("Kafka consumer error: %v \n", err.Error())
				time.Sleep(time.Second) // 1s 后重启
				s.NewConsumer(configs...)
				s.ProcessMessage(msgFunc)
			}
		case ntf, ok := <-s.Consumer.Notifications():
			if ok {
				fmt.Printf("kafka consumer rebalance: %v \n", ntf)
			}
		}
	}
	return err
}
