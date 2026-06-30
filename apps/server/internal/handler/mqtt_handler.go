package handler

import (
	"fmt"
	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"
	mqttPkg "fullstack-engineering-lab/server/pkg/mqtt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MQTTHandler MQTT HTTP 处理器
type MQTTHandler struct {
	svc *service.MQTTService
}

// NewMQTTHandler 创建 MQTT Handler
func NewMQTTHandler(svc *service.MQTTService) *MQTTHandler {
	return &MQTTHandler{svc: svc}
}

// --- 设备管理 ---

// GetDevices 获取设备列表
func (h *MQTTHandler) GetDevices(c *gin.Context) {
	devices := h.svc.GetDevices()
	response.Success(c, devices)
}

// GetDevice 获取设备详情
func (h *MQTTHandler) GetDevice(c *gin.Context) {
	deviceID := c.Param("id")
	dev, err := h.svc.GetDevice(deviceID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}
	response.Success(c, dev)
}

// AddDevice 添加模拟设备
func (h *MQTTHandler) AddDevice(c *gin.Context) {
	var req model.AddDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Type = "temperature_sensor"
	}

	dev, err := h.svc.AddDevice(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, dev)
}

// RemoveDevice 移除设备
func (h *MQTTHandler) RemoveDevice(c *gin.Context) {
	deviceID := c.Param("id")
	if err := h.svc.RemoveDevice(deviceID); err != nil {
		response.Error(c, http.StatusNotFound, 40400, err.Error())
		return
	}
	response.Success(c, gin.H{"device_id": deviceID, "removed": true})
}

// SendDeviceCommand 向设备下发指令
func (h *MQTTHandler) SendDeviceCommand(c *gin.Context) {
	deviceID := c.Param("id")
	var req model.DeviceCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.SendDeviceCommand(deviceID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	response.Success(c, gin.H{
		"device_id": deviceID,
		"command":   req.Command,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// --- SSE 订阅 ---

// Subscribe 全局 SSE 事件流
func (h *MQTTHandler) Subscribe(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = uuid.New().String()
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch, unsub := h.svc.Subscribe(clientID)
	defer unsub()

	c.Stream(func(w io.Writer) bool {
		c.SSEvent("connected", gin.H{
			"client_id": clientID,
		})

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-ch:
				if !ok {
					return false
				}
				c.SSEvent("mqtt_event", event)
			case <-ticker.C:
				c.SSEvent("heartbeat", gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
				})
			case <-c.Request.Context().Done():
				return false
			}
		}
	})
}

// SubscribeDevice 订阅单个设备的 SSE 事件流
func (h *MQTTHandler) SubscribeDevice(c *gin.Context) {
	deviceID := c.Param("id")
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = uuid.New().String()
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch, unsub := h.svc.Subscribe(clientID)
	defer unsub()

	c.Stream(func(w io.Writer) bool {
		c.SSEvent("connected", gin.H{
			"client_id": clientID,
			"device_id": deviceID,
		})

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-ch:
				if !ok {
					return false
				}
				if event.DeviceID == "" || event.DeviceID == deviceID {
					c.SSEvent("mqtt_event", event)
				}
			case <-ticker.C:
				c.SSEvent("heartbeat", gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
				})
			case <-c.Request.Context().Done():
				return false
			}
		}
	})
}

// --- 消息 ---

// GetMessages 获取历史消息
func (h *MQTTHandler) GetMessages(c *gin.Context) {
	limit := 50
	limitStr := c.Query("limit")
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
		if limit > 200 {
			limit = 200
		}
	}

	messages := h.svc.GetMessages(limit)
	if messages == nil {
		messages = []mqttPkg.SSEEvent{}
	}

	response.Success(c, messages)
}

// Publish 手动发布消息
func (h *MQTTHandler) Publish(c *gin.Context) {
	var req model.MQTTPublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.Publish(req.Topic, req.Payload, req.QoS); err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}

	response.Success(c, gin.H{
		"topic":     req.Topic,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
