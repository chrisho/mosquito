package alikafka

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"errors"
	"github.com/chrisho/mosquito/helper"
	"time"
)

// NewConsumer
func (s *AliKafka) NewConsumer(cfg *KafkaConfig) (err error) {
	// 消息队列配置
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
func (s *AliKafka) ProcessMessage(cfg *KafkaConfig, msgFunc []func(msg *sarama.ConsumerMessage)) (err error) {
	// 没有监听
	if len(msgFunc) == 0 {
		return err
	}
	// 初始化 消费者
	if s.Consumer == nil {
		if err = s.NewConsumer(cfg); err != nil {
			return err
		}
	}
	fmt.Println("consumer is start...")
	// 开始消费
	for {
		select {
		case msg, ok := <-s.Consumer.Messages():
			if ok {
				for _, f := range msgFunc {
					f(msg)
				}
			}
		case err, ok := <-s.Consumer.Errors():
			if ok {
				helper.GrpcError(500, "Kafka consumer error: "+err.Error())
				// 错误：重连kafka队列
				//return s.reconnectionConsumer(cfg, msgFunc)
			}
		case ntf, ok := <-s.Consumer.Notifications():
			if ok {
				fmt.Printf("kafka consumer rebalance: %v \n", ntf)
			}
		}
	}
	return err
}

// 重连kafka队列
func (s *AliKafka) reconnectionConsumer(cfg *KafkaConfig, msgFunc []func(msg *sarama.ConsumerMessage)) (err error) {
	if err := s.NewConsumer(cfg); err != nil {
		time.Sleep(3 * time.Second)
		s.reconnectionConsumer(cfg, msgFunc)
	}
	return s.ProcessMessage(cfg, msgFunc)
}
