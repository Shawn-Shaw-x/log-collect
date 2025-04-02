package test

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestEtcdClient(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("cannot connect to etcd", err)
		return
	}
	defer cli.Close()
	////	 put value
	//cli.Put(ctx, "/go", "go hello")
	//defer cancel()
	//// get
	//response, err := cli.Get(ctx, "/go")
	//if err != nil {
	//	fmt.Println("cannot get key from etcd", err)
	//	return
	//}
	//for index, value := range response.Kvs {
	//	logrus.WithFields(logrus.Fields{
	//		"index": index,
	//		"key":   string(value.Key),
	//		"value": string(value.Value),
	//	}).Info("get from etcd,")
	//}
	//fmt.Println("etcd client success")

	// watch key: go
	watchChan := cli.Watch(context.Background(), "/collect_log_conf")
	for {
		select {
		case resp := <-watchChan:
			for _, ev := range resp.Events {
				logrus.WithFields(logrus.Fields{
					"key":     string(ev.Kv.Key),
					"value":   string(ev.Kv.Value),
					"version": ev.Kv.Version,
					"type":    ev.Type,
				}).Info("watch event")
			}

		case <-ctx.Done():
		}
	}

}
