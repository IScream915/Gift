package kafkautil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	addr   string
	config = Config{}
	reader *kafka.Reader
	group  = "test-group"
	topic  = "test-topic"
)

// SecKillMessage 秒杀消息结构
type SecKillMessage struct {
	InventoryId uint64 `json:"inventory_id"` // 商品ID
	UserId      uint64 `json:"user_id"`      // 用户ID
	Count       uint64 `json:"count"`        // 购买数量
	OrderId     string `json:"order_id"`     // 订单ID（用于幂等性）
	CreateTime  int64  `json:"create_time"`  // 创建时间戳
}

type Config struct {
	Address string `mapstructure:"address"`
}

func (c *Config) loadKafkaConfig() error {
	configPath := DefaultConfigPath

	// 使用 Viper 加载配置文件
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	// 获取配置文件内容
	if err := viper.Unmarshal(c); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	addr = c.Address

	return nil
}

// 测试与kafka的连接
func dial() error {
	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

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

// 创建topic，并且限定topic仅有一个partition
func createTopic(topic string) error {
	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

	// 连接到kafka服务端
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
	if err = conn.CreateTopics(topicConfig); err != nil {
		return fmt.Errorf("create topic error: %w", err)
	}
	fmt.Println("✅ Topic 创建成功")
	return nil
}

func isThisBrokerController() (bool, error) {
	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return false, err
	}

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
	fmt.Printf("Controller info: %+v\n", controller)

	// 3. 获取当前连接的远程地址，并与 Controller 地址对比
	remote := conn.RemoteAddr().String()
	fmt.Println("Remote addr:", remote)

	ctrlAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	return remote == ctrlAddr, nil
}

func listTopics() error {

	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

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
		fmt.Println("partition:", partition.ID, "leader:", partition.Leader.Host, ":", partition.Leader.Port)
	}

	return nil
}

func listTopicPartitions(topic string) error {

	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

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

func produceMessage(topic string) error {
	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

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

	ctx := context.Background()

	for i := 0; i < 20; i++ {
		if err := writer.WriteMessages(ctx, msg); err != nil {
			return fmt.Errorf("write message error: %w", err)
		}

		fmt.Println(i, "✅ 消息发送成功")

		time.Sleep(1 * time.Second)
	}

	return nil
}

func consumeMessage(topic string) error {

	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

	reader = kafka.NewReader(kafka.ReaderConfig{
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

func listenSignal() error {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	s, ok := <-done
	if !ok {
		return errors.New("中止信道意外关闭")
	}
	fmt.Println("收到用户操作:", s.String())
	// 判断reader是否存在
	if reader != nil {
		if err := reader.Close(); err != nil {
			return err
		}
	}
	defer os.Exit(0)
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func InitKafkaBroker() error {
	return dial()
}

func InventoryDeduct(ctx context.Context, inventoryId, userId uint64) error {
	return InventoryDeductWithOrder(ctx, inventoryId, userId, 1, "")
}

// InventoryDeductWithOrder 库存扣减消息发送（支持订单信息）
func InventoryDeductWithOrder(ctx context.Context, inventoryId, userId, count uint64, orderId string) error {
	// 读入配置
	if err := config.loadKafkaConfig(); err != nil {
		fmt.Println("load kafka config err: ", err)
		return err
	}

	writer := kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		WriteTimeout:           3 * time.Second,
		RequiredAcks:           kafka.RequireAll,
		Async:                  false, // 改为同步，确保消息发送成功
		AllowAutoTopicCreation: false,
	}

	defer writer.Close()

	// 生成订单ID（如果没有提供）
	if orderId == "" {
		orderId = fmt.Sprintf("order_%d_%d_%d", inventoryId, userId, time.Now().UnixNano())
	}

	// 构造秒杀消息
	seckillMsg := SecKillMessage{
		InventoryId: inventoryId,
		UserId:      userId,
		Count:       count,
		OrderId:     orderId,
		CreateTime:  time.Now().Unix(),
	}

	// 序列化消息
	msgBytes, err := json.Marshal(seckillMsg)
	if err != nil {
		return fmt.Errorf("marshal message error: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(orderId), // 使用订单ID作为key，保证同一订单的消息有序
		Value: msgBytes,
	}

	if err := writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write message error: %w", err)
	}

	fmt.Println("✅ 秒杀消息发送成功, orderId:", orderId)
	return nil
}
