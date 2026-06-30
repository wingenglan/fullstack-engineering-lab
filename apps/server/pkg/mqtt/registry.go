package mqtt

import (
	"sync"
	"time"
)

// DeviceRegistry 设备注册表（线程安全）
type DeviceRegistry struct {
	devices map[string]*DeviceInfo // deviceID → DeviceInfo
	mu      sync.RWMutex
}

// NewDeviceRegistry 创建设备注册表
func NewDeviceRegistry() *DeviceRegistry {
	return &DeviceRegistry{
		devices: make(map[string]*DeviceInfo),
	}
}

// Register 注册或更新设备
func (r *DeviceRegistry) Register(payload DeviceRegisterPayload) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.devices[payload.DeviceID]; ok {
		// 更新已有设备
		existing.Name = payload.Name
		existing.Type = payload.Type
		existing.LastSeen = now()
		existing.Online = true
	} else {
		r.devices[payload.DeviceID] = &DeviceInfo{
			DeviceID:   payload.DeviceID,
			Type:       payload.Type,
			Name:       payload.Name,
			Online:     true,
			LastSeen:   now(),
			Properties: make(map[string]any),
		}
	}
}

// UpdateStatus 更新设备在线状态
func (r *DeviceRegistry) UpdateStatus(deviceID string, online bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if dev, ok := r.devices[deviceID]; ok {
		dev.Online = online
		dev.LastSeen = now()
	}
}

// UpdateTelemetry 更新设备遥测属性快照
func (r *DeviceRegistry) UpdateTelemetry(deviceID string, data map[string]any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if dev, ok := r.devices[deviceID]; ok {
		dev.LastSeen = now()
		dev.Online = true
		if dev.Properties == nil {
			dev.Properties = make(map[string]any)
		}
		for k, v := range data {
			dev.Properties[k] = v
		}
	}
}

// Get 获取设备信息
func (r *DeviceRegistry) Get(deviceID string) (*DeviceInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dev, ok := r.devices[deviceID]
	return dev, ok
}

// List 列出所有设备
func (r *DeviceRegistry) List() []*DeviceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*DeviceInfo, 0, len(r.devices))
	for _, dev := range r.devices {
		result = append(result, dev)
	}
	return result
}

// Remove 移除设备
func (r *DeviceRegistry) Remove(deviceID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.devices, deviceID)
}

// Count 在线设备数
func (r *DeviceRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, dev := range r.devices {
		if dev.Online {
			count++
		}
	}
	return count
}

// MarkOffline 将指定设备标记为离线
func (r *DeviceRegistry) MarkOffline(deviceID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dev, ok := r.devices[deviceID]; ok {
		dev.Online = false
		dev.LastSeen = now()
	}
}

// CleanupOffline 清理离线超过指定时间的设备
func (r *DeviceRegistry) CleanupOffline(maxAge time.Duration) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-maxAge).Format(time.RFC3339)
	removed := 0
	for id, dev := range r.devices {
		if !dev.Online && dev.LastSeen < cutoff {
			delete(r.devices, id)
			removed++
		}
	}
	return removed
}
