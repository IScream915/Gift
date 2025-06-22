package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"time"
)

var (
	addr        = "198.19.249.172:9092"
	CacheReader *kafka.Reader
	CacheWriter *kafka.Writer
	topic       = "test-topic"
	group       = "test-group"
)

func dial() error {
	conn, err := kafka.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	fmt.Println("本地地址", conn.LocalAddr().String())
	fmt.Println("远端地址", conn.RemoteAddr().String())

	defer conn.Close()

	fmt.Println("✅ Kafka broker 连接成功")
	return nil
}

func createTopic() error {
	conn, err := kafka.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// CreateTopics 仅在请求的 Broker 为 Controller 时才生效
	if err := conn.CreateTopics(topicConfig); err != nil {
		return fmt.Errorf("create topic error: %w", err)
	}
	fmt.Println("✅ Topic 创建成功")
	return nil
}

func isThisBrokerController() (bool, error) {
	// 1. 建立到某个 Broker（可能是任意节点）的连接
	conn, err := kafka.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		return false, fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	// 2. 查询当前集群中的 Controller
	controller, err := conn.Controller()
	if err != nil {
		return false, fmt.Errorf("controller lookup error: %w", err)
	}
	// controller.Host, controller.Port, controller.ID 可用来定位 Controller 节点  [oai_citation:0‡pkg.go.dev](https://pkg.go.dev/gopkg.in/segmentio/kafka-go.v0)

	// 3. 获取当前连接的远程地址，并与 Controller 地址对比
	remote := conn.RemoteAddr().String()

	ctrlAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	return remote == ctrlAddr, nil
}

func listTopics() error {
	conn, err := kafka.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		fmt.Println("kafka connect failed", err)
	}

	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		fmt.Println("kafka read partitions failed", err)
	}

	for _, partition := range partitions {
		fmt.Println(partition.Topic)
	}

	return nil
}

func listPartitions() error {
	conn, err := kafka.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		fmt.Println("kafka connect failed", err)
	}

	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		fmt.Println("kafka read partitions failed", err)
	}
	for _, partition := range partitions {
		fmt.Println("partition:", partition.ID, "leader:", partition.Leader.Host, ":", partition.Leader.Port)
	}

	return nil
}

func produceMessage() error {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireAll,
		Async:                  true,
		AllowAutoTopicCreation: false,
	}

	defer writer.Close()

	msg := kafka.Message{
		Key:   []byte("test"),
		Value: []byte("Hello Kafka-Go!"),
	}

	for i := 0; i < 100; i++ {
		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			return fmt.Errorf("write message error: %w", err)
		}

		fmt.Println(i, "✅ 消息发送成功")

		time.Sleep(1 * time.Second)
	}

	return nil
}

func consumeMessage() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{addr},
		Topic:          topic,
		GroupID:        group,
		CommitInterval: 0,
		StartOffset:    kafka.FirstOffset,
	})
	defer reader.Close()

	ctx := context.Background()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read message error: %w", err)
		}
		fmt.Printf("✅ 消费到消息: key=%s value=%s\n topic=%s partition=%d offset=%d\n",
			string(m.Key), string(m.Value), m.Topic, m.Partition, m.Offset)
	}
}
