<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import mqtt from 'mqtt'
import { getDevices, addDevice, removeDevice, sendDeviceCommand } from '@/api/mqtt'
import GlowCard from '@/components/GlowCard.vue'
import type { MQTTDeviceInfo } from '@/types'

const devices = ref<MQTTDeviceInfo[]>([])
const selectedDeviceId = ref<string>('')
const selectedDevice = ref<MQTTDeviceInfo | null>(null)
const addDialogVisible = ref(false)
const newDevice = ref({ type: 'temperature_sensor', name: '' })
const isLoading = ref(false)

// 指令
const commandDialogVisible = ref(false)
const commandForm = ref({ command: 'set_interval', params: { interval_ms: 5000 } })

// MQTT 客户端（浏览器直连 Mosquitto WebSocket）
let mqttClient: mqtt.MqttClient | null = null

const deviceTypes = [
  { value: 'temperature_sensor', label: '温度传感器' },
  { value: 'humidity_sensor', label: '湿度传感器' },
  { value: 'smart_switch', label: '智能开关' },
  { value: 'environment_sensor', label: '环境传感器' },
]

const commandOptions = [
  { value: 'set_interval', label: '设置上报间隔', hasParams: true },
  { value: 'reboot', label: '重启设备', hasParams: false },
  { value: 'toggle', label: '切换开关', hasParams: false },
]

const typeLabels: Record<string, string> = {
  temperature_sensor: '温度传感器',
  humidity_sensor: '湿度传感器',
  smart_switch: '智能开关',
  environment_sensor: '环境传感器',
}
const typeIcons: Record<string, string> = {
  temperature_sensor: '🌡️',
  humidity_sensor: '💧',
  smart_switch: '⚡',
  environment_sensor: '🌍',
}

function propLabel(k: string, v: any): string {
  const labels: Record<string, string> = {
    temperature: '温度', humidity: '湿度', pressure: '气压',
    power: '功率', voltage: '电压', current: '电流',
    relay_state: '继电器', unit: '单位',
  }
  const units: Record<string, string> = {
    temperature: '°C', humidity: '%', pressure: 'hPa',
    power: 'W', voltage: 'V', current: 'A',
  }
  const unit = units[k] || ''
  return `${labels[k] || k}: ${v} ${unit}`
}

async function loadDevices() {
  try {
    const { data } = await getDevices()
    devices.value = data || []
    if (selectedDeviceId.value) {
      const found = devices.value.find(d => d.device_id === selectedDeviceId.value)
      selectedDevice.value = found || null
      if (!found) selectedDeviceId.value = ''
    }
  } catch { /* 忽略错误 */ }
}

function connectMQTT() {
  mqttClient = mqtt.connect('ws://localhost:9001', {
    clientId: 'browser_' + Math.random().toString(16).slice(2, 10),
    clean: true,
  })

  mqttClient.on('connect', () => {
    // 订阅所有设备的状态和遥测主题（通配符 +）
    mqttClient?.subscribe(['devices/+/status', 'devices/+/telemetry'])
  })

  mqttClient.on('message', (topic, payload) => {
    try {
      const msg = JSON.parse(payload.toString())
      const deviceId: string = msg.device_id
      const device = devices.value.find(d => d.device_id === deviceId)
      if (!device) return

      if (topic.includes('/status')) {
        // 设备上下线状态变更
        device.online = msg.online
        device.last_seen = msg.timestamp
      } else if (topic.includes('/telemetry')) {
        // 遥测数据更新
        device.properties = msg.data
        device.last_seen = msg.timestamp
      }
      // selectedDevice 与 devices 数组中同一对象引用，自动同步
    } catch { /* 忽略解析错误 */ }
  })

  mqttClient.on('error', (err) => {
    console.error('MQTT 连接错误:', err)
  })
}

function disconnectMQTT() {
  if (mqttClient) {
    mqttClient.end()
    mqttClient = null
  }
}

async function handleAddDevice() {
  isLoading.value = true
  try {
    await addDevice(newDevice.value)
    addDialogVisible.value = false
    newDevice.value = { type: 'temperature_sensor', name: '' }
    await loadDevices()
    ElMessage.success('设备已添加')
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '添加失败')
  } finally {
    isLoading.value = false
  }
}

async function handleRemoveDevice(id: string) {
  try {
    await removeDevice(id)
    if (selectedDeviceId.value === id) {
      selectedDeviceId.value = ''
      selectedDevice.value = null
    }
    await loadDevices()
    ElMessage.success('设备已移除')
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '移除失败')
  }
}

function selectDevice(id: string) {
  selectedDeviceId.value = id
  selectedDevice.value = devices.value.find(d => d.device_id === id) || null
}

async function handleSendCommand() {
  if (!selectedDeviceId.value) return
  try {
    await sendDeviceCommand(selectedDeviceId.value, commandForm.value)
    commandDialogVisible.value = false
    ElMessage.success('指令已下发')
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '指令下发失败')
  }
}

function onCommandChange(cmd: string) {
  const opt = commandOptions.find(o => o.value === cmd)
  if (opt?.hasParams) {
    commandForm.value.params = { interval_ms: 5000 }
  } else {
    commandForm.value.params = {}
  }
}

onMounted(() => {
  loadDevices()
  connectMQTT()
})

onUnmounted(() => {
  disconnectMQTT()
})
</script>

<template>
  <div class="mqtt-container">
    <div class="page-header">
      <h2>MQTT IoT 设备管理平台</h2>
      <p class="subtitle">模拟真实 IoT 场景：设备自主上报 + 平台远程控制</p>
    </div>

    <div class="content-layout">
      <!-- 左侧：设备列表 -->
      <div class="left-panel">
        <GlowCard>
          <template #header>
            <div class="panel-header">
              <span>设备列表 ({{ devices.length }})</span>
              <el-button type="primary" size="small" @click="addDialogVisible = true">+ 添加设备</el-button>
            </div>
          </template>

          <div class="device-list">
            <div
              v-for="dev in devices"
              :key="dev.device_id"
              class="device-item"
              :class="{ active: selectedDeviceId === dev.device_id, offline: !dev.online }"
              @click="selectDevice(dev.device_id)"
            >
              <span class="device-icon">{{ typeIcons[dev.type] || '📡' }}</span>
              <div class="device-info">
                <div class="device-name">{{ dev.name || dev.device_id }}</div>
                <div class="device-type">{{ typeLabels[dev.type] || dev.type }}</div>
              </div>
              <el-tag :type="dev.online ? 'success' : 'info'" size="small">
                {{ dev.online ? '在线' : '离线' }}
              </el-tag>
              <el-button
                type="danger"
                size="small"
                circle
                :icon="'Close'"
                @click.stop="handleRemoveDevice(dev.device_id)"
                style="margin-left: 8px"
              />
            </div>

            <el-empty v-if="devices.length === 0" description="暂无设备，点击右上角添加" :image-size="60" />
          </div>
        </GlowCard>
      </div>

      <!-- 右侧：设备详情 + 遥测 + 控制 -->
      <div class="right-panel">
        <template v-if="selectedDevice">
          <GlowCard>
            <template #header>
              <div class="panel-header">
                <span>
                  {{ typeIcons[selectedDevice.type] || '📡' }}
                  {{ selectedDevice.name || selectedDevice.device_id }}
                </span>
                <div>
                  <el-tag :type="selectedDevice.online ? 'success' : 'info'" size="small" style="margin-right: 8px">
                    {{ selectedDevice.online ? '在线' : '离线' }}
                  </el-tag>
                  <el-button type="primary" size="small" @click="commandDialogVisible = true">下发指令</el-button>
                </div>
              </div>
            </template>

            <!-- 遥测数据 -->
            <div class="telemetry-section">
              <h4>实时遥测</h4>
              <div class="telemetry-grid">
                <div v-for="(val, key) in selectedDevice.properties" :key="key" class="telemetry-card">
                  <div class="telemetry-value">
                    <template v-if="key === 'relay_state'">
                      <span :style="{ color: val === 'on' ? '#10b981' : '#ef4444', fontSize: '24px' }">●</span>
                      {{ val === 'on' ? '已开启' : '已关闭' }}
                    </template>
                    <template v-else-if="key === 'temperature'">
                      <span :style="{ color: (val as number) > 30 ? '#ef4444' : (val as number) < 18 ? '#3b82f6' : '#10b981' }">{{ val }}</span>
                      <span class="unit">°C</span>
                    </template>
                    <template v-else>
                      {{ val }}<span v-if="key === 'humidity'" class="unit">%</span>
                      <span v-else-if="key === 'pressure'" class="unit">hPa</span>
                      <span v-else-if="key === 'power'" class="unit">W</span>
                      <span v-else-if="key === 'voltage'" class="unit">V</span>
                      <span v-else-if="key === 'current'" class="unit">A</span>
                    </template>
                  </div>
                  <div class="telemetry-label">{{ {
                    temperature: '温度', humidity: '湿度', pressure: '气压',
                    power: '功率', voltage: '电压', current: '电流',
                    relay_state: '继电器状态', unit: '单位'
                  }[key as string] || key }}</div>
                </div>
              </div>

              <el-empty v-if="!selectedDevice.properties || Object.keys(selectedDevice.properties).length === 0"
                description="等待设备上报数据..." :image-size="40" />
            </div>

            <!-- 设备信息 -->
            <div class="device-meta">
              <div class="meta-item">
                <span class="meta-label">设备ID</span>
                <code>{{ selectedDevice.device_id }}</code>
              </div>
              <div class="meta-item">
                <span class="meta-label">类型</span>
                <span>{{ typeLabels[selectedDevice.type] || selectedDevice.type }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">最后上报</span>
                <span>{{ selectedDevice.last_seen ? new Date(selectedDevice.last_seen).toLocaleTimeString() : '-' }}</span>
              </div>
            </div>
          </GlowCard>
        </template>

        <el-empty v-else description="选择左侧设备查看详情和遥测数据" :image-size="80" />
      </div>
    </div>

    <!-- 添加设备对话框 -->
    <el-dialog v-model="addDialogVisible" title="添加模拟设备" width="420px">
      <el-form label-width="80px">
        <el-form-item label="设备类型">
          <el-select v-model="newDevice.type" style="width: 100%">
            <el-option v-for="t in deviceTypes" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备名称">
          <el-input v-model="newDevice.name" placeholder="可选，留空自动生成" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="isLoading" @click="handleAddDevice">确定</el-button>
      </template>
    </el-dialog>

    <!-- 下发指令对话框 -->
    <el-dialog v-model="commandDialogVisible" title="下发给令" width="420px">
      <el-form label-width="80px">
        <el-form-item label="指令类型">
          <el-select v-model="commandForm.command" style="width: 100%" @change="onCommandChange">
            <el-option v-for="c in commandOptions" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="commandForm.command === 'set_interval'" label="间隔(ms)">
          <el-input-number v-model="commandForm.params.interval_ms" :min="1000" :max="60000" :step="1000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="commandDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSendCommand">下发</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.mqtt-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px 16px;
}

.page-header {
  margin-bottom: 24px;
}
.page-header h2 {
  margin: 0 0 4px 0;
  font-size: 22px;
  color: #e2e8f0;
}
.subtitle {
  color: #94a3b8;
  margin: 0;
  font-size: 14px;
}

.content-layout {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 20px;
  align-items: start;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.device-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.device-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
  border: 1px solid transparent;
}
.device-item:hover {
  background: rgba(255,255,255,0.04);
}
.device-item.active {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}
.device-item.offline {
  opacity: 0.5;
}
.device-icon {
  font-size: 22px;
  width: 32px;
  text-align: center;
}
.device-info {
  flex: 1;
  min-width: 0;
}
.device-name {
  font-size: 14px;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.device-type {
  font-size: 12px;
  color: #94a3b8;
}

.telemetry-section {
  margin-bottom: 20px;
}
.telemetry-section h4 {
  margin: 0 0 12px 0;
  color: #94a3b8;
  font-size: 13px;
  font-weight: 500;
}
.telemetry-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}
.telemetry-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  padding: 16px;
  text-align: center;
}
.telemetry-value {
  font-size: 28px;
  font-weight: 700;
  color: #e2e8f0;
  font-variant-numeric: tabular-nums;
}
.unit {
  font-size: 14px;
  font-weight: 400;
  color: #94a3b8;
  margin-left: 2px;
}
.telemetry-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.device-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid rgba(255,255,255,0.06);
}
.meta-item {
  display: flex;
  gap: 12px;
  font-size: 13px;
}
.meta-label {
  color: #64748b;
  min-width: 70px;
}
.meta-item code {
  background: rgba(255,255,255,0.05);
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #93c5fd;
}

@media (max-width: 900px) {
  .content-layout {
    grid-template-columns: 1fr;
  }
}
</style>
