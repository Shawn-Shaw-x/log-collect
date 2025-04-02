package test

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"testing"
)

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJson(t *testing.T) {
	marshal, err := json.Marshal(People{"John", 20})
	if err != nil {
		t.Error(err)
	}
	logrus.Info(string(marshal))
}
