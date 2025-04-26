package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

var infuraURL = "https://sepolia.infura.io/v3/xxx"

func main() {
	client, err := ethclient.DialContext(context.Background(), infuraURL)
	if err != nil {
		log.Fatalf("Error to create a ether client:%v", err)
	}
	defer client.Close()
	blockNumber := big.NewInt(8199285)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("Error to get a block:%v", err)
	}

	fmt.Printf("区块号: %d\n", block.Number())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("时间戳: %d (Unix 时间)\n", block.Time())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))

}
