package instructions

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/MeteoraAg/damm-v2-go/common"
	"github.com/MeteoraAg/damm-v2-go/helpers"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetPool(ctx context.Context, poolAddress solana.PublicKey, rpcClient *rpc.Client) (*common.Pool, error) {
	account, err := rpcClient.GetAccountInfo(ctx, poolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool account: %w", err)
	}

	if account == nil || account.Value == nil {
		return nil, fmt.Errorf("pool account not found")
	}

	data := account.Value.Data.GetBinary()

	if len(data) < 8 {
		return nil, fmt.Errorf("data too short")
	}

	expectedDiscriminator := []byte{241, 154, 109, 4, 17, 177, 109, 188}
	if !bytes.Equal(data[:8], expectedDiscriminator) {
		return nil, fmt.Errorf("invalid discriminator, not a pool account")
	}

	return helpers.DeserializePool(data)
}

func GetPositionsByUser(
	ctx context.Context,
	rpcClient *rpc.Client,
	user solana.PublicKey,
) ([]common.PositionResult, error) {
	// Get all position NFT accounts owned by the user
	userPositionAccounts, err := common.GetAllPositionNftAccountByOwner(ctx, rpcClient, user)
	if err != nil {
		return nil, fmt.Errorf("failed to get position NFT accounts: %w", err)
	}

	if len(userPositionAccounts) == 0 {
		return []common.PositionResult{}, nil
	}

	// Get position addresses for each NFT
	positionAddresses := make([]solana.PublicKey, len(userPositionAccounts))
	for i, account := range userPositionAccounts {
		positionAddress, err := helpers.DerivePositionPDA(account.PositionNft)
		if err != nil {
			return nil, fmt.Errorf("failed to derive position address: %w", err)
		}
		positionAddresses[i] = positionAddress
	}

	// Fetch all position states
	positionStates := make([]*common.PositionState, len(positionAddresses))
	for i, address := range positionAddresses {
		account, err := rpcClient.GetAccountInfo(ctx, address)
		if err != nil {
			return nil, fmt.Errorf("failed to get position account: %w", err)
		}
		if account == nil || account.Value == nil {
			continue
		}

		positionState, err := helpers.DeserializePosition(account.Value.Data.GetBinary())
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize position state: %w", err)
		}
		positionStates[i] = positionState
	}

	// Create position results
	positionResults := make([]common.PositionResult, 0, len(userPositionAccounts))
	for i, account := range userPositionAccounts {
		if positionStates[i] == nil {
			continue
		}

		positionResults = append(positionResults, common.PositionResult{
			PositionNftAccount: account.PositionNftAccount,
			Position:           positionAddresses[i],
			PositionState:      *positionStates[i],
		})
	}

	// Sort positions by total liquidity
	sort.Slice(positionResults, func(i, j int) bool {
		totalLiquidityI := positionResults[i].PositionState.VestedLiquidity.
			Add(positionResults[i].PositionState.PermanentLockedLiquidity).
			Add(positionResults[i].PositionState.UnlockedLiquidity)

		totalLiquidityJ := positionResults[j].PositionState.VestedLiquidity.
			Add(positionResults[j].PositionState.PermanentLockedLiquidity).
			Add(positionResults[j].PositionState.UnlockedLiquidity)

		return totalLiquidityJ.Cmp(totalLiquidityI) < 0
	})

	return positionResults, nil
}

func GetPosition(ctx context.Context, positionAddress solana.PublicKey, rpcClient *rpc.Client) (*common.PositionState, error) {
	account, err := rpcClient.GetAccountInfo(ctx, positionAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get position account: %w", err)
	}

	if account == nil || account.Value == nil {
		return nil, fmt.Errorf("position account not found")
	}

	data := account.Value.Data.GetBinary()

	if len(data) < 8 {
		return nil, fmt.Errorf("data too short")
	}

	expectedDiscriminator := []byte{170, 188, 143, 228, 122, 64, 247, 208}
	if !bytes.Equal(data[:8], expectedDiscriminator) {
		return nil, fmt.Errorf("invalid discriminator, not a position account")
	}

	return helpers.DeserializePosition(data)
}

func GetUserPositionByPool(
	ctx context.Context,
	rpcClient *rpc.Client,
	pool solana.PublicKey,
	user solana.PublicKey,
) ([]common.PositionResult, error) {
	// Get all positions for the user
	allPositions, err := GetPositionsByUser(ctx, rpcClient, user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user positions: %w", err)
	}

	// Filter positions by pool
	filteredPositions := make([]common.PositionResult, 0)
	for _, position := range allPositions {
		if position.PositionState.Pool.Equals(pool) {
			filteredPositions = append(filteredPositions, position)
		}
	}

	return filteredPositions, nil
}
