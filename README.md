# Ethereum Scanner

一個高性能的以太坊區塊鏈掃描器，用 Go 語言開發，能夠實時掃描以太坊主網的區塊和交易，並提供 RESTful API 服務。

## 🚀 功能特色

- **實時區塊掃描**: 自動掃描最新的以太坊區塊
- **交易處理**: 並發處理區塊中的所有交易
- **數據持久化**: 使用 PostgreSQL 存儲區塊和交易數據
- **RESTful API**: 提供區塊和交易查詢 API
- **高並發處理**: 使用 Worker Pool 模式處理交易
- **優雅關閉**: 支持優雅關閉和資源清理
- **自動遷移**: 數據庫結構自動遷移

## 🏗️ 系統架構

```
ethereum-scanner/
├── cmd/eth-scanner/          # 主程序入口
├── scanner/                  # 區塊掃描器
├── controllers/              # API 控制器
│   └── handler/             # API 處理器
├── domain/entity/           # 數據實體
├── database/                # 數據庫連接
├── eth/                     # 以太坊客戶端
├── workerpool/              # 工作池
├── repositories/            # 數據訪問層
├── usecase/                 # 業務邏輯層
└── test/                    # 測試文件
```

## 📋 技術棧

- **語言**: Go 1.24.4
- **Web 框架**: Gin
- **ORM**: GORM
- **數據庫**: PostgreSQL
- **以太坊**: go-ethereum
- **並發**: Goroutines + Channels + WaitGroup

## 🛠️ 安裝與配置

### 前置要求

- Go 1.24.4 或更高版本
- PostgreSQL 數據庫
- 以太坊節點訪問權限（或 Infura API Key）

### 1. 克隆專案

```bash
git clone <repository-url>
cd ethereum-scanner
```

### 2. 安裝依賴

```bash
go mod download
```

### 3. 配置數據庫

修改 `database/postgres.go` 中的數據庫連接字符串：

```go
dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
```

### 4. 配置以太坊節點

修改 `eth/client.go` 中的以太坊節點 URL：

```go
Client, err = ethclient.Dial("https://mainnet.infura.io/v3/YOUR_API_KEY")
```

### 5. 運行應用

```bash
go run cmd/eth-scanner/main.go
```

## 📊 數據模型

### Block 實體

```go
type Block struct {
    gorm.Model
    Number     uint64 `gorm:"uniqueIndex;not null"`
    Hash       string `gorm:"size:66;not null"`
    ParentHash string `gorm:"size:66"`
    Timestamp  uint64 `gorm:"not null"`
    Miner      string `gorm:"size:42"`
    GasLimit   uint64
    GasUsed    uint64
    BaseFee    string `gorm:"type:text"`
    TxCount    int    `gorm:"not null"`
}
```

### Transaction 實體

```go
type Transaction struct {
    gorm.Model
    Hash        string `gorm:"uniqueIndex;type:varchar(66)"`
    BlockNumber uint64 `gorm:"index"`
    From        string `gorm:"type:varchar(42)"`
    To          string `gorm:"type:varchar(42)"`
    Value       string `gorm:"type:varchar(32)"`
    GasPrice    string `gorm:"type:varchar(32)"`
    GasLimit    uint64
    GasUsed     uint64
    Nonce       uint64
    Data        []byte `gorm:"type:bytea"`
    Status      bool
}
```

## 🔌 API 文檔

### 獲取區塊信息

**GET** `/block/:blockNumber`

獲取指定區塊號的區塊信息。

**參數:**
- `blockNumber` (path): 區塊號碼

**響應示例:**
```json
{
    "ID": 1,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z",
    "DeletedAt": null,
    "Number": 19000000,
    "Hash": "0x1234567890abcdef...",
    "ParentHash": "0xabcdef1234567890...",
    "Timestamp": 1704067200,
    "Miner": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
    "GasLimit": 30000000,
    "GasUsed": 15000000,
    "BaseFee": "20000000000",
    "TxCount": 150
}
```

### 獲取區塊交易

**GET** `/transaction/:blockNumber`

獲取指定區塊號的所有交易。

**參數:**
- `blockNumber` (path): 區塊號碼

**響應示例:**
```json
[
    {
        "ID": 1,
        "CreatedAt": "2024-01-01T00:00:00Z",
        "UpdatedAt": "2024-01-01T00:00:00Z",
        "DeletedAt": null,
        "Hash": "0xabcdef1234567890...",
        "BlockNumber": 19000000,
        "From": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
        "To": "0x1234567890abcdef...",
        "Value": "1000000000000000000",
        "GasPrice": "20000000000",
        "GasLimit": 21000,
        "GasUsed": 21000,
        "Nonce": 0,
        "Data": [],
        "Status": true
    }
]
```

### 數據庫遷移

**POST** `/migrate`

執行數據庫結構遷移。

**響應示例:**
```json
{
    "message": "Database migration completed"
}
```

### 事件日誌

**GET** `/eventlog`

獲取事件日誌（待實現）。

## ⚡ 性能優化

### 並發處理

- 使用 Worker Pool 模式並發處理交易
- 可配置工作線程數量（默認 10 個）
- 使用 Channel 進行任務分發

### 批量處理

- 交易批量寫入數據庫
- 減少數據庫連接開銷
- 提高寫入性能

### 優雅關閉

- 支持 SIGINT 和 SIGTERM 信號
- 等待所有 goroutine 完成後關閉
- 正確釋放數據庫和以太坊客戶端資源

## 🔧 配置選項

### Worker Pool 配置

在 `main.go` 中修改工作線程數量：

```go
workerpool.InitWorkerPool(10, &wg) // 10 個工作線程
```

### 掃描間隔

在 `scanner/scanner.go` 中修改掃描間隔：

```go
time.Sleep(1 * time.Second) // 1 秒間隔
```

### API 端口

在 `controllers/ginServer.go` 中修改服務端口：

```go
Addr: ":8080" // 默認 8080 端口
```

## 🧪 測試

運行測試：

```bash
go test ./...
```

## 📝 日誌

應用會輸出詳細的日誌信息，包括：

- 區塊掃描進度
- 交易處理狀態
- 數據庫操作結果
- 錯誤信息

## 🚨 注意事項

1. **API Key 安全**: 請妥善保管以太坊節點的 API Key
2. **數據庫備份**: 建議定期備份 PostgreSQL 數據庫
3. **資源監控**: 監控內存和 CPU 使用情況
4. **網絡連接**: 確保以太坊節點連接穩定

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 許可證

MIT License

## 📞 聯繫

如有問題，請提交 Issue 或聯繫開發團隊。 