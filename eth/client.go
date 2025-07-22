// eth/client.go
package eth

import (
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 全域變數client 用來連接以太坊節點，寫在init()外面，代表全域，其他包也看的到、能用
var Clients []*ethclient.Client
var currentIndex int
var mutex sync.Mutex

func InitEthClient() {
	currentIndex = 0

	urls := []string{
		"https://mainnet.infura.io/v3/cf65d2dab655468e9b600b288726e837",
		"https://eth-mainnet.g.alchemy.com/v2/QKG7rwN9_PsMIFVkJORHP2qSPOuNjOc4",
		"https://rpc.ankr.com/eth/651d3be568769c131b20ca3fb4752ed393cda5dde02e5f49d060e384e5c2ae11",
	}

	for i, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			log.Printf("RPC %d 連線失敗: %v", i, err)
			continue
		}
		log.Printf("RPC %d 連線成功: %s", i, url)
		Clients = append(Clients, client)
	}

	if len(Clients) == 0 {
		log.Fatal("無可用 RPC")
	}
}

func CloseEthClient() {
	for i, client := range Clients {
		if client != nil {
			client.Close()
		}
		log.Printf("Ethereum client %d closed.", i)
	}
}

func GetNextEthClient() *ethclient.Client {
	mutex.Lock()
	defer mutex.Unlock()

	client := Clients[currentIndex]
	currentIndex = (currentIndex + 1) % len(Clients)

	return client
}
