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

func GetUserPositionByPool() {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	poolAddressStr := "YOUR_POOL_ADDRESS"
	userAddressStr := "YOUR_WALLET_ADDRESS"

	fmt.Println("Getting pool state...")
	poolAddress := solana.MustPublicKeyFromBase58(poolAddressStr)
	userAddress := solana.MustPublicKeyFromBase58(userAddressStr)

	ctx := context.Background()

	// get pool state
	_, err := instructions.GetPool(ctx, poolAddress, rpcClient)
	if err != nil {
		log.Fatalf("Failed to get pool state: %v", err)
	}

	// get user positions for this pool
	positions, err := instructions.GetUserPositionByPool(ctx, rpcClient, poolAddress, userAddress)
	if err != nil {
		log.Fatalf("Failed to get user positions: %v", err)
	}

	if len(positions) == 0 {
		fmt.Println("No positions found for this user.")
		return
	}

	// get position state for the first position
	positionState, err := instructions.GetPosition(ctx, positions[0].Position, rpcClient)
	if err != nil {
		log.Fatalf("Failed to get position state: %v", err)
	}

	jsonData, err := json.MarshalIndent(positionState, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal position state to JSON: %v", err)
	}

	fmt.Printf("Found %d positions for this pool\n", len(positions))
	fmt.Printf("Position state: %s\n", string(jsonData))
}

//func main() {
//	GetUserPositionByPool()
//}
