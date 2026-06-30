package mqtt

import "time"

// SensorData 传感器数据
type SensorData struct {
	DeviceID   string  `json:"device_id"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"`
	SensorType string  `json:"sensor_type"`
	Timestamp  string  `json:"timestamp"`
}

// SSEEvent SSE 推送事件（增强版，支持设备事件类型）
type SSEEvent struct {
	Type      string `json:"type"` // device_online / device_offline / telemetry / command_response / status
	DeviceID  string `json:"device_id,omitempty"`
	Topic     string `json:"topic,omitempty"`
	Payload   string `json:"payload,omitempty"`
	Timestamp string `json:"timestamp"`
}

// DeviceInfo 设备信息（注册表中维护）
type DeviceInfo struct {
	DeviceID   string            `json:"device_id"`
	Type       string            `json:"type"` // temperature_sensor / humidity_sensor / smart_switch
	Name       string            `json:"name"`
	Online     bool              `json:"online"`
	LastSeen   string            `json:"last_seen"`
	Properties map[string]any    `json:"properties"` // 最新遥测属性快照
}

// DeviceRegisterPayload 设备注册消息负载
type DeviceRegisterPayload struct {
	DeviceID     string   `json:"device_id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"` // ["telemetry", "control"]
}

// DeviceStatusPayload 设备状态上报负载
type DeviceStatusPayload struct {
	DeviceID string `json:"device_id"`
	Online   bool   `json:"online"`
	Uptime   int64  `json:"uptime_sec"`
}

// DeviceTelemetryPayload 设备遥测负载
type DeviceTelemetryPayload struct {
	DeviceID string         `json:"device_id"`
	Data     map[string]any `json:"data"` // {"temperature": 25.5, "humidity": 60}
}

// DeviceCommandPayload 平台下发指令负载
type DeviceCommandPayload struct {
	Command string         `json:"command"` // set_interval / reboot / toggle / set_threshold
	Params  map[string]any `json:"params,omitempty"`
}

// DeviceCommandResponse 设备指令响应负载
type DeviceCommandResponsePayload struct {
	DeviceID string `json:"device_id"`
	Command  string `json:"command"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
}

// Topic 常量
const (
	// 设备管理 topic
	TopicDeviceRegister = "devices/register"
	TopicDeviceStatus   = "devices/+/status"
	TopicDeviceTelemetry = "devices/+/telemetry"
	TopicDeviceCommand  = "devices/+/commands"
	TopicDeviceResponse = "devices/+/response"

	// 旧版传感器 topic（兼容保留）
	TopicTemperature = "sensors/temperature"
	TopicHumidity    = "sensors/humidity"
	TopicPressure    = "sensors/pressure"
	TopicAll         = "sensors/all"
)

// 传感器单位
var (
	UnitTemperature = "°C"
	UnitHumidity    = "%"
	UnitPressure    = "hPa"
)

// 生成设备 topic
func DeviceStatusTopic(deviceID string) string {
	return "devices/" + deviceID + "/status"
}
func DeviceTelemetryTopic(deviceID string) string {
	return "devices/" + deviceID + "/telemetry"
}
func DeviceCommandTopic(deviceID string) string {
	return "devices/" + deviceID + "/commands"
}
func DeviceResponseTopic(deviceID string) string {
	return "devices/" + deviceID + "/response"
}

func now() string {
	return time.Now().Format(time.RFC3339)
}
