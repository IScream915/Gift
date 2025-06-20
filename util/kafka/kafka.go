package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	reader *kafka.Reader
	topic  = "user_click"
)

func WriteKafka(ctx context.Context) {
	writer := kafka.Writer{
		Addr:                   kafka.TCP("198.19.249.225:9092"),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}

	defer writer.Close()

	for i := 0; i < 3; i++ {
		if err := writer.WriteMessages(ctx,
			kafka.Message{Key: []byte("1"), Value: []byte("c")},
			kafka.Message{Key: []byte("2"), Value: []byte("y")},
			kafka.Message{Key: []byte("3"), Value: []byte("x")},
		); err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				fmt.Println("kafka write error:", err)
			}
		} else {
			fmt.Println("kafka write success")
			break
		}
	}
}

func ReadKafka(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"198.19.249.225:9092"},
		Topic:          topic,
		CommitInterval: time.Second,
		GroupID:        "rec_team",
		StartOffset:    kafka.FirstOffset,
	})

	//defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("read message error:", err)
			break
		} else {
			fmt.Println("topic:", msg.Topic, "partition", msg.Partition, "offset", msg.Offset, "key", string(msg.Key), "value", string(msg.Value))
		}
	}
}

func ListenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	fmt.Println("接收到信号:", sig.String())
	if reader != nil {
		reader.Close()
	}
	os.Exit(0)
}
