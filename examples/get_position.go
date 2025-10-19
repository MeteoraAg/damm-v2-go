package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MeteoraAg/damm-v2-go/instructions"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetPosition() {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	userAddressStr := "YOUR_WALLET_ADDRESS"

	fmt.Println("Getting positions for user...")
	userAddress := solana.MustPublicKeyFromBase58(userAddressStr)

	ctx := context.Background()

	positions, err := instructions.GetPositionsByUser(ctx, rpcClient, userAddress)
	if err != nil {
		log.Fatalf("Failed to get positions: %v", err)
	}

	jsonData, err := json.MarshalIndent(positions, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal positions to JSON: %v", err)
	}

	fmt.Printf("Found %d positions\n", len(positions))
	fmt.Printf("Positions JSON: %s\n", string(jsonData))
}

// func main() {
// 	GetPosition()
// }
