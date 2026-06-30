package mqtt

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Simulator 传感器数据模拟器
type Simulator struct {
	client       *MQTTClient
	intervalMs   int
	ticker       *time.Ticker
	stopCh       chan struct{}
	runningMu    sync.Mutex
	running      bool
	messageCount atomic.Int64
}

// NewSimulator 创建模拟器
func NewSimulator(client *MQTTClient, intervalMs int) *Simulator {
	if intervalMs <= 0 {
		intervalMs = 2000
	}
	return &Simulator{
		client:     client,
		intervalMs: intervalMs,
	}
}

// Start 启动模拟器
func (s *Simulator) Start(intervalMs int) error {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if s.running {
		return fmt.Errorf("模拟器已在运行中")
	}

	if intervalMs > 0 {
		s.intervalMs = intervalMs
	}

	s.stopCh = make(chan struct{})
	s.ticker = time.NewTicker(time.Duration(s.intervalMs) * time.Millisecond)
	s.running = true

	go s.run()

	zap.L().Info("MQTT 模拟器已启动", zap.Int("interval_ms", s.intervalMs))
	return nil
}

// Stop 停止模拟器
func (s *Simulator) Stop() {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if !s.running {
		return
	}

	s.ticker.Stop()
	close(s.stopCh)
	s.running = false
	zap.L().Info("MQTT 模拟器已停止")
}

// IsRunning 是否运行中
func (s *Simulator) IsRunning() bool {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()
	return s.running
}

// MessageCount 已发送消息数
func (s *Simulator) MessageCount() int64 {
	return s.messageCount.Load()
}

// IntervalMs 当前间隔
func (s *Simulator) IntervalMs() int {
	return s.intervalMs
}

func (s *Simulator) run() {
	for {
		select {
		case <-s.stopCh:
			return
		case <-s.ticker.C:
			s.generateAndPublish()
		}
	}
}

func (s *Simulator) generateAndPublish() {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	deviceID := fmt.Sprintf("sensor-%03d", rand.Intn(10)+1)

	// 温度：15~35°C
	temp := SensorData{
		DeviceID:   deviceID,
		Value:      round(15 + rand.Float64()*20, 1),
		Unit:       UnitTemperature,
		SensorType: "temperature",
		Timestamp:  timestamp,
	}
	s.publishSensor(TopicTemperature, temp)

	// 湿度：40~80%
	humid := SensorData{
		DeviceID:   deviceID,
		Value:      round(40+rand.Float64()*40, 1),
		Unit:       UnitHumidity,
		SensorType: "humidity",
		Timestamp:  timestamp,
	}
	s.publishSensor(TopicHumidity, humid)

	// 气压：980~1020hPa
	pressure := SensorData{
		DeviceID:   deviceID,
		Value:      round(980+rand.Float64()*40, 1),
		Unit:       UnitPressure,
		SensorType: "pressure",
		Timestamp:  timestamp,
	}
	s.publishSensor(TopicPressure, pressure)

	// 聚合 all topic
	all := map[string]SensorData{
		"temperature": temp,
		"humidity":    humid,
		"pressure":    pressure,
	}
	allPayload, _ := json.Marshal(all)
	s.client.Publish(TopicAll, 1, false, allPayload)

	s.messageCount.Add(1)
}

func (s *Simulator) publishSensor(topic string, data SensorData) {
	payload, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("序列化传感器数据失败", zap.Error(err))
		return
	}
	if err := s.client.Publish(topic, 1, false, payload); err != nil {
		zap.L().Error("发布传感器数据失败", zap.String("topic", topic), zap.Error(err))
	}
}

func round(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}
