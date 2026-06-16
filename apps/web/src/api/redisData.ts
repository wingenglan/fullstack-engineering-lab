import request from './request'
import type {
  ApiResponse,
  HashFieldUpdateRequest,
  HashMultiSetRequest,
  HashProfileResponse,
  SetAddMembersRequest,
  SetRemoveMembersRequest,
  SetListResponse,
  SetOperationRequest,
  SetOperationResponse,
  ZSetAddScoreRequest,
  ZSetRankRequest,
  ZSetTopNRequest,
  ZSetMemberResponse,
  ZSetTopNResponse,
  ListPushRequest,
  ListPopRequest,
  ListRangeResponse,
  ListPopResponse,
  StringSetRequest,
  StringSetResponse,
  StringGetResponse,
  StringIncrRequest,
  StringIncrResponse,
} from '@/types'

// ===== Hash: 用户画像缓存 =====

/** 设置单个 Hash 字段 */
export function setHashField(data: HashFieldUpdateRequest): Promise<ApiResponse<HashProfileResponse>> {
  return request.post('/redis-data/hash/field', data)
}

/** 查询完整 Hash 画像 */
export function getHashProfile(key: string): Promise<ApiResponse<HashProfileResponse>> {
  return request.get('/redis-data/hash/profile', { params: { key } })
}

/** 批量设置 Hash 字段 */
export function multiSetHash(data: HashMultiSetRequest): Promise<ApiResponse<HashProfileResponse>> {
  return request.post('/redis-data/hash/multi-set', data)
}

/** 删除 Hash 字段 */
export function deleteHashField(key: string, field: string): Promise<ApiResponse<HashProfileResponse>> {
  return request.post('/redis-data/hash/delete-field', { key, field })
}

// ===== Set: 标签/收藏夹管理 =====

/** 添加集合成员 */
export function setAddMembers(data: SetAddMembersRequest): Promise<ApiResponse<SetListResponse>> {
  return request.post('/redis-data/set/add', data)
}

/** 移除集合成员 */
export function setRemoveMembers(data: SetRemoveMembersRequest): Promise<ApiResponse<SetListResponse>> {
  return request.post('/redis-data/set/remove', data)
}

/** 查询集合所有成员 */
export function getSetMembers(key: string): Promise<ApiResponse<SetListResponse>> {
  return request.get('/redis-data/set/members', { params: { key } })
}

/** 计算交集 */
export function setIntersect(data: SetOperationRequest): Promise<ApiResponse<SetOperationResponse>> {
  return request.post('/redis-data/set/intersect', data)
}

/** 计算并集 */
export function setUnion(data: SetOperationRequest): Promise<ApiResponse<SetOperationResponse>> {
  return request.post('/redis-data/set/union', data)
}

/** 计算差集 */
export function setDiff(data: SetOperationRequest): Promise<ApiResponse<SetOperationResponse>> {
  return request.post('/redis-data/set/diff', data)
}

// ===== ZSet: 实时排行榜 =====

/** 添加/增加成员分数 */
export function zsetAddScore(data: ZSetAddScoreRequest): Promise<ApiResponse<ZSetMemberResponse>> {
  return request.post('/redis-data/zset/add-score', data)
}

/** 获取前 N 名 */
export function zsetTopN(data: ZSetTopNRequest): Promise<ApiResponse<ZSetTopNResponse>> {
  return request.post('/redis-data/zset/top-n', data)
}

/** 查询某成员排名 */
export function zsetGetRank(data: ZSetRankRequest): Promise<ApiResponse<ZSetMemberResponse>> {
  return request.post('/redis-data/zset/rank', data)
}

// ===== List: 最新活动流 / 简易消息队列 =====

/** 推入消息 */
export function listPush(data: ListPushRequest): Promise<ApiResponse<ListRangeResponse>> {
  return request.post('/redis-data/list/push', data)
}

/** 弹出消息 */
export function listPop(data: ListPopRequest): Promise<ApiResponse<ListPopResponse>> {
  return request.post('/redis-data/list/pop', data)
}

/** 查询列表范围 */
export function listRange(key: string, start = 0, stop = 19): Promise<ApiResponse<ListRangeResponse>> {
  return request.get('/redis-data/list/range', { params: { key, start, stop } })
}

// ===== String: 验证码存储 / 计数器 =====

/** 设置键值 */
export function stringSet(data: StringSetRequest): Promise<ApiResponse<StringSetResponse>> {
  return request.post('/redis-data/string/set', data)
}

/** 获取键值 */
export function stringGet(key: string): Promise<ApiResponse<StringGetResponse>> {
  return request.get('/redis-data/string/get', { params: { key } })
}

/** 自增计数 */
export function stringIncr(data: StringIncrRequest): Promise<ApiResponse<StringIncrResponse>> {
  return request.post('/redis-data/string/incr', data)
}
