package test

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"log"
	"testing"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func TestElastic(t *testing.T) {
	// 创建 Elasticsearch 客户端
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // 你的 ES 地址
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	logrus.Info("connect success")
	p1 := Person{
		Name:    "huanghua",
		Age:     20,
		Married: false,
	}
	doc, err := json.Marshal(p1)
	if err != nil {
		log.Fatal(err)
	}
	res, err := es.Index("user", bytes.NewReader(doc))
	if err != nil {
		logrus.Fatal(err)
	}
	defer res.Body.Close()

	logrus.WithFields(logrus.Fields{
		"res": res,
	}).Info("index")
}
