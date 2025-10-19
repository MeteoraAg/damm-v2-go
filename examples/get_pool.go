package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/MeteoraAg/damm-v2-go/instructions"
)

func GetPool() {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	poolAddressStr := "YOUR_POOL_ADDRESS"

	fmt.Println("Getting pool...")
	poolAddress := solana.MustPublicKeyFromBase58(poolAddressStr)

	ctx := context.Background()

	pool, err := instructions.GetPool(ctx, poolAddress, rpcClient)
	if err != nil {
		log.Fatalf("Failed to get pool: %v", err)
	}

	jsonData, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal pool to JSON: %v", err)
	}

	fmt.Printf("Pool JSON: %s\n", string(jsonData))
}

// func main() {
// 	GetPool()
// }
