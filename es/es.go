package es

import (
	"bytes"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"log"
	"logtransfer/config"
)

type ESClient struct {
	ESProducer *elasticsearch.Client
}

// transfer messages form msgChan to elasticsearch
func (c *ESClient) SendMsg2ESBatch(esConfig config.ESConfig, msgChan chan *sarama.ConsumerMessage) {
	ctx, cancel := context.WithCancel(context.Background())
	// execute 100 goroutines to send request to elasticsearch
	for i := 0; i < esConfig.GoroutineSize; i++ {
		go func(goroutineSize int) {
			for {
				select {
				case msg := <-msgChan:
					// send to es
					ret, err := c.ESProducer.Index(
						esConfig.Index,
						bytes.NewReader(msg.Value),
					)
					if err != nil {
						logrus.Errorf("send msg to es fail")
					}
					retBytes, _ := json.Marshal(ret)
					logrus.WithFields(logrus.Fields{
						"res": string(retBytes),
					}).Info("index")
					defer ret.Body.Close()
				case <-ctx.Done():
					cancel()
					logrus.Error("context cancel, ES client exit")
					return
				}
			}
		}(i)

	}

}

// create elastic client
func NewESClient(cfg *config.Config) (*ESClient, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.ESConfig.Address,
		},
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// new esclient
	return &ESClient{
		ESProducer: es,
	}, nil
}
