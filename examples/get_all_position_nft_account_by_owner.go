package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/MeteoraAg/damm-v2-go/common"
)

func GetAllPositionNftAccountByOwner() {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	ownerAddressStr := "YOUR_WALLET_ADDRESS"

	fmt.Println("Getting position NFT accounts...")
	ownerAddress := solana.MustPublicKeyFromBase58(ownerAddressStr)

	ctx := context.Background()

	positionNftAccounts, err := common.GetAllPositionNftAccountByOwner(ctx, rpcClient, ownerAddress)
	if err != nil {
		log.Fatalf("Failed to get position NFT accounts: %v", err)
	}

	jsonData, err := json.MarshalIndent(positionNftAccounts, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal position NFT accounts to JSON: %v", err)
	}

	fmt.Printf("Found %d position NFT accounts\n", len(positionNftAccounts))
	fmt.Printf("Position NFT Accounts JSON: %s\n", string(jsonData))
}

// func main() {
// 	GetAllPositionNftAccountByOwner()
// }
