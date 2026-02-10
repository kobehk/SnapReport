# SnapReport

SnapReport 是一个后端服务，用于抓取行车记录仪视频并生成基于位置的报告。它集成了盯盯拍 (DDPAI) 行车记录仪以抓取视频片段，并使用 Nominatim (OpenStreetMap) 进行逆地理编码，以识别位置和道路类型。

## 功能特性

- **行车记录仪集成**：从连接的盯盯拍设备抓取最近的视频片段。
- **地理编码**：根据 GPS 坐标自动确定城市和道路名称。
- **高速公路检测**：分类当前位置是否位于高速公路/快速路上。
- **报告管理**：提供 API 用于准备、发送和列出报告。
- **模拟模式**：支持在没有物理行车记录仪的情况下进行开发。

## 项目结构

```
SnapReport/
├── config.yaml          # 配置文件
├── internal/
│   ├── api/             # HTTP API 处理程序
│   ├── config/          # 配置加载
│   ├── ddpai/           # DDPAI 设备客户端
│   ├── geo/             # 地理编码和高速公路分类
│   ├── model/           # 数据模型
│   ├── service/         # 业务逻辑
│   └── store/           # 内存数据存储
└── main.go              # 入口点
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本

### 配置

编辑 `config.yaml` 以配置服务器、盯盯拍设备连接和地理编码设置。

```yaml
server:
  port: 8081

ddpai:
  base_url: "http://193.168.0.1"
  timeout_seconds: 5
  mock_mode: true # 设置为 true 以模拟设备连接

nominatim:
  user_agent: "SnapReport/1.0"
```

### 运行应用

```bash
go run main.go
```

服务器将在 `config.yaml` 中指定的端口上启动（默认为 8081）。

## API 接口

### 1. 健康检查 (Health Check)
检查服务是否正在运行。

- **URL**: `/health`
- **Method**: `GET`
- **Response**: `{"status": "ok"}`
- **Example**:
  ```bash
  curl http://localhost:8081/health
  ```

### 2. 准备报告 (Prepare Report)
抓取视频并创建包含位置数据的初步报告。

- **URL**: `/reports/prepare`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "device_id": "device_123",
    "lat": 39.9042,
    "lng": 116.4074,
    "duration_sec": 20,
    "tags": ["traffic", "accident"]
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:8081/reports/prepare \
    -H "Content-Type: application/json" \
    -d '{
      "device_id": "device_123",
      "lat": 39.9042,
      "lng": 116.4074,
      "duration_sec": 20,
      "tags": ["traffic", "accident"]
    }'
  ```
- **Response**:
  ```json
  {
    "id": "rep_...",
    "timestamp": "2023-10-27T10:00:00Z",
    "lat": 39.9042,
    "lng": 116.4074,
    "city": "北京市",
    "road_name": "长安街",
    "is_highway": false,
    "video_url": "http://...",
    "status": "prepared",
    "device_id": "device_123"
  }
  ```

### 3. 发送报告 (Send Report)
将报告标记为已提交。

- **URL**: `/reports/send`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "id": "rep_..."
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:8081/reports/send \
    -H "Content-Type: application/json" \
    -d '{
      "id": "rep_..."
    }'
  ```

### 4. 获取报告列表 (List Reports)
获取所有存储的报告。

- **URL**: `/reports`
- **Method**: `GET`
- **Example**:
  ```bash
  curl http://localhost:8081/reports
  ```

## 许可证

[MIT](LICENSE)
