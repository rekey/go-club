# API 接口文档

## 基础信息

- **服务地址**: `http://localhost:8888`
- **基础路径**: `/api`

## 接口列表

### 1. 获取所有任务

**接口**: `GET /api/task/all`

**描述**: 获取系统中所有的任务列表

**请求参数**: 无

**响应示例**:

成功响应 (200):
```json
[
  {
    "ID": 1,
    "CreatedAt": "2025-12-30T10:00:00Z",
    "UpdatedAt": "2025-12-30T10:05:00Z",
    "DeletedAt": null,
    "url": "https://example.com/video.mp4",
    "name": "video.mp4",
    "progress": 50.5,
    "err": "",
    "status": 2
  }
]
```

失败响应 (200):
```json
{
  "code": 3,
  "msg": "get tasks error"
}
```

---

### 2. 添加任务

**接口**: `GET /api/task/add`

**描述**: 根据URL创建一个新的下载任务。如果URL已存在则返回已有任务，如果URL格式不支持则返回错误。

**请求参数**:

| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| url | string | 是 | 要下载的资源URL |

**请求示例**:
```
GET /api/task/add?url=https://example.com/video.mp4
```

**响应示例**:

成功响应 (200):
```json
{
  "ID": 1,
  "CreatedAt": "2025-12-30T10:00:00Z",
  "UpdatedAt": "2025-12-30T10:00:00Z",
  "DeletedAt": null,
  "url": "https://example.com/video.mp4",
  "name": "video.mp4",
  "progress": 0,
  "err": "",
  "status": 0
}
```

失败响应 (200):
```json
{
  "code": 1,
  "msg": "url not support"
}
```

---

### 3. 获取单个任务

**接口**: `GET /api/task/item`

**描述**: 根据URL获取单个任务详情

**请求参数**:

| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| url | string | 是 | 任务的URL |

**请求示例**:
```
GET /api/task/item?url=https://example.com/video.mp4
```

**响应示例**:

成功响应 (200):
```json
{
  "ID": 1,
  "CreatedAt": "2025-12-30T10:00:00Z",
  "UpdatedAt": "2025-12-30T10:05:00Z",
  "DeletedAt": null,
  "url": "https://example.com/video.mp4",
  "name": "video.mp4",
  "progress": 25,
  "err": "",
  "status": 2
}
```

失败响应 (200):
```json
{
  "code": 2,
  "msg": "task not found"
}
```

---

### 4. 启动任务

**接口**: `GET /api/task/start`

**描述**: 根据URL查找任务并强制启动

**请求参数**:

| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| url | string | 是 | 任务的URL |

**请求示例**:
```
GET /api/task/start?url=https://example.com/video.mp4
```

**响应示例**:

成功响应 (200):
```json
{
  "ID": 1,
  "CreatedAt": "2025-12-30T10:00:00Z",
  "UpdatedAt": "2025-12-30T10:05:00Z",
  "DeletedAt": null,
  "url": "https://example.com/video.mp4",
  "name": "video.mp4",
  "progress": 25,
  "err": "",
  "status": 2
}
```

失败响应 (200):
```json
{
  "code": 2,
  "msg": "task not found"
}
```

---

## 数据模型

### Task 对象

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | uint | 任务ID（GORM自动生成） |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |
| DeletedAt | datetime | 删除时间 |
| url | string | 任务URL（主键） |
| name | string | 任务名称 |
| progress | float64 | 下载进度（0-100） |
| err | string | 错误信息 |
| status | int | 任务状态 |

### 任务状态值 (Status)

| 状态值 | 常量名 | 描述 |
|--------|--------|------|
| -1 | Err | 错误 |
| 0 | Wait | 等待中 |
| 1 | Parser | 解析中 |
| 2 | Download | 下载中 |
| 3 | Merge | 合并中 |
| 4 | Move | 移动中 |
| 5 | Complete | 完成 |

---

## 响应状态码说明

| HTTP状态码 | 描述 |
|------------|------|
| 200 | 请求成功 |

## 业务错误码说明

| code | msg | 描述 |
|------|-----|------|
| 1 | url not support | 不支持的URL格式 |
| 2 | task not found | 任务未找到 |
| 3 | get tasks error | 获取任务列表失败 |
| 4 | start task failed | 启动任务失败 |

## 注意事项

1. 所有接口均使用GET方法
2. 服务监听端口为 8888
3. 任务URL通过查询参数传递
4. 所有响应的HTTP状态码均为 200，通过业务code区分成功失败
5. Task对象的url字段为主键，同一URL只能有一个任务
