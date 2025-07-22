# Ethereum Scanner

高性能以太坊區塊鏈掃描器，實時掃描區塊和交易數據並提供 RESTful API。

## ✨ 功能

- 🔍 實時掃描以太坊區塊
- 💾 PostgreSQL 數據存儲
- 🚀 高並發交易處理
- 📡 RESTful API 查詢
- 📊 Kafka 消息推送
- 🔄 優雅關閉機制

## 🏗️ 架構

```
ethereum-scanner/
├── cmd/eth-scanner/     # 主程序
├── scanner/            # 區塊掃描器
├── controllers/        # API 控制器
├── domain/entity/     # 數據模型
├── database/          # 數據庫
├── eth/              # 以太坊客戶端
├── kafka/            # Kafka 消息
└── workerpool/       # 工作池
```

## 🚀 快速開始

### 1. 安裝依賴
```bash
go mod download
```

### 2. 配置環境
```bash
# 數據庫連接
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASS=postgres
export DB_NAME=ethereum_scanner

# 以太坊節點
export ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_API_KEY

# Kafka
export KAFKA_BROKERS=localhost:9092
```

### 3. 運行
```bash
go run cmd/eth-scanner/main.go
```

## 📡 API

| 端點 | 方法 | 描述 |
|------|------|------|
| `/block/:number` | GET | 獲取區塊信息 |
| `/transaction/:blockNumber` | GET | 獲取區塊交易 |
| `/migrate` | POST | 數據庫遷移 |

## 🛠️ 技術棧

- **Go 1.24.4** - 後端語言
- **Gin** - Web 框架
- **GORM** - ORM
- **PostgreSQL** - 數據庫
- **go-ethereum** - 以太坊客戶端
- **Kafka** - 消息隊列

## 📊 數據模型

### Block
```go
type Block struct {
    Number     uint64 `gorm:"uniqueIndex"`
    Hash       string `gorm:"size:66"`
    Timestamp  uint64
    GasLimit   uint64
    GasUsed    uint64
    TxCount    int
}
```

### Transaction
```go
type Transaction struct {
    Hash        string `gorm:"uniqueIndex"`
    BlockNumber uint64 `gorm:"index"`
    From        string `gorm:"type:varchar(42)"`
    To          string `gorm:"type:varchar(42)"`
    Value       string
    GasPrice    string
    GasUsed     uint64
    Status      bool
}
```

## ⚡ 性能特性

- **Worker Pool** - 並發處理交易
- **批量寫入** - 優化數據庫性能
- **限流機制** - 防止 RPC 過載
- **連接池** - 高效資源管理

## 🔧 配置

### Worker Pool
```go
workerpool.InitWorkerPool(20, &wg) // 20 個工作線程
```

### 掃描間隔
```go
time.Sleep(1 * time.Second) // 1 秒間隔
```

## 📝 日誌

應用會輸出詳細日誌：
- 區塊掃描進度
- 交易處理狀態
- 錯誤信息

## 🚨 注意事項

- 確保以太坊節點連接穩定
- 定期備份數據庫
- 監控系統資源使用

## 📄 License

MIT 