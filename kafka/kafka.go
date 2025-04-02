package kafka

import (
	"context"
	"logtransfer/config"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type KFKClient struct {
	Consumer sarama.Consumer
}

// new kfkClient instance
func NewKFKClient(config *config.Config) (*KFKClient, error) {
	// init all consumers
	consumer, err := sarama.NewConsumer([]string{config.KafkaConfig.Address}, nil)
	if err != nil {
		logrus.Errorf("NewConsumer err: %v", err)
		panic(err)
	}
	//  return
	return &KFKClient{
		Consumer: consumer,
	}, nil
}

// transfer messages from kafka to msgChan
func (k *KFKClient) AsyncReadMessageToChan(topic string, msgChan chan<- *sarama.ConsumerMessage) error {
	//  get patition list of this topic
	partitionList, err := k.Consumer.Partitions(topic)
	if err != nil {
		logrus.Errorf("Error getting list of partitions: %v", err)
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"count": len(partitionList),
		"value": partitionList,
	}).Info("get list of partitions")

	ctx, cancel := context.WithCancel(context.Background())

	for _, partition := range partitionList {
		partitionConsumer, err := k.Consumer.ConsumePartition("web_log", partition, sarama.OffsetNewest)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"partition": partition,
				"error":     err,
			}).Error("Failed to consume partition")
			continue
		}
		// for a single partition,
		// execute goroutine to transfer message to msgChan
		go func(pc sarama.PartitionConsumer, p int32) {
			defer pc.AsyncClose()
			for {
				select {
				case msg := <-pc.Messages():
					logrus.WithFields(logrus.Fields{
						"partition": p,
						"message":   string(msg.Value),
					}).Info("Consume messages")
					msgChan <- msg
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
	return nil
}
