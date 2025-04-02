package test

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_kfk_producer(test *testing.T) {
	fmt.Println("Kafka Client")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll          // wait for all partition
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // random partition
	saramaConfig.Producer.Return.Successes = true
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, saramaConfig)
	defer client.Close()
	if err != nil {
		panic(err)
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = "test"
	msg.Value = sarama.StringEncoder("this is a test")
	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(partition, offset)

}

func Test_kfk_consumer(test *testing.T) {
	logrus.Info("Kafka consumer start")

	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	defer consumer.Close()
	if err != nil {
		logrus.Errorf("NewConsumer err: %v", err)
		panic(err)
	}

	// 获取所有分区
	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		logrus.Errorf("Error getting list of partitions: %v", err)
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"count": len(partitionList),
		"value": partitionList,
	}).Info("get list of partitions")

	ctx, cancel := context.WithCancel(context.Background())

	for _, partition := range partitionList { // ✅ 修正遍历方式
		partitionConsumer, err := consumer.ConsumePartition("web_log", partition, sarama.OffsetNewest)
		defer partitionConsumer.AsyncClose()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"partition": partition,
				"error":     err,
			}).Error("Failed to consume partition")
			continue
		}

		go func(pc sarama.PartitionConsumer, p int32) {
			defer pc.AsyncClose()
			for {
				select {
				case msg := <-pc.Messages():
					logrus.WithFields(logrus.Fields{
						"partition": p,
						"message":   string(msg.Value),
					}).Info("Consume messages")
				case err := <-pc.Errors():
					logrus.Error(err)
					cancel()
					return
				case <-ctx.Done():
					logrus.Errorf("Got ctx.Done signal, shutting down")
					return
				}
			}
		}(partitionConsumer, partition)
	}
	make(chan struct{}) <- struct{}{}

}
