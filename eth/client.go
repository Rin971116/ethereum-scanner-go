// eth/client.go
package eth

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 全域變數client 用來連接以太坊節點，寫在init()外面，代表全域，其他包也看的到、能用
var Client *ethclient.Client

func InitEthClient() {
	var err error
	Client, err = ethclient.Dial("https://mainnet.infura.io/v3/cf65d2dab655468e9b600b288726e837")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum RPC:", err)
	}
}

func CloseEthClient() {
	if Client != nil {
		Client.Close()
		log.Println("Ethereum client closed.")
	}
}
