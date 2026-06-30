package mqtt

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// SimulatedDevice 模拟 IoT 设备（独立 MQTT 客户端）
type SimulatedDevice struct {
	DeviceID      string
	deviceType    string // temperature_sensor / humidity_sensor / smart_switch
	name          string
	client        mqtt.Client
	reportInterval time.Duration
	stopCh        chan struct{}
	running       atomic.Bool
	msgCount      atomic.Int64
	mu            sync.Mutex
}

// NewSimulatedDevice 创建模拟设备（不立即连接）
func NewSimulatedDevice(deviceID, deviceType, name string, reportInterval time.Duration) *SimulatedDevice {
	return &SimulatedDevice{
		DeviceID:       deviceID,
		deviceType:     deviceType,
		name:           name,
		reportInterval: reportInterval,
	}
}

// Connect 连接 MQTT Broker 并开始上报
func (d *SimulatedDevice) Connect(broker string) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(d.DeviceID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		zap.L().Info("设备已连接 MQTT Broker", zap.String("device_id", d.DeviceID))
		d.onConnected(client)
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		zap.L().Warn("设备 MQTT 连接断开", zap.String("device_id", d.DeviceID), zap.Error(err))
	})

	d.client = mqtt.NewClient(opts)
	token := d.client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("设备 %s 连接 MQTT 失败: %w", d.DeviceID, err)
	}
	return nil
}

// Disconnect 断开设备连接
func (d *SimulatedDevice) Disconnect() {
	d.Stop()
	if d.client != nil && d.client.IsConnected() {
		// 发布离线状态
		offlinePayload, _ := json.Marshal(DeviceStatusPayload{
			DeviceID: d.DeviceID,
			Online:   false,
		})
		d.client.Publish(DeviceStatusTopic(d.DeviceID), 1, false, offlinePayload)
		d.client.Disconnect(500)
	}
}

// onConnected 连接成功后：注册设备 + 订阅指令 + 启动上报
func (d *SimulatedDevice) onConnected(client mqtt.Client) {
	// 注册设备
	regPayload, _ := json.Marshal(DeviceRegisterPayload{
		DeviceID:     d.DeviceID,
		Type:         d.deviceType,
		Name:         d.name,
		Capabilities: []string{"telemetry", "control"},
	})
	client.Publish(TopicDeviceRegister, 1, false, regPayload)

	// 订阅平台下发指令
	cmdTopic := DeviceCommandTopic(d.DeviceID)
	client.Subscribe(cmdTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
		d.handleCommand(msg.Payload())
	})
	zap.L().Debug("设备已订阅指令 topic", zap.String("device_id", d.DeviceID), zap.String("topic", cmdTopic))

	// 启动定时上报
	d.startReporting(client)
}

// startReporting 定时上报状态和遥测
func (d *SimulatedDevice) startReporting(client mqtt.Client) {
	if !d.running.CompareAndSwap(false, true) {
		return
	}
	d.stopCh = make(chan struct{})

	// 立即上报一次在线状态
	d.reportStatus(client)
	d.reportTelemetry(client)

	go func() {
		ticker := time.NewTicker(d.reportInterval)
		defer ticker.Stop()
		for {
			select {
			case <-d.stopCh:
				return
			case <-ticker.C:
				d.reportStatus(client)
				d.reportTelemetry(client)
				d.msgCount.Add(2)
			}
		}
	}()
}

// Stop 停止上报
func (d *SimulatedDevice) Stop() {
	if d.running.CompareAndSwap(true, false) {
		close(d.stopCh)
	}
}

// IsRunning 是否运行中
func (d *SimulatedDevice) IsRunning() bool {
	return d.running.Load()
}

// MessageCount 已上报消息数
func (d *SimulatedDevice) MessageCount() int64 {
	return d.msgCount.Load()
}

// reportStatus 上报设备在线状态
func (d *SimulatedDevice) reportStatus(client mqtt.Client) {
	payload, _ := json.Marshal(DeviceStatusPayload{
		DeviceID: d.DeviceID,
		Online:   true,
		Uptime:   time.Now().Unix(),
	})
	client.Publish(DeviceStatusTopic(d.DeviceID), 1, false, payload)
}

// reportTelemetry 上报遥测数据
func (d *SimulatedDevice) reportTelemetry(client mqtt.Client) {
	data := d.generateTelemetry()
	payload, _ := json.Marshal(DeviceTelemetryPayload{
		DeviceID: d.DeviceID,
		Data:     data,
	})
	client.Publish(DeviceTelemetryTopic(d.DeviceID), 1, false, payload)
}

// generateTelemetry 根据设备类型生成遥测数据
func (d *SimulatedDevice) generateTelemetry() map[string]any {
	data := make(map[string]any)
	switch d.deviceType {
	case "temperature_sensor":
		data["temperature"] = round(15+rand.Float64()*20, 1)
		data["unit"] = "°C"
	case "humidity_sensor":
		data["humidity"] = round(40+rand.Float64()*40, 1)
		data["unit"] = "%"
	case "smart_switch":
		data["power"] = round(100+rand.Float64()*400, 1)
		data["voltage"] = round(215+rand.Float64()*15, 1)
		data["current"] = round(0.5+rand.Float64()*2, 2)
		data["relay_state"] = d.getSwitchState()
	case "environment_sensor":
		data["temperature"] = round(15+rand.Float64()*20, 1)
		data["humidity"] = round(40+rand.Float64()*40, 1)
		data["pressure"] = round(980+rand.Float64()*40, 1)
	}
	return data
}

func (d *SimulatedDevice) getSwitchState() string {
	if rand.Intn(2) == 0 {
		return "on"
	}
	return "off"
}

// handleCommand 处理平台下发的指令
func (d *SimulatedDevice) handleCommand(payload []byte) {
	var cmd DeviceCommandPayload
	if err := json.Unmarshal(payload, &cmd); err != nil {
		zap.L().Error("解析设备指令失败", zap.String("device_id", d.DeviceID), zap.Error(err))
		return
	}

	zap.L().Info("设备收到指令", zap.String("device_id", d.DeviceID), zap.String("command", cmd.Command))

	var resp DeviceCommandResponsePayload
	resp.DeviceID = d.DeviceID
	resp.Command = cmd.Command
	resp.Success = true

	switch cmd.Command {
	case "set_interval":
		if intervalMs, ok := cmd.Params["interval_ms"].(float64); ok {
			d.reportInterval = time.Duration(intervalMs) * time.Millisecond
			resp.Message = fmt.Sprintf("上报间隔已设为 %dms", int(intervalMs))
			// 重启上报
			d.Stop()
			d.startReporting(d.client)
		}
	case "reboot":
		resp.Message = "设备正在重启..."
		// 模拟断开重连
		go func() {
			time.Sleep(2 * time.Second)
			d.Stop()
			d.startReporting(d.client)
		}()
	case "toggle":
		resp.Message = "开关已切换"
	default:
		resp.Success = false
		resp.Message = fmt.Sprintf("未知指令: %s", cmd.Command)
	}

	respPayload, _ := json.Marshal(resp)
	d.client.Publish(DeviceResponseTopic(d.DeviceID), 1, false, respPayload)
}
