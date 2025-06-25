package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/dannwee/dbc-go/helpers"
	"github.com/dannwee/dbc-go/instructions"
)

func ClaimPositionFee() {
	ctx := context.Background()
	client := rpc.New("https://api.mainnet-beta.solana.com")

	// 1) load user keypair
	userKeypair := solana.MustPrivateKeyFromBase58("YOUR_PRIVATE_KEY")
	userWallet := userKeypair.PublicKey()

	// 2) pool address
	poolAddress := solana.MustPublicKeyFromBase58("YOUR_POOL_ADDRESS")

	// 3) get pool state
	poolState, err := instructions.GetPool(ctx, poolAddress, client)
	if err != nil {
		log.Fatalf("Failed to get pool state: %v", err)
	}

	// 4) get user positions for this pool
	positions, err := instructions.GetUserPositionByPool(ctx, client, poolAddress, userWallet)
	if err != nil {
		log.Fatalf("Failed to get user positions: %v", err)
	}

	if len(positions) == 0 {
		fmt.Println("No positions found for this user.")
		return
	}

	// 5) get position state for the first position
	positionState, err := instructions.GetPosition(ctx, positions[0].Position, client)
	if err != nil {
		log.Fatalf("Failed to get position state: %v", err)
	}

	// 6) get unclaimed rewards (TODO: check math for this)
	unclaimedReward, err := helpers.GetUnclaimReward(poolState, positionState)
	if err != nil {
		log.Fatalf("Failed to get unclaimed rewards: %v", err)
	}

	fmt.Printf("Fee Token A: %s\n", unclaimedReward.FeeTokenA.String())
	fmt.Printf("Fee Token B: %s\n", unclaimedReward.FeeTokenB.String())

	// 7) derive token accounts
	tokenAAccount, _, _ := solana.FindAssociatedTokenAddress(
		userWallet,
		poolState.TokenAMint,
	)
	tokenBAccount, _, _ := solana.FindAssociatedTokenAddress(
		userWallet,
		poolState.TokenBMint,
	)

	// 8) create ATAs if they dont exist
	var createTokenAAtaIx, createTokenBAtaIx solana.Instruction

	// check if token A ATA exists
	accountInfo, err := client.GetAccountInfo(ctx, tokenAAccount)
	if err != nil || accountInfo == nil || accountInfo.Value == nil {
		createTokenAAtaIx = associatedtokenaccount.NewCreateInstruction(
			userWallet,
			userWallet,
			poolState.TokenAMint,
		).Build()
	}

	// check if token B ATA exists
	accountInfo, err = client.GetAccountInfo(ctx, tokenBAccount)
	if err != nil || accountInfo == nil || accountInfo.Value == nil {
		createTokenBAtaIx = associatedtokenaccount.NewCreateInstruction(
			userWallet,
			userWallet,
			poolState.TokenBMint,
		).Build()
	}

	// 9) build claim position fee instruction
	ixClaim := instructions.ClaimPositionFee(
		poolAddress,
		positions[0].Position,
		tokenAAccount,
		tokenBAccount,
		poolState.TokenAVault,
		poolState.TokenBVault,
		poolState.TokenAMint,
		poolState.TokenBMint,
		solana.TokenProgramID,
		solana.TokenProgramID,
		positions[0].PositionNftAccount,
		userWallet,
	)

	// 10) assemble transaction
	bh, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("GetLatestBlockhash: %v", err)
	}

	// create instructions slice
	instructions := []solana.Instruction{ixClaim}

	// add ATA creation instructions only if needed
	if createTokenAAtaIx != nil {
		instructions = append([]solana.Instruction{createTokenAAtaIx}, instructions...)
	}
	if createTokenBAtaIx != nil {
		instructions = append([]solana.Instruction{createTokenBAtaIx}, instructions...)
	}

	tx, err := solana.NewTransaction(
		instructions,
		bh.Value.Blockhash,
		solana.TransactionPayer(userWallet),
	)
	if err != nil {
		log.Fatalf("NewTransaction: %v", err)
	}

	// 11) sign transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(userWallet) {
			return &userKeypair
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Sign: %v", err)
	}

	// 12) send & confirm
	sig, err := client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("SendTransaction: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", sig)

	// wait for confirmation by polling
	for i := 0; i < 30; i++ { // try for 30 secs
		time.Sleep(time.Second)
		resp, err := client.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
			Commitment: rpc.CommitmentFinalized,
		})
		if err != nil {
			continue
		}
		if resp != nil {
			if resp.Meta != nil && resp.Meta.Err != nil {
				log.Fatalf("Transaction failed: %v", resp.Meta.Err)
			}
			fmt.Printf("Transaction confirmed: %s\n", `https://solscan.io/tx/`+sig.String())
			return
		}
	}
	log.Fatalf("Transaction confirmation timeout")
}

// func main() {
// 	ClaimPositionFee()
// }
