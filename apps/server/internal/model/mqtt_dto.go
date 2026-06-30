package model

// MQTTPublishRequest 手动发布消息请求
type MQTTPublishRequest struct {
	Topic   string `json:"topic" binding:"required"`
	Payload string `json:"payload" binding:"required"`
	QoS     byte   `json:"qos"`
}

// MQTTMessageItem 消息记录项
type MQTTMessageItem struct {
	Type      string `json:"type"`
	DeviceID  string `json:"device_id,omitempty"`
	Topic     string `json:"topic,omitempty"`
	Payload   string `json:"payload,omitempty"`
	Timestamp string `json:"timestamp"`
}

// AddDeviceRequest 添加模拟设备请求
type AddDeviceRequest struct {
	Type string `json:"type"` // temperature_sensor / humidity_sensor / smart_switch / environment_sensor
	Name string `json:"name"`
}

// DeviceCommandRequest 设备指令请求
type DeviceCommandRequest struct {
	Command string         `json:"command" binding:"required"`
	Params  map[string]any `json:"params,omitempty"`
}

// DeviceCommandResponse 设备指令响应
type DeviceCommandResponse struct {
	DeviceID string `json:"device_id"`
	Command  string `json:"command"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
}

// DeviceInfoResponse 设备详情响应
type DeviceInfoResponse struct {
	DeviceID   string         `json:"device_id"`
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Online     bool           `json:"online"`
	LastSeen   string         `json:"last_seen"`
	Properties map[string]any `json:"properties"`
}

// SimulatorStartRequest 启动模拟器请求
type SimulatorStartRequest struct {
	IntervalMs int `json:"interval_ms"`
}

// SimulatorStatusResponse 模拟器状态响应
type SimulatorStatusResponse struct {
	Running      bool   `json:"running"`
	IntervalMs   int    `json:"interval_ms"`
	MessageCount int64  `json:"message_count"`
}
