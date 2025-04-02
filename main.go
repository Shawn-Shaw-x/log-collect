package main

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"logtransfer/config"
	"logtransfer/es"
	"logtransfer/kafka"
)

// consumer message from kafka
// transfer message to elasticsearch
func main() {
	// init config
	newConfig, err := config.NewConfig()
	if err != nil {
		logrus.Error("load config file failed.", err.Error())
		return
	}
	logrus.Info("load config file succeed.", newConfig)

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Error("load config file failed.", err.Error())
		return
	}

	// core message channel
	// kafka --> msgChan --> elasticsearch
	msgChan := make(chan *sarama.ConsumerMessage, cfg.ESConfig.MaxChanSize)

	// init kafka consumer
	kfkClient, err := kafka.NewKFKClient(cfg)
	defer kfkClient.Consumer.Close()
	if err != nil {
		logrus.Error("create kafka client failed.", err.Error())
		return
	}
	// transfer messages from kafka to msgChan
	kfkClient.AsyncReadMessageToChan(cfg.KafkaConfig.Topic, msgChan)

	//  init elasticsearch
	esClient, err := es.NewESClient(cfg)
	if err != nil {
		logrus.Error("create ES client failed.", err.Error())
		return
	}
	// transfer messages form msgChan to elasticsearch
	esClient.SendMsg2ESBatch(cfg.ESConfig, msgChan)

	//block main goroutine
	select {}

}
