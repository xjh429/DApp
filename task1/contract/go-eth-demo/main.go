package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/local/go-eth-demo/counter" // 导入生成的绑定代码
)

func main() {
	// 连接到 Sepolia
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/xxx")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// 设置私钥
	privateKey, err := crypto.HexToECDSA("xxx")
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// 从私钥获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// 获取链ID（Sepolia 的链ID是 11155111）
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	// 创建授权交易签名者
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // 不发送以太币
	auth.GasLimit = uint64(300000) // 设置合理的 gas limit
	auth.GasPrice, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	// 替换为已部署的合约地址
	contractAddress := common.HexToAddress("xxx")
	instance, err := counter.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	// 1. 首先获取当前计数器值
	currentCount, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get current count: %v", err)
	}
	fmt.Printf("Current count: %d\n", currentCount)

	// 2. 调用增加计数器的方法
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatalf("Failed to increment count: %v", err)
	}
	fmt.Printf("Increment transaction sent! Hash: %s\n", tx.Hash().Hex())

	// 等待交易被挖出
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction to be mined: %v", err)
	}
	if receipt.Status != 1 {
		log.Fatal("Transaction failed")
	}

	// 3. 再次获取计数器值查看变化
	updatedCount, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Failed to get updated count: %v", err)
	}
	fmt.Printf("Updated count: %d\n", updatedCount)
}
