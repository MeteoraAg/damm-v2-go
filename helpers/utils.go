package helpers

import (
	"math/big"

	"github.com/MeteoraAg/damm-v2-go/common"
	"lukechampine.com/uint128"
)

// reverseBytes reverses a byte slice to match TypeScript Buffer.reverse() behavior
func reverseBytes(b []byte) []byte {
	reversed := make([]byte, len(b))
	for i, j := 0, len(b)-1; i < len(b); i, j = i+1, j-1 {
		reversed[i] = b[j]
	}
	return reversed
}

// GetUnclaimReward calculates the unclaimed rewards for a position
func GetUnclaimReward(poolState *common.Pool, positionState *common.PositionState) (*common.UnclaimReward, error) {
	// Calculate total position liquidity
	totalPositionLiquidity := positionState.UnlockedLiquidity.
		Add(positionState.VestedLiquidity).
		Add(positionState.PermanentLockedLiquidity)

	// Convert byte arrays to uint128 with proper byte order
	feeAPerLiquidity := uint128.FromBytesBE(reverseBytes(poolState.FeeAPerLiquidity[:16]))
	feeBPerLiquidity := uint128.FromBytesBE(reverseBytes(poolState.FeeBPerLiquidity[:16]))
	feeAPerTokenCheckpoint := uint128.FromBytesBE(reverseBytes(positionState.FeeAPerTokenCheckpoint[:16]))
	feeBPerTokenCheckpoint := uint128.FromBytesBE(reverseBytes(positionState.FeeBPerTokenCheckpoint[:16]))

	// Calculate fee per token stored for token A
	feeAPerTokenStored := feeAPerLiquidity.Sub(feeAPerTokenCheckpoint)

	// Calculate fee per token stored for token B
	feeBPerTokenStored := feeBPerLiquidity.Sub(feeBPerTokenCheckpoint)

	// Convert to big.Int for calculation
	totalLiquidityBig := totalPositionLiquidity.Big()
	feeAPerTokenStoredBig := feeAPerTokenStored.Big()
	feeBPerTokenStoredBig := feeBPerTokenStored.Big()

	// Calculate fees: (totalLiquidity * feePerTokenStored) >> LIQUIDITY_SCALE
	feeABig := new(big.Int).Mul(totalLiquidityBig, feeAPerTokenStoredBig)
	feeABig.Rsh(feeABig, common.LIQUIDITY_SCALE)

	feeBBig := new(big.Int).Mul(totalLiquidityBig, feeBPerTokenStoredBig)
	feeBBig.Rsh(feeBBig, common.LIQUIDITY_SCALE)

	// Convert back to uint128
	feeA := uint128.FromBig(feeABig)
	feeB := uint128.FromBig(feeBBig)

	// Convert pending fees to uint128
	feeAPending := uint128.From64(positionState.FeeAPending)
	feeBPending := uint128.From64(positionState.FeeBPending)

	// Calculate total fees including pending
	totalFeeA := feeAPending.Add(feeA)
	totalFeeB := feeBPending.Add(feeB)

	// Get rewards from reward infos
	rewards := make([]uint128.Uint128, 0)
	if len(positionState.RewardInfos) > 0 {
		for _, info := range positionState.RewardInfos {
			rewards = append(rewards, uint128.From64(info.RewardPendings))
		}
	}

	return &common.UnclaimReward{
		FeeTokenA: totalFeeA,
		FeeTokenB: totalFeeB,
		Rewards:   rewards,
	}, nil
}
