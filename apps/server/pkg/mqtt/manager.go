package mqtt

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DeviceManager 设备管理器：管理模拟设备的生命周期
type DeviceManager struct {
	devices  map[string]*SimulatedDevice
	mu       sync.Mutex
	broker   string
}

// NewDeviceManager 创建设备管理器
func NewDeviceManager(broker string) *DeviceManager {
	return &DeviceManager{
		devices: make(map[string]*SimulatedDevice),
		broker:  broker,
	}
}

// AddDevice 添加模拟设备
func (m *DeviceManager) AddDevice(deviceType, name string) (*SimulatedDevice, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	deviceID := fmt.Sprintf("device-%s", uuid.New().String()[:8])
	if name == "" {
		name = deviceID
	}

	dev := NewSimulatedDevice(deviceID, deviceType, name, 5*time.Second)
	if err := dev.Connect(m.broker); err != nil {
		return nil, fmt.Errorf("启动设备失败: %w", err)
	}

	m.devices[deviceID] = dev
	zap.L().Info("模拟设备已添加",
		zap.String("device_id", deviceID),
		zap.String("type", deviceType),
		zap.String("name", name),
	)
	return dev, nil
}

// RemoveDevice 移除设备
func (m *DeviceManager) RemoveDevice(deviceID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dev, ok := m.devices[deviceID]
	if !ok {
		return fmt.Errorf("设备 %s 不存在", deviceID)
	}

	dev.Disconnect()
	delete(m.devices, deviceID)
	zap.L().Info("模拟设备已移除", zap.String("device_id", deviceID))
	return nil
}

// GetDevice 获取设备
func (m *DeviceManager) GetDevice(deviceID string) (*SimulatedDevice, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	dev, ok := m.devices[deviceID]
	return dev, ok
}

// ListDevices 列出所有设备
func (m *DeviceManager) ListDevices() []*SimulatedDevice {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]*SimulatedDevice, 0, len(m.devices))
	for _, dev := range m.devices {
		result = append(result, dev)
	}
	return result
}

// DeviceCount 设备数量
func (m *DeviceManager) DeviceCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.devices)
}

// StopAll 停止所有设备
func (m *DeviceManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, dev := range m.devices {
		dev.Disconnect()
		delete(m.devices, id)
	}
	zap.L().Info("所有模拟设备已停止")
}
