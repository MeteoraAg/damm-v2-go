package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dannwee/dbc-go/helpers"
	"github.com/dannwee/dbc-go/instructions"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetUnclaimReward() {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	poolAddressStr := "YOUR_POOL_ADDRESS"
	userAddressStr := "YOUR_WALLET_ADDRESS"

	fmt.Println("Getting pool state...")
	poolAddress := solana.MustPublicKeyFromBase58(poolAddressStr)
	userAddress := solana.MustPublicKeyFromBase58(userAddressStr)

	ctx := context.Background()

	// get pool state
	poolState, err := instructions.GetPool(ctx, poolAddress, rpcClient)
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

	// get unclaimed rewards
	unclaimedReward, err := helpers.GetUnclaimReward(poolState, positionState)
	if err != nil {
		log.Fatalf("Failed to get unclaimed rewards: %v", err)
	}

	positionJSON, err := json.MarshalIndent(positionState, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal position state to JSON: %v", err)
	}

	unclaimedJSON, err := json.MarshalIndent(unclaimedReward, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal unclaimed rewards to JSON: %v", err)
	}

	fmt.Printf("Found %d positions for this pool\n", len(positions))
	fmt.Printf("Position state: %s\n", string(positionJSON))
	fmt.Printf("Unclaimed rewards: %s\n", string(unclaimedJSON))

	fmt.Printf("Fee Token A: %s\n", unclaimedReward.FeeTokenA.String())
	fmt.Printf("Fee Token B: %s\n", unclaimedReward.FeeTokenB.String())
	fmt.Printf("Rewards: %s\n", unclaimedReward.Rewards)
}

func main() {
	GetUnclaimReward()
}
