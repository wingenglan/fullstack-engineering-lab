package service

import (
	"encoding/json"
	"fullstack-engineering-lab/server/internal/model"
	mqttPkg "fullstack-engineering-lab/server/pkg/mqtt"
	"time"

	"go.uber.org/zap"
)

// MQTTService MQTT 业务逻辑
type MQTTService struct {
	client  *mqttPkg.MQTTClient
	broker  *mqttPkg.SSEBroker
	manager *mqttPkg.DeviceManager
	registry *mqttPkg.DeviceRegistry
}

// NewMQTTService 创建 MQTT 服务
func NewMQTTService(client *mqttPkg.MQTTClient, brokerAddr string) *MQTTService {
	broker := mqttPkg.NewSSEBroker(200)
	registry := mqttPkg.NewDeviceRegistry()
	manager := mqttPkg.NewDeviceManager(brokerAddr)

	svc := &MQTTService{
		client:   client,
		broker:   broker,
		manager:  manager,
		registry: registry,
	}

	if client != nil {
		// 订阅设备注册
		client.Subscribe(mqttPkg.TopicDeviceRegister, 1, func(topic string, payload []byte) {
			var reg mqttPkg.DeviceRegisterPayload
			if err := json.Unmarshal(payload, &reg); err != nil {
				return
			}
			svc.registry.Register(reg)
			svc.broker.Publish(mqttPkg.SSEEvent{
				Type:      "device_online",
				DeviceID:  reg.DeviceID,
				Payload:   string(payload),
				Timestamp: time.Now().Format(time.RFC3339),
			})
		})

		// 订阅设备状态
		client.Subscribe(mqttPkg.TopicDeviceStatus, 1, func(topic string, payload []byte) {
			var status mqttPkg.DeviceStatusPayload
			if err := json.Unmarshal(payload, &status); err != nil {
				return
			}
			svc.registry.UpdateStatus(status.DeviceID, status.Online)
			eventType := "device_online"
			if !status.Online {
				eventType = "device_offline"
			}
			svc.broker.Publish(mqttPkg.SSEEvent{
				Type:      eventType,
				DeviceID:  status.DeviceID,
				Payload:   string(payload),
				Timestamp: time.Now().Format(time.RFC3339),
			})
		})

		// 订阅设备遥测
		client.Subscribe(mqttPkg.TopicDeviceTelemetry, 1, func(topic string, payload []byte) {
			var telemetry mqttPkg.DeviceTelemetryPayload
			if err := json.Unmarshal(payload, &telemetry); err != nil {
				return
			}
			svc.registry.UpdateTelemetry(telemetry.DeviceID, telemetry.Data)
			svc.broker.Publish(mqttPkg.SSEEvent{
				Type:      "telemetry",
				DeviceID:  telemetry.DeviceID,
				Payload:   string(payload),
				Timestamp: time.Now().Format(time.RFC3339),
			})
		})

		// 订阅设备指令响应
		client.Subscribe(mqttPkg.TopicDeviceResponse, 1, func(topic string, payload []byte) {
			svc.broker.Publish(mqttPkg.SSEEvent{
				Type:      "command_response",
				Payload:   string(payload),
				Timestamp: time.Now().Format(time.RFC3339),
			})
		})

		zap.L().Info("MQTT 设备管理服务已初始化")
	}

	return svc
}

// --- 设备管理 ---

// GetDevices 获取设备列表
func (s *MQTTService) GetDevices() []*mqttPkg.DeviceInfo {
	return s.registry.List()
}

// GetDevice 获取设备详情
func (s *MQTTService) GetDevice(deviceID string) (*mqttPkg.DeviceInfo, error) {
	dev, ok := s.registry.Get(deviceID)
	if !ok {
		return nil, errDeviceNotFound
	}
	return dev, nil
}

// AddDevice 添加模拟设备
func (s *MQTTService) AddDevice(req model.AddDeviceRequest) (*mqttPkg.DeviceInfo, error) {
	if s.client == nil || !s.client.IsConnected() {
		return nil, errMQTTNotAvailable
	}

	deviceType := req.Type
	if deviceType == "" {
		deviceType = "temperature_sensor"
	}

	dev, err := s.manager.AddDevice(deviceType, req.Name)
	if err != nil {
		return nil, err
	}

	// 返回注册表中的设备信息（等待注册消息到达）
	time.Sleep(200 * time.Millisecond)
	info, ok := s.registry.Get(dev.DeviceID)
	if !ok {
		info = &mqttPkg.DeviceInfo{
			DeviceID: dev.DeviceID,
			Type:     deviceType,
			Name:     req.Name,
			Online:   false,
		}
	}
	return info, nil
}

// RemoveDevice 移除模拟设备
func (s *MQTTService) RemoveDevice(deviceID string) error {
	if err := s.manager.RemoveDevice(deviceID); err != nil {
		return err
	}
	s.registry.MarkOffline(deviceID)

	s.broker.Publish(mqttPkg.SSEEvent{
		Type:      "device_offline",
		DeviceID:  deviceID,
		Timestamp: time.Now().Format(time.RFC3339),
	})
	return nil
}

// SendDeviceCommand 向设备下发指令
func (s *MQTTService) SendDeviceCommand(deviceID string, req model.DeviceCommandRequest) error {
	if s.client == nil || !s.client.IsConnected() {
		return errMQTTNotAvailable
	}

	cmdPayload, _ := json.Marshal(mqttPkg.DeviceCommandPayload{
		Command: req.Command,
		Params:  req.Params,
	})

	return s.client.Publish(mqttPkg.DeviceCommandTopic(deviceID), 1, false, cmdPayload)
}

// --- SSE 订阅 ---

// Subscribe 订阅全局 SSE 事件流
func (s *MQTTService) Subscribe(clientID string) (<-chan mqttPkg.SSEEvent, func()) {
	return s.broker.Subscribe(clientID)
}

// GetMessages 获取历史消息
func (s *MQTTService) GetMessages(limit int) []mqttPkg.SSEEvent {
	history := s.broker.History()
	if limit > 0 && limit < len(history) {
		return history[len(history)-limit:]
	}
	return history
}

// Publish 手动发布消息（用于兼容旧的 /publish 端点）
func (s *MQTTService) Publish(topic string, payload string, qos byte) error {
	if s.client == nil || !s.client.IsConnected() {
		return errMQTTNotAvailable
	}
	return s.client.Publish(topic, qos, false, []byte(payload))
}

// --- 清理 ---

// Shutdown 关闭服务
func (s *MQTTService) Shutdown() {
	s.manager.StopAll()
}

var (
	errMQTTNotAvailable = &MQTTError{Message: "MQTT 服务未连接，请检查 Broker 状态"}
	errDeviceNotFound   = &MQTTError{Message: "设备不存在"}
)

type MQTTError struct {
	Message string
}

func (e *MQTTError) Error() string {
	return e.Message
}
