package mqtt

import (
	"context"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// MQTTClient 封装 paho MQTT 客户端
type MQTTClient struct {
	client     mqtt.Client
	cfg        *ClientConfig
	mu         sync.Mutex
	subHandlers map[string]func(topic string, payload []byte)
}

// ClientConfig MQTT 客户端配置
type ClientConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
	Topics   []string
}

// NewMQTTClient 创建 MQTT 客户端
func NewMQTTClient(cfg *ClientConfig) *MQTTClient {
	c := &MQTTClient{
		cfg:         cfg,
		subHandlers: make(map[string]func(topic string, payload []byte)),
	}
	return c
}

// Connect 连接 MQTT Broker
func (c *MQTTClient) Connect(ctx context.Context) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.cfg.Broker)
	opts.SetClientID(c.cfg.ClientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		zap.L().Warn("MQTT 连接断开", zap.Error(err))
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		zap.L().Info("MQTT 已连接", zap.String("broker", c.cfg.Broker))
		// 重连后重新订阅
		c.resubscribe(client)
	})

	if c.cfg.Username != "" {
		opts.SetUsername(c.cfg.Username)
		opts.SetPassword(c.cfg.Password)
	}

	opts.SetConnectTimeout(5 * time.Second)

	c.client = mqtt.NewClient(opts)
	token := c.client.Connect()
	token.Wait()

	select {
	case <-ctx.Done():
		c.client.Disconnect(250)
		return ctx.Err()
	default:
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("MQTT 连接失败: %w", err)
	}

	// 订阅默认 topic
	for _, topic := range c.cfg.Topics {
		if err := c.doSubscribe(topic, 1); err != nil {
			zap.L().Warn("MQTT 订阅失败", zap.String("topic", topic), zap.Error(err))
		}
	}

	return nil
}

// Publish 发布消息
func (c *MQTTClient) Publish(topic string, qos byte, retained bool, payload []byte) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("MQTT 未连接")
	}

	token := c.client.Publish(topic, qos, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("发布失败: %w", err)
	}
	return nil
}

// Subscribe 订阅 topic（带回调）
func (c *MQTTClient) Subscribe(topic string, qos byte, handler func(topic string, payload []byte)) error {
	c.mu.Lock()
	c.subHandlers[topic] = handler
	c.mu.Unlock()

	if c.client != nil && c.client.IsConnected() {
		return c.doSubscribe(topic, qos)
	}
	return nil
}

// IsConnected 检查连接状态
func (c *MQTTClient) IsConnected() bool {
	return c.client != nil && c.client.IsConnected()
}

// Disconnect 断开连接
func (c *MQTTClient) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(500)
	}
}

func (c *MQTTClient) doSubscribe(topic string, qos byte) error {
	token := c.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		c.mu.Lock()
		handler, ok := c.subHandlers[msg.Topic()]
		// 也尝试匹配通配符 topic 的 handler
		if !ok {
			for t, h := range c.subHandlers {
				if matchTopic(t, msg.Topic()) {
					handler = h
					break
				}
			}
		}
		c.mu.Unlock()

		if handler != nil {
			handler(msg.Topic(), msg.Payload())
		}
	})
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("订阅 %s 失败: %w", topic, err)
	}
	zap.L().Info("MQTT 已订阅", zap.String("topic", topic))
	return nil
}

func (c *MQTTClient) resubscribe(client mqtt.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for topic := range c.subHandlers {
		token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
			if handler, ok := c.subHandlers[msg.Topic()]; ok && handler != nil {
				handler(msg.Topic(), msg.Payload())
			}
		})
		token.Wait()
	}
}

// matchTopic 简单通配符匹配（支持 # 匹配所有，支持 + 匹配单层）
func matchTopic(pattern, topic string) bool {
	if pattern == "#" || pattern == topic {
		return true
	}
	// sensors/# 匹配 sensors/ 开头的所有 topic
	if len(pattern) > 2 && pattern[len(pattern)-2:] == "/#" {
		prefix := pattern[:len(pattern)-1] // 去掉 # 保留 /
		return len(topic) >= len(prefix) && topic[:len(prefix)] == prefix
	}
	return false
}
