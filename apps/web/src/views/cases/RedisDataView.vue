<template>
  <div class="redis-data-view py-8 px-6">
    <div class="max-w-7xl mx-auto">
      <!-- 面包屑 -->
      <div class="flex items-center gap-2 text-sm text-text-muted mb-6">
        <router-link to="/cases" class="hover:text-text-primary transition-colors">案例列表</router-link>
        <i class="i-lucide-chevron-right w-4 h-4" />
        <span class="text-text-primary">Redis 数据类型实战</span>
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <!-- 左栏 -->
        <div class="lg:col-span-1 space-y-6">
          <!-- 案例信息 -->
          <GlowCard>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-lg bg-brand/10 flex items-center justify-center">
                <i class="i-lucide-database w-5 h-5 text-brand" />
              </div>
              <div>
                <h2 class="text-text-primary font-bold text-lg">Redis 数据类型实战</h2>
                <span class="text-xs text-success">已上线</span>
              </div>
            </div>
            <p class="text-text-secondary text-sm mb-4">
              覆盖 Redis 五大核心数据结构的实战场景：String 验证码/计数器、Hash 用户画像、List 活动流、Set 标签管理、ZSet 排行榜。
            </p>
            <div class="flex flex-wrap gap-2 mb-4">
              <span v-for="tag in ['Redis', 'Go', 'String', 'Hash', 'List', 'Set', 'ZSet']" :key="tag" class="px-2.5 py-0.5 text-xs rounded-full bg-white/5 text-text-secondary border border-border">
                {{ tag }}
              </span>
            </div>
            <div class="text-xs text-text-muted">
              难度：<span class="text-warning font-medium">进阶</span>
            </div>
          </GlowCard>

          <!-- 操作步骤 -->
          <GlowCard title="操作步骤">
            <div class="space-y-2">
              <div
                v-for="(step, i) in steps"
                :key="i"
                class="flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer transition-all"
                :class="activeStep === i ? 'bg-brand/10 text-brand' : 'text-text-secondary hover:bg-white/5'"
                @click="activeStep = i"
              >
                <div
                  class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                  :class="step.completed ? 'bg-success text-white' : activeStep === i ? 'bg-brand text-white' : 'bg-white/10 text-text-muted'"
                >
                  <i v-if="step.completed" class="i-lucide-check w-3 h-3" />
                  <span v-else>{{ i + 1 }}</span>
                </div>
                <span class="text-sm font-medium">{{ step.label }}</span>
              </div>
            </div>
          </GlowCard>

          <!-- 核心原理 -->
          <GlowCard title="核心原理">
            <div class="space-y-3 text-sm">
              <div v-for="(item, i) in concepts" :key="i" class="flex items-start gap-3">
                <div class="w-7 h-7 rounded-full bg-brand/10 flex items-center justify-center shrink-0 mt-0.5">
                  <i :class="item.icon" class="w-3.5 h-3.5 text-brand" />
                </div>
                <div>
                  <span class="text-text-primary font-medium block">{{ item.title }}</span>
                  <span class="text-text-muted text-xs">{{ item.desc }}</span>
                </div>
              </div>
            </div>
          </GlowCard>
        </div>

        <!-- 右栏 -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 数据类型切换标签 -->
          <div class="flex gap-1 p-1 bg-bg-elevated rounded-xl border border-border overflow-x-auto">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              class="flex-1 px-3 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap"
              :class="activeTab === tab.key ? 'bg-brand text-white shadow-lg' : 'text-text-secondary hover:text-text-primary'"
              @click="activeTab = tab.key"
            >
              <i :class="tab.icon" class="w-4 h-4 mr-1.5 inline-block" />
              {{ tab.label }}
            </button>
          </div>

          <!-- ===== String: 验证码/计数器 ===== -->
          <template v-if="activeTab === 'string'">
            <GlowCard>
              <h3 class="text-text-primary font-semibold text-lg mb-2">String — 验证码存储 / 计数器</h3>
              <p class="text-text-secondary text-sm mb-6">
                String 是最基础的类型。本场景演示：① 存储验证码并设置 TTL 自动过期；② 使用 INCR 实现原子自增计数器（如文章阅读量）。
              </p>

              <div class="space-y-6">
                <!-- 验证码存储 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-key-round w-4 h-4 text-brand" />
                    验证码存储（SET with TTL）
                  </h4>
                  <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">Key</label>
                      <el-input v-model="stringForm.key" placeholder="sms:code:13800001111" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">验证码</label>
                      <el-input v-model="stringForm.value" placeholder="583214" />
                    </div>
                  </div>
                  <div class="mb-4">
                    <div class="flex items-center justify-between mb-1">
                      <label class="text-xs text-text-muted">过期时间（秒）</label>
                      <span class="text-xs font-mono text-brand">{{ stringForm.ttl }}s</span>
                    </div>
                    <el-slider v-model="stringForm.ttl" :min="10" :max="600" :step="10" />
                  </div>
                  <div class="flex gap-3">
                    <el-button type="primary" :loading="stringLoading" @click="handleStringSet">
                      <i class="i-lucide-save w-4 h-4 mr-1" />
                      存储验证码
                    </el-button>
                    <el-button :loading="stringLoading" @click="handleStringGet">
                      <i class="i-lucide-search w-4 h-4 mr-1" />
                      查询验证码
                    </el-button>
                  </div>
                  <div v-if="stringGetResult" class="mt-4 p-3 rounded-lg border" :class="stringGetResult.exists ? 'bg-success/5 border-success/20' : 'bg-danger/5 border-danger/20'">
                    <div class="flex items-center gap-2 mb-1">
                      <i :class="stringGetResult.exists ? 'i-lucide-check-circle text-success' : 'i-lucide-x-circle text-danger'" class="w-4 h-4" />
                      <span class="text-sm font-medium" :class="stringGetResult.exists ? 'text-success' : 'text-danger'">
                        {{ stringGetResult.exists ? '查询成功' : 'Key 不存在或已过期' }}
                      </span>
                    </div>
                    <pre v-if="stringGetResult.exists" class="text-xs font-mono text-text-secondary bg-bg-card rounded p-2 mt-2">{{ JSON.stringify(stringGetResult, null, 2) }}</pre>
                  </div>
                </div>

                <!-- 计数器 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-trending-up w-4 h-4 text-brand" />
                    原子计数器（INCR）
                  </h4>
                  <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">计数器 Key</label>
                      <el-input v-model="incrForm.key" placeholder="article:10086:views" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">步长（负数为递减）</label>
                      <el-input-number v-model="incrForm.delta" :min="-100" :max="100" class="w-full" />
                    </div>
                  </div>
                  <el-button type="primary" :loading="incrLoading" @click="handleStringIncr">
                    <i class="i-lucide-plus w-4 h-4 mr-1" />
                    执行 INCR
                  </el-button>
                  <div v-if="incrResult" class="mt-4 grid grid-cols-3 gap-3">
                    <div class="p-3 rounded-lg bg-bg-card border border-border text-center">
                      <span class="text-lg font-bold text-text-muted block">{{ incrResult.before }}</span>
                      <span class="text-xs text-text-muted">操作前</span>
                    </div>
                    <div class="p-3 rounded-lg bg-brand/10 border border-brand/20 text-center">
                      <span class="text-lg font-bold text-brand block">{{ incrResult.delta > 0 ? '+' : '' }}{{ incrResult.delta }}</span>
                      <span class="text-xs text-brand/70">变化量</span>
                    </div>
                    <div class="p-3 rounded-lg bg-success/5 border border-success/20 text-center">
                      <span class="text-lg font-bold text-success block">{{ incrResult.after }}</span>
                      <span class="text-xs text-success/70">操作后</span>
                    </div>
                  </div>
                </div>
              </div>
            </GlowCard>
          </template>

          <!-- ===== Hash: 用户画像缓存 ===== -->
          <template v-if="activeTab === 'hash'">
            <GlowCard>
              <h3 class="text-text-primary font-semibold text-lg mb-2">Hash — 用户画像缓存</h3>
              <p class="text-text-secondary text-sm mb-6">
                Hash 适合存储对象属性。本场景演示：用 Hash 缓存用户画像（昵称、邮箱、等级等），支持单字段读写、批量设置和字段删除。
              </p>

              <div class="space-y-6">
                <!-- 批量设置 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-user w-4 h-4 text-brand" />
                    设置用户画像
                  </h4>
                  <div class="mb-4">
                    <label class="text-xs text-text-muted mb-1 block">Key</label>
                    <el-input v-model="hashForm.key" placeholder="user:10086" />
                  </div>
                  <div class="mb-4">
                    <label class="text-xs text-text-muted mb-2 block">字段列表</label>
                    <div class="space-y-2">
                      <div v-for="(val, field, i) in hashForm.fields" :key="i" class="flex items-center gap-2">
                        <el-input :model-value="field as string" disabled class="w-32" />
                        <el-input v-model="hashForm.fields[field as string]" class="flex-1" />
                        <el-button size="small" type="danger" plain @click="delete hashForm.fields[field as string]">
                          <i class="i-lucide-trash-2 w-3 h-3" />
                        </el-button>
                      </div>
                    </div>
                    <div class="flex gap-2 mt-3">
                      <el-input v-model="newHashField.key" placeholder="字段名" class="w-32" size="small" />
                      <el-input v-model="newHashField.value" placeholder="字段值" class="flex-1" size="small" />
                      <el-button size="small" @click="addHashField">
                        <i class="i-lucide-plus w-3 h-3" />
                      </el-button>
                    </div>
                  </div>
                  <div class="flex gap-3">
                    <el-button type="primary" :loading="hashLoading" @click="handleHashMultiSet">
                      <i class="i-lucide-save w-4 h-4 mr-1" />
                      批量设置
                    </el-button>
                    <el-button :loading="hashLoading" @click="handleHashGetProfile">
                      <i class="i-lucide-search w-4 h-4 mr-1" />
                      查询画像
                    </el-button>
                  </div>
                </div>

                <!-- 画像展示 -->
                <div v-if="hashResult" class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-3">
                    画像数据 <span class="text-text-muted text-xs">（{{ hashResult.num_of_fld }} 个字段）</span>
                  </h4>
                  <div class="grid grid-cols-2 gap-3">
                    <div v-for="(val, field) in hashResult.fields" :key="field" class="p-3 rounded-lg bg-bg-card border border-border">
                      <span class="text-xs text-text-muted block mb-1">{{ field }}</span>
                      <span class="text-sm text-text-primary font-mono">{{ val }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </GlowCard>
          </template>

          <!-- ===== List: 最新活动流 ===== -->
          <template v-if="activeTab === 'list'">
            <GlowCard>
              <h3 class="text-text-primary font-semibold text-lg mb-2">List — 最新活动流 / 简易消息队列</h3>
              <p class="text-text-secondary text-sm mb-6">
                List 是有序的字符串列表。本场景演示：① 用 LPUSH/LRANGE 构建最新活动流；② 用 LPOP/RPOP 实现简易消息队列的消费语义。
              </p>

              <div class="space-y-6">
                <!-- 推入消息 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-arrow-down-to-line w-4 h-4 text-brand" />
                    推入消息
                  </h4>
                  <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">Key</label>
                      <el-input v-model="listForm.key" placeholder="activity:feed" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">消息内容</label>
                      <el-input v-model="listForm.value" placeholder="用户张三购买了商品" />
                    </div>
                  </div>
                  <div class="mb-4">
                    <label class="text-xs text-text-muted mb-1 block">推入方向</label>
                    <el-radio-group v-model="listForm.pos">
                      <el-radio value="left">LPUSH（头部插入，最新优先）</el-radio>
                      <el-radio value="right">RPUSH（尾部插入，先进先出）</el-radio>
                    </el-radio-group>
                  </div>
                  <el-button type="primary" :loading="listLoading" @click="handleListPush">
                    <i class="i-lucide-plus w-4 h-4 mr-1" />
                    推入
                  </el-button>
                </div>

                <!-- 消费消息 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-arrow-up-from-line w-4 h-4 text-brand" />
                    弹出消息（消费）
                  </h4>
                  <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">Key</label>
                      <el-input v-model="listPopForm.key" placeholder="activity:feed" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">弹出方向</label>
                      <el-radio-group v-model="listPopForm.pos">
                        <el-radio value="left">LPOP（头部弹出）</el-radio>
                        <el-radio value="right">RPOP（尾部弹出）</el-radio>
                      </el-radio-group>
                    </div>
                  </div>
                  <el-button type="danger" plain :loading="listLoading" @click="handleListPop">
                    <i class="i-lucide-minus w-4 h-4 mr-1" />
                    弹出
                  </el-button>
                  <div v-if="listPopResult" class="mt-3 p-3 rounded-lg bg-brand/5 border border-brand/20">
                    <span class="text-xs text-text-muted">弹出内容：</span>
                    <span class="text-sm font-mono text-brand">{{ listPopResult.value }}</span>
                  </div>
                </div>

                <!-- 列表展示 -->
                <div v-if="listResult" class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <div class="flex items-center justify-between mb-3">
                    <h4 class="text-sm font-medium text-text-primary">
                      列表内容
                    </h4>
                    <span class="text-xs text-text-muted">{{ listResult.length }} 条记录</span>
                  </div>
                  <div class="space-y-1">
                    <div
                      v-for="(val, i) in listResult.values"
                      :key="i"
                      class="flex items-center gap-3 px-3 py-2 rounded-lg bg-bg-card border border-border text-sm"
                    >
                      <span class="text-xs text-text-muted font-mono w-6 shrink-0">#{{ i }}</span>
                      <span class="text-text-secondary font-mono text-xs flex-1 truncate">{{ val }}</span>
                    </div>
                    <div v-if="listResult.values.length === 0" class="text-center py-4 text-text-muted text-xs">
                      列表为空
                    </div>
                  </div>
                </div>
              </div>
            </GlowCard>
          </template>

          <!-- ===== Set: 标签管理 ===== -->
          <template v-if="activeTab === 'set'">
            <GlowCard>
              <h3 class="text-text-primary font-semibold text-lg mb-2">Set — 标签/收藏夹管理</h3>
              <p class="text-text-secondary text-sm mb-6">
                Set 是无序且去重的集合。本场景演示：① 管理文章标签（添加/移除/查询）；② 集合运算——求两个标签集合的交集、并集、差集。
              </p>

              <div class="space-y-6">
                <!-- 标签管理 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-tag w-4 h-4 text-brand" />
                    标签管理
                  </h4>
                  <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">Key</label>
                      <el-input v-model="setForm.key" placeholder="article:100:tags" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">标签（逗号分隔）</label>
                      <el-input v-model="setForm.membersInput" placeholder="Go,Redis,微服务" />
                    </div>
                  </div>
                  <div class="flex gap-3">
                    <el-button type="primary" :loading="setLoading" @click="handleSetAdd">
                      <i class="i-lucide-plus w-4 h-4 mr-1" />
                      添加
                    </el-button>
                    <el-button type="danger" plain :loading="setLoading" @click="handleSetRemove">
                      <i class="i-lucide-minus w-4 h-4 mr-1" />
                      移除
                    </el-button>
                    <el-button :loading="setLoading" @click="handleSetMembers">
                      <i class="i-lucide-search w-4 h-4 mr-1" />
                      查询
                    </el-button>
                  </div>
                  <div v-if="setResult" class="mt-4">
                    <div class="flex flex-wrap gap-2">
                      <span
                        v-for="m in setResult.members"
                        :key="m"
                        class="px-2.5 py-1 text-xs rounded-full bg-brand/10 text-brand border border-brand/20"
                      >
                        {{ m }}
                      </span>
                    </div>
                    <span class="text-xs text-text-muted mt-2 block">共 {{ setResult.count }} 个标签</span>
                  </div>
                </div>

                <!-- 集合运算 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-git-merge w-4 h-4 text-brand" />
                    集合运算
                  </h4>
                  <div class="mb-4">
                    <label class="text-xs text-text-muted mb-1 block">Keys（逗号分隔，至少 2 个）</label>
                    <el-input v-model="setOpForm.keysInput" placeholder="article:100:tags,article:200:tags" />
                  </div>
                  <div class="flex gap-3">
                    <el-button :loading="setLoading" @click="handleSetOp('intersect')">交集</el-button>
                    <el-button :loading="setLoading" @click="handleSetOp('union')">并集</el-button>
                    <el-button :loading="setLoading" @click="handleSetOp('diff')">差集</el-button>
                  </div>
                  <div v-if="setOpResult" class="mt-4">
                    <span class="text-xs text-text-muted block mb-2">{{ setOpResult.op }} 结果（{{ setOpResult.count }} 个）：</span>
                    <div class="flex flex-wrap gap-2">
                      <span
                        v-for="m in setOpResult.members"
                        :key="m"
                        class="px-2.5 py-1 text-xs rounded-full bg-success/10 text-success border border-success/20"
                      >
                        {{ m }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </GlowCard>
          </template>

          <!-- ===== ZSet: 实时排行榜 ===== -->
          <template v-if="activeTab === 'zset'">
            <GlowCard>
              <h3 class="text-text-primary font-semibold text-lg mb-2">ZSet — 实时排行榜</h3>
              <p class="text-text-secondary text-sm mb-6">
                Sorted Set 按分数自动排序，天然适合排行榜场景。本场景演示：① ZINCRBY 增加分数；② ZREVRANGE 获取 Top N；③ ZREVRANK 查询排名。
              </p>

              <div class="space-y-6">
                <!-- 增加分数 -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-trophy w-4 h-4 text-brand" />
                    增加分数
                  </h4>
                  <div class="grid grid-cols-3 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">排行榜 Key</label>
                      <el-input v-model="zsetForm.key" placeholder="rank:game:daily" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">玩家名</label>
                      <el-input v-model="zsetForm.member" placeholder="player_001" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">分数增量</label>
                      <el-input-number v-model="zsetForm.score" :min="1" :max="1000" class="w-full" />
                    </div>
                  </div>
                  <div class="flex gap-3">
                    <el-button type="primary" :loading="zsetLoading" @click="handleZsetAddScore">
                      <i class="i-lucide-plus w-4 h-4 mr-1" />
                      增加分数
                    </el-button>
                    <el-button :loading="zsetLoading" @click="handleZsetGetRank">
                      <i class="i-lucide-search w-4 h-4 mr-1" />
                      查询排名
                    </el-button>
                  </div>
                  <div v-if="zsetRankResult" class="mt-4 p-3 rounded-lg bg-brand/5 border border-brand/20">
                    <div class="flex items-center gap-4">
                      <span class="text-xs text-text-muted">玩家：</span>
                      <span class="text-sm font-mono text-text-primary">{{ zsetRankResult.member }}</span>
                      <span class="text-xs text-text-muted">分数：</span>
                      <span class="text-sm font-mono text-brand">{{ zsetRankResult.score }}</span>
                      <span class="text-xs text-text-muted">排名：</span>
                      <span class="text-lg font-bold text-brand">#{{ zsetRankResult.rank }}</span>
                    </div>
                  </div>
                </div>

                <!-- Top N -->
                <div class="p-4 rounded-lg bg-bg-elevated border border-border">
                  <h4 class="text-sm font-medium text-text-primary mb-4 flex items-center gap-2">
                    <i class="i-lucide-bar-chart-3 w-4 h-4 text-brand" />
                    排行榜 Top N
                  </h4>
                  <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">排行榜 Key</label>
                      <el-input v-model="zsetTopForm.key" placeholder="rank:game:daily" />
                    </div>
                    <div>
                      <label class="text-xs text-text-muted mb-1 block">Top N</label>
                      <el-input-number v-model="zsetTopForm.n" :min="1" :max="50" class="w-full" />
                    </div>
                  </div>
                  <el-button type="primary" :loading="zsetLoading" @click="handleZsetTopN">
                    <i class="i-lucide-list-ordered w-4 h-4 mr-1" />
                    获取排行榜
                  </el-button>
                  <div v-if="zsetTopResult" class="mt-4 space-y-2">
                    <div class="flex items-center justify-between mb-2">
                      <span class="text-xs text-text-muted">共 {{ zsetTopResult.total }} 人</span>
                    </div>
                    <div
                      v-for="entry in zsetTopResult.members"
                      :key="entry.member"
                      class="flex items-center gap-3 px-3 py-2.5 rounded-lg"
                      :class="entry.rank <= 3 ? 'bg-brand/5 border border-brand/20' : 'bg-bg-card border border-border'"
                    >
                      <span
                        class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                        :class="entry.rank === 1 ? 'bg-yellow-500 text-black' : entry.rank === 2 ? 'bg-gray-400 text-black' : entry.rank === 3 ? 'bg-amber-700 text-white' : 'bg-white/10 text-text-muted'"
                      >
                        {{ entry.rank }}
                      </span>
                      <span class="text-sm text-text-primary font-mono flex-1">{{ entry.member }}</span>
                      <span class="text-sm font-mono text-brand">{{ entry.score }}</span>
                    </div>
                    <div v-if="zsetTopResult.members.length === 0" class="text-center py-4 text-text-muted text-xs">
                      排行榜为空
                    </div>
                  </div>
                </div>
              </div>
            </GlowCard>
          </template>

          <!-- 请求日志 -->
          <div class="rounded-lg border border-border overflow-hidden">
            <div class="flex items-center justify-between px-4 py-2.5 bg-bg-elevated border-b border-border">
              <div class="flex items-center gap-2">
                <i class="i-lucide-terminal w-4 h-4 text-brand" />
                <span class="text-sm font-medium text-text-primary">请求日志</span>
              </div>
              <span class="text-xs text-text-muted">{{ requestLogs.length }} 条请求</span>
            </div>
            <RequestLog :logs="requestLogs" />
          </div>

          <!-- 接口列表 -->
          <GlowCard title="接口列表">
            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-border">
                    <th class="text-left py-3 text-text-muted font-medium">方法</th>
                    <th class="text-left py-3 text-text-muted font-medium">接口</th>
                    <th class="text-left py-3 text-text-muted font-medium">说明</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="api in apiList" :key="api.endpoint" class="border-b border-border/50">
                    <td class="py-2.5">
                      <span class="px-2 py-0.5 rounded text-xs font-bold" :class="api.method === 'GET' ? 'bg-success/20 text-success' : 'bg-brand/20 text-brand'">
                        {{ api.method }}
                      </span>
                    </td>
                    <td class="py-2.5 font-mono text-text-secondary text-xs">{{ api.endpoint }}</td>
                    <td class="py-2.5 text-text-secondary">{{ api.desc }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </GlowCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import * as redisApi from '@/api/redisData'
import { requestLogs } from '@/api/request'
import GlowCard from '@/components/GlowCard.vue'
import RequestLog from '@/components/RequestLog.vue'
import type {
  HashProfileResponse,
  SetListResponse,
  SetOperationResponse,
  ListRangeResponse,
  ListPopResponse,
  StringGetResponse,
  StringIncrResponse,
  ZSetMemberResponse,
  ZSetTopNResponse,
} from '@/types'

const activeStep = ref(0)
const activeTab = ref('string')

const tabs = [
  { key: 'string', label: 'String', icon: 'i-lucide-key-round' },
  { key: 'hash', label: 'Hash', icon: 'i-lucide-user' },
  { key: 'list', label: 'List', icon: 'i-lucide-list' },
  { key: 'set', label: 'Set', icon: 'i-lucide-tag' },
  { key: 'zset', label: 'ZSet', icon: 'i-lucide-trophy' },
]

const steps = computed(() => [
  { label: 'String 验证码/计数器', completed: !!stringGetResult.value || !!incrResult.value },
  { label: 'Hash 用户画像', completed: !!hashResult.value },
  { label: 'List 活动流', completed: !!listResult.value },
  { label: 'Set 标签管理', completed: !!setResult.value },
  { label: 'ZSet 排行榜', completed: !!zsetTopResult.value },
])

const concepts = [
  { title: 'String + TTL', desc: '验证码存储：SET key code EX 300，5 分钟自动过期', icon: 'i-lucide-key-round' },
  { title: 'INCR 原子自增', desc: '计数器：多个客户端并发 INCR 不会丢失计数', icon: 'i-lucide-trending-up' },
  { title: 'Hash 字段操作', desc: '对象缓存：HSET/HGET/HDEL 操作单个字段，HMSET 批量设置', icon: 'i-lucide-user' },
  { title: 'List 有序列表', desc: 'LPUSH + LRANGE 实现最新列表，LPOP/RPOP 实现消息队列', icon: 'i-lucide-list' },
  { title: 'Set 去重 + 运算', desc: 'SADD 自动去重，SINTER/SUNION/SDIFF 集合运算', icon: 'i-lucide-tag' },
  { title: 'ZSet 排序', desc: 'ZINCRBY 增分，ZREVRANGE Top N，ZREVRANK 排名查询', icon: 'i-lucide-trophy' },
]

// ===== String 状态 =====
const stringForm = reactive({ key: 'sms:code:13800001111', value: '583214', ttl: 300 })
const stringLoading = ref(false)
const stringGetResult = ref<StringGetResponse | null>(null)
const incrForm = reactive({ key: 'article:10086:views', delta: 1 })
const incrLoading = ref(false)
const incrResult = ref<StringIncrResponse | null>(null)

// ===== Hash 状态 =====
const hashForm = reactive({
  key: 'user:10086',
  fields: {} as Record<string, string>,
})
const newHashField = reactive({ key: '', value: '' })
const hashLoading = ref(false)
const hashResult = ref<HashProfileResponse | null>(null)

// ===== List 状态 =====
const listForm = reactive({ key: 'activity:feed', value: '', pos: 'left' as 'left' | 'right' })
const listPopForm = reactive({ key: 'activity:feed', pos: 'left' as 'left' | 'right' })
const listLoading = ref(false)
const listResult = ref<ListRangeResponse | null>(null)
const listPopResult = ref<ListPopResponse | null>(null)

// ===== Set 状态 =====
const setForm = reactive({ key: 'article:100:tags', membersInput: 'Go,Redis,微服务' })
const setOpForm = reactive({ keysInput: 'article:100:tags,article:200:tags' })
const setLoading = ref(false)
const setResult = ref<SetListResponse | null>(null)
const setOpResult = ref<SetOperationResponse | null>(null)

// ===== ZSet 状态 =====
const zsetForm = reactive({ key: 'rank:game:daily', member: 'player_001', score: 10 })
const zsetTopForm = reactive({ key: 'rank:game:daily', n: 10 })
const zsetLoading = ref(false)
const zsetRankResult = ref<ZSetMemberResponse | null>(null)
const zsetTopResult = ref<ZSetTopNResponse | null>(null)

// ===== String 操作 =====
async function handleStringSet() {
  if (!stringForm.key || !stringForm.value) {
    ElMessage.warning('请填写 Key 和验证码')
    return
  }
  stringLoading.value = true
  try {
    await redisApi.stringSet({ key: stringForm.key, value: stringForm.value, ttl: stringForm.ttl })
    ElMessage.success(`验证码已存储，${stringForm.ttl}s 后过期`)
    activeStep.value = 0
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    stringLoading.value = false
  }
}

async function handleStringGet() {
  if (!stringForm.key) {
    ElMessage.warning('请填写 Key')
    return
  }
  stringLoading.value = true
  try {
    const res = await redisApi.stringGet(stringForm.key)
    stringGetResult.value = res.data
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '查询失败')
  } finally {
    stringLoading.value = false
  }
}

async function handleStringIncr() {
  if (!incrForm.key) {
    ElMessage.warning('请填写计数器 Key')
    return
  }
  incrLoading.value = true
  try {
    const res = await redisApi.stringIncr({ key: incrForm.key, delta: incrForm.delta })
    incrResult.value = res.data
    activeStep.value = 0
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    incrLoading.value = false
  }
}

// ===== Hash 操作 =====
function addHashField() {
  if (newHashField.key && newHashField.value) {
    hashForm.fields[newHashField.key] = newHashField.value
    newHashField.key = ''
    newHashField.value = ''
  }
}

async function handleHashMultiSet() {
  if (!hashForm.key || Object.keys(hashForm.fields).length === 0) {
    ElMessage.warning('请填写 Key 和至少一个字段')
    return
  }
  hashLoading.value = true
  try {
    const res = await redisApi.multiSetHash({ key: hashForm.key, fields: { ...hashForm.fields } })
    hashResult.value = res.data
    ElMessage.success('画像设置成功')
    activeStep.value = 1
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    hashLoading.value = false
  }
}

async function handleHashGetProfile() {
  if (!hashForm.key) {
    ElMessage.warning('请填写 Key')
    return
  }
  hashLoading.value = true
  try {
    const res = await redisApi.getHashProfile(hashForm.key)
    hashResult.value = res.data
    activeStep.value = 1
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '查询失败')
  } finally {
    hashLoading.value = false
  }
}

// ===== List 操作 =====
async function handleListPush() {
  if (!listForm.key || !listForm.value) {
    ElMessage.warning('请填写 Key 和消息内容')
    return
  }
  listLoading.value = true
  try {
    const res = await redisApi.listPush({ key: listForm.key, value: listForm.value, pos: listForm.pos })
    listResult.value = res.data
    listForm.value = ''
    ElMessage.success('消息已推入')
    activeStep.value = 2
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    listLoading.value = false
  }
}

async function handleListPop() {
  if (!listPopForm.key) {
    ElMessage.warning('请填写 Key')
    return
  }
  listLoading.value = true
  try {
    const res = await redisApi.listPop({ key: listPopForm.key, pos: listPopForm.pos })
    listPopResult.value = res.data
    // 刷新列表
    const listRes = await redisApi.listRange(listPopForm.key)
    listResult.value = listRes.data
    ElMessage.success('消息已弹出')
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    listLoading.value = false
  }
}

// ===== Set 操作 =====
function parseMembers(input: string): string[] {
  return input.split(',').map(s => s.trim()).filter(Boolean)
}

async function handleSetAdd() {
  const members = parseMembers(setForm.membersInput)
  if (!setForm.key || members.length === 0) {
    ElMessage.warning('请填写 Key 和至少一个标签')
    return
  }
  setLoading.value = true
  try {
    const res = await redisApi.setAddMembers({ key: setForm.key, members })
    setResult.value = res.data
    ElMessage.success('标签已添加')
    activeStep.value = 3
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    setLoading.value = false
  }
}

async function handleSetRemove() {
  const members = parseMembers(setForm.membersInput)
  if (!setForm.key || members.length === 0) {
    ElMessage.warning('请填写 Key 和至少一个标签')
    return
  }
  setLoading.value = true
  try {
    const res = await redisApi.setRemoveMembers({ key: setForm.key, members })
    setResult.value = res.data
    ElMessage.success('标签已移除')
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    setLoading.value = false
  }
}

async function handleSetMembers() {
  if (!setForm.key) {
    ElMessage.warning('请填写 Key')
    return
  }
  setLoading.value = true
  try {
    const res = await redisApi.getSetMembers(setForm.key)
    setResult.value = res.data
    activeStep.value = 3
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '查询失败')
  } finally {
    setLoading.value = false
  }
}

async function handleSetOp(op: 'intersect' | 'union' | 'diff') {
  const keys = parseMembers(setOpForm.keysInput)
  if (keys.length < 2) {
    ElMessage.warning('请至少输入 2 个 Key')
    return
  }
  setLoading.value = true
  try {
    const apiFn = op === 'intersect' ? redisApi.setIntersect : op === 'union' ? redisApi.setUnion : redisApi.setDiff
    const res = await apiFn({ keys })
    setOpResult.value = res.data
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    setLoading.value = false
  }
}

// ===== ZSet 操作 =====
async function handleZsetAddScore() {
  if (!zsetForm.key || !zsetForm.member) {
    ElMessage.warning('请填写排行榜 Key 和玩家名')
    return
  }
  zsetLoading.value = true
  try {
    const res = await redisApi.zsetAddScore({ key: zsetForm.key, member: zsetForm.member, score: zsetForm.score })
    zsetRankResult.value = res.data
    ElMessage.success('分数已增加')
    activeStep.value = 4
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败')
  } finally {
    zsetLoading.value = false
  }
}

async function handleZsetGetRank() {
  if (!zsetForm.key || !zsetForm.member) {
    ElMessage.warning('请填写排行榜 Key 和玩家名')
    return
  }
  zsetLoading.value = true
  try {
    const res = await redisApi.zsetGetRank({ key: zsetForm.key, member: zsetForm.member })
    zsetRankResult.value = res.data
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '查询失败')
  } finally {
    zsetLoading.value = false
  }
}

async function handleZsetTopN() {
  if (!zsetTopForm.key) {
    ElMessage.warning('请填写排行榜 Key')
    return
  }
  zsetLoading.value = true
  try {
    const res = await redisApi.zsetTopN({ key: zsetTopForm.key, n: zsetTopForm.n })
    zsetTopResult.value = res.data
    activeStep.value = 4
  } catch (err: unknown) {
    ElMessage.error(err instanceof Error ? err.message : '查询失败')
  } finally {
    zsetLoading.value = false
  }
}

// ===== API 列表 =====
const apiList = [
  { method: 'POST', endpoint: '/api/v1/redis-data/string/set', desc: '存储键值（支持 TTL）' },
  { method: 'GET', endpoint: '/api/v1/redis-data/string/get', desc: '查询键值' },
  { method: 'POST', endpoint: '/api/v1/redis-data/string/incr', desc: '原子自增计数' },
  { method: 'POST', endpoint: '/api/v1/redis-data/hash/field', desc: '设置 Hash 单字段' },
  { method: 'GET', endpoint: '/api/v1/redis-data/hash/profile', desc: '查询 Hash 完整画像' },
  { method: 'POST', endpoint: '/api/v1/redis-data/hash/multi-set', desc: '批量设置 Hash 字段' },
  { method: 'POST', endpoint: '/api/v1/redis-data/hash/delete-field', desc: '删除 Hash 字段' },
  { method: 'POST', endpoint: '/api/v1/redis-data/list/push', desc: '推入消息（LPUSH/RPUSH）' },
  { method: 'POST', endpoint: '/api/v1/redis-data/list/pop', desc: '弹出消息（LPOP/RPOP）' },
  { method: 'GET', endpoint: '/api/v1/redis-data/list/range', desc: '查询列表范围' },
  { method: 'POST', endpoint: '/api/v1/redis-data/set/add', desc: '添加集合成员' },
  { method: 'POST', endpoint: '/api/v1/redis-data/set/remove', desc: '移除集合成员' },
  { method: 'GET', endpoint: '/api/v1/redis-data/set/members', desc: '查询集合成员' },
  { method: 'POST', endpoint: '/api/v1/redis-data/set/intersect', desc: '集合交集' },
  { method: 'POST', endpoint: '/api/v1/redis-data/set/union', desc: '集合并集' },
  { method: 'POST', endpoint: '/api/v1/redis-data/set/diff', desc: '集合差集' },
  { method: 'POST', endpoint: '/api/v1/redis-data/zset/add-score', desc: '增加成员分数' },
  { method: 'POST', endpoint: '/api/v1/redis-data/zset/top-n', desc: '获取 Top N' },
  { method: 'POST', endpoint: '/api/v1/redis-data/zset/rank', desc: '查询成员排名' },
]
</script>
