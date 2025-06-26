package kafkautil

import (
	"context"
	"fmt"
	"testing"
)

var ctx = context.Background()

func TestKafkaDial(t *testing.T) {
	if err := dial(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTopic(t *testing.T) {
	topic := "seckill"
	if err := createTopic(topic); err != nil {
		t.Fatal(err)
	}
}

func TestIsController(t *testing.T) {
	fmt.Println(isThisBrokerController())
}

func TestListTopic(t *testing.T) {
	if err := listTopics(); err != nil {
		t.Fatal(err)
	}
}

func TestListPartition(t *testing.T) {
	topic := "seckill"
	if err := listPartitions(topic); err != nil {
		t.Fatal(err)
	}
}

func TestProduceMessage(t *testing.T) {
	topic := "seckill"
	if err := produceMessage(topic); err != nil {
		t.Fatal(err)
	}
}

func TestConsumeMessage(t *testing.T) {
	topic := "seckill"
	if err := consumeMessage(topic); err != nil {
		t.Fatal(err)
	}
}
