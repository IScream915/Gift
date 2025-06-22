package kafka

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
	if err := createTopic(); err != nil {
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
	if err := listPartitions(); err != nil {
		t.Fatal(err)
	}
}

func TestProduceMessage(t *testing.T) {
	if err := produceMessage(); err != nil {
		t.Fatal(err)
	}
}

func TestConsumeMessage(t *testing.T) {
	if err := consumeMessage(); err != nil {
		t.Fatal(err)
	}
}
