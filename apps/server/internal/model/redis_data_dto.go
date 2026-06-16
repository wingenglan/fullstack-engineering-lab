package model

// ===== Redis 数据类型 DTO =====

// ----- Hash: 用户画像缓存 -----

// HashFieldUpdateRequest 单字段更新请求
type HashFieldUpdateRequest struct {
	Key   string `json:"key" binding:"required"`
	Field string `json:"field" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// HashMultiSetRequest 批量设置请求
type HashMultiSetRequest struct {
	Key    string            `json:"key" binding:"required"`
	Fields map[string]string `json:"fields" binding:"required,min=1"`
}

// HashProfileResponse 用户画像响应
type HashProfileResponse struct {
	Key      string            `json:"key"`
	Fields   map[string]string `json:"fields"`
	NumOfFld int               `json:"num_of_fld"`
}

// ----- Set: 标签/收藏夹管理 -----

// SetAddMembersRequest 添加成员请求
type SetAddMembersRequest struct {
	Key     string   `json:"key" binding:"required"`
	Members []string `json:"members" binding:"required,min=1"`
}

// SetRemoveMembersRequest 移除成员请求
type SetRemoveMembersRequest struct {
	Key     string   `json:"key" binding:"required"`
	Members []string `json:"members" binding:"required,min=1"`
}

// SetListRequest 查询请求
type SetListRequest struct {
	Key string `json:"key" binding:"required"`
}

// SetOperationRequest 集合运算请求
type SetOperationRequest struct {
	Keys []string `json:"keys" binding:"required,min=2,max=5"`
}

// SetListResponse 成员列表响应
type SetListResponse struct {
	Key     string   `json:"key"`
	Members []string `json:"members"`
	Count   int      `json:"count"`
}

// SetOperationResponse 集合运算响应
type SetOperationResponse struct {
	Keys    []string `json:"keys"`
	Op      string   `json:"op"`
	Members []string `json:"members"`
	Count   int      `json:"count"`
}

// ----- ZSet: 实时排行榜 -----

// ZSetAddScoreRequest 添加/增加分数请求
type ZSetAddScoreRequest struct {
	Key    string  `json:"key" binding:"required"`
	Member string  `json:"member" binding:"required"`
	Score  float64 `json:"score"`
}

// ZSetRankRequest 排名查询请求
type ZSetRankRequest struct {
	Key    string `json:"key" binding:"required"`
	Member string `json:"member" binding:"required"`
}

// ZSetTopNRequest Top N 请求
type ZSetTopNRequest struct {
	Key string `json:"key" binding:"required"`
	N   int    `json:"n" binding:"required,min=1,max=100"`
}

// ZSetMemberResponse 单个成员响应
type ZSetMemberResponse struct {
	Key    string  `json:"key"`
	Member string  `json:"member"`
	Score  float64 `json:"score"`
	Rank   int     `json:"rank"`
}

// ZSetTopNResponse Top N 响应
type ZSetTopNResponse struct {
	Key     string            `json:"key"`
	Total   int               `json:"total"`
	Members []ZSetRankEntry   `json:"members"`
}

// ZSetRankEntry 排行条目
type ZSetRankEntry struct {
	Rank   int     `json:"rank"`
	Member string  `json:"member"`
	Score  float64 `json:"score"`
}

// ----- List: 最新活动流 / 简易消息队列 -----

// ListPushRequest 推入消息请求
type ListPushRequest struct {
	Key    string `json:"key" binding:"required"`
	Value  string `json:"value" binding:"required"`
	Pos    string `json:"pos" binding:"required,oneof=left right"` // left=LPUSH, right=RPUSH
}

// ListPopRequest 弹出消息请求
type ListPopRequest struct {
	Key string `json:"key" binding:"required"`
	Pos string `json:"pos" binding:"required,oneof=left right"` // left=LPOP, right=RPOP
}

// ListRangeRequest 查询范围请求
type ListRangeRequest struct {
	Key   string `json:"key" binding:"required"`
	Start int64  `json:"start"`
	Stop  int64  `json:"stop"`
}

// ListPopResponse 弹出消息响应
type ListPopResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Pos   string `json:"pos"`
}

// ListRangeResponse 范围查询响应
type ListRangeResponse struct {
	Key     string   `json:"key"`
	Values  []string `json:"values"`
	Length  int64    `json:"length"`
}

// ----- String: 验证码 / 计数器 -----

// StringSetRequest 设置键值请求（支持 TTL）
type StringSetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl"` // 过期秒数，0 表示不过期
}

// StringGetRequest 获取键值请求
type StringGetRequest struct {
	Key string `json:"key" binding:"required"`
}

// StringIncrRequest 自增计数请求
type StringIncrRequest struct {
	Key   string `json:"key" binding:"required"`
	Delta int64  `json:"delta"` // 自增步长，默认 1
}

// StringSetResponse 设置响应
type StringSetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

// StringGetResponse 获取响应
type StringGetResponse struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	TTL     int    `json:"ttl"`  // -1=无过期，-2=不存在
	Exists  bool   `json:"exists"`
}

// StringIncrResponse 自增响应
type StringIncrResponse struct {
	Key    string `json:"key"`
	Before int64  `json:"before"`
	After  int64  `json:"after"`
	Delta  int64  `json:"delta"`
}
