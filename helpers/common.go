package helpers

import (
	"encoding/binary"
	"fmt"

	"github.com/MeteoraAg/damm-v2-go/common"
	"github.com/gagliardetto/solana-go"
	"lukechampine.com/uint128"
)

// Deserializes pool data
func DeserializePool(data []byte) (*common.Pool, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("data too short")
	}

	// Skip discriminator
	data = data[8:]

	pool := &common.Pool{}

	// Read PoolFees
	// BaseFee
	pool.PoolFees.BaseFee.CliffFeeNumerator = binary.LittleEndian.Uint64(data[0:8])
	pool.PoolFees.BaseFee.FeeSchedulerMode = data[8]
	copy(pool.PoolFees.BaseFee.Padding0[:], data[9:14])
	pool.PoolFees.BaseFee.NumberOfPeriod = binary.LittleEndian.Uint16(data[14:16])
	pool.PoolFees.BaseFee.PeriodFrequency = binary.LittleEndian.Uint64(data[16:24])
	pool.PoolFees.BaseFee.ReductionFactor = binary.LittleEndian.Uint64(data[24:32])
	pool.PoolFees.BaseFee.Padding1 = binary.LittleEndian.Uint64(data[32:40])
	data = data[40:]

	// ProtocolFeePercent, PartnerFeePercent, ReferralFeePercent
	pool.PoolFees.ProtocolFeePercent = data[0]
	pool.PoolFees.PartnerFeePercent = data[1]
	pool.PoolFees.ReferralFeePercent = data[2]
	data = data[3:]

	// Padding0
	copy(pool.PoolFees.Padding0[:], data[0:5])
	data = data[5:]

	// DynamicFee
	pool.PoolFees.DynamicFee.Initialized = data[0]
	copy(pool.PoolFees.DynamicFee.Padding[:], data[1:8])
	pool.PoolFees.DynamicFee.MaxVolatilityAccumulator = binary.LittleEndian.Uint32(data[8:12])
	pool.PoolFees.DynamicFee.VariableFeeControl = binary.LittleEndian.Uint32(data[12:16])
	pool.PoolFees.DynamicFee.BinStep = binary.LittleEndian.Uint16(data[16:18])
	pool.PoolFees.DynamicFee.FilterPeriod = binary.LittleEndian.Uint16(data[18:20])
	pool.PoolFees.DynamicFee.DecayPeriod = binary.LittleEndian.Uint16(data[20:22])
	pool.PoolFees.DynamicFee.ReductionFactor = binary.LittleEndian.Uint16(data[22:24])
	pool.PoolFees.DynamicFee.LastUpdateTimestamp = binary.LittleEndian.Uint64(data[24:32])
	pool.PoolFees.DynamicFee.BinStepU128 = uint128.From64(binary.LittleEndian.Uint64(data[32:40])).Add(uint128.From64(binary.LittleEndian.Uint64(data[40:48])).Lsh(64))
	pool.PoolFees.DynamicFee.SqrtPriceReference = uint128.From64(binary.LittleEndian.Uint64(data[48:56])).Add(uint128.From64(binary.LittleEndian.Uint64(data[56:64])).Lsh(64))
	pool.PoolFees.DynamicFee.VolatilityAccumulator = uint128.From64(binary.LittleEndian.Uint64(data[64:72])).Add(uint128.From64(binary.LittleEndian.Uint64(data[72:80])).Lsh(64))
	pool.PoolFees.DynamicFee.VolatilityReference = uint128.From64(binary.LittleEndian.Uint64(data[80:88])).Add(uint128.From64(binary.LittleEndian.Uint64(data[88:96])).Lsh(64))
	data = data[96:]

	// Padding1
	for i := 0; i < 2; i++ {
		pool.PoolFees.Padding1[i] = binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
	}
	data = data[16:]

	// Read PublicKeys (32 bytes each)
	pool.TokenAMint = solana.PublicKeyFromBytes(data[0:32])
	pool.TokenBMint = solana.PublicKeyFromBytes(data[32:64])
	pool.TokenAVault = solana.PublicKeyFromBytes(data[64:96])
	pool.TokenBVault = solana.PublicKeyFromBytes(data[96:128])
	pool.WhitelistedVault = solana.PublicKeyFromBytes(data[128:160])
	pool.Partner = solana.PublicKeyFromBytes(data[160:192])
	data = data[192:]

	// Read uint128 values
	pool.Liquidity = uint128.From64(binary.LittleEndian.Uint64(data[0:8])).Add(uint128.From64(binary.LittleEndian.Uint64(data[8:16])).Lsh(64))
	pool.Padding = uint128.From64(binary.LittleEndian.Uint64(data[16:24])).Add(uint128.From64(binary.LittleEndian.Uint64(data[24:32])).Lsh(64))
	data = data[32:]

	// Read uint64 values
	pool.ProtocolAFee = binary.LittleEndian.Uint64(data[0:8])
	pool.ProtocolBFee = binary.LittleEndian.Uint64(data[8:16])
	pool.PartnerAFee = binary.LittleEndian.Uint64(data[16:24])
	pool.PartnerBFee = binary.LittleEndian.Uint64(data[24:32])
	data = data[32:]

	// Read uint128 values for prices
	pool.SqrtMinPrice = uint128.From64(binary.LittleEndian.Uint64(data[0:8])).Add(uint128.From64(binary.LittleEndian.Uint64(data[8:16])).Lsh(64))
	pool.SqrtMaxPrice = uint128.From64(binary.LittleEndian.Uint64(data[16:24])).Add(uint128.From64(binary.LittleEndian.Uint64(data[24:32])).Lsh(64))
	pool.SqrtPrice = uint128.From64(binary.LittleEndian.Uint64(data[32:40])).Add(uint128.From64(binary.LittleEndian.Uint64(data[40:48])).Lsh(64))
	data = data[48:]

	// Read uint64 and uint8 values
	pool.ActivationPoint = binary.LittleEndian.Uint64(data[0:8])
	pool.ActivationType = data[8]
	pool.PoolStatus = data[9]
	pool.TokenAFlag = data[10]
	pool.TokenBFlag = data[11]
	pool.CollectFeeMode = data[12]
	pool.PoolType = data[13]
	data = data[14:]

	// Read padding0
	copy(pool.Padding0[:], data[0:2])
	data = data[2:]

	// Read fee per liquidity arrays
	copy(pool.FeeAPerLiquidity[:], data[0:32])
	copy(pool.FeeBPerLiquidity[:], data[32:64])
	data = data[64:]

	// Read permanent lock liquidity
	pool.PermanentLockLiquidity = uint128.From64(binary.LittleEndian.Uint64(data[0:8])).Add(uint128.From64(binary.LittleEndian.Uint64(data[8:16])).Lsh(64))
	data = data[16:]

	// Read metrics
	pool.Metrics.TotalLpAFee = uint128.From64(binary.LittleEndian.Uint64(data[0:8])).Add(uint128.From64(binary.LittleEndian.Uint64(data[8:16])).Lsh(64))
	pool.Metrics.TotalLpBFee = uint128.From64(binary.LittleEndian.Uint64(data[16:24])).Add(uint128.From64(binary.LittleEndian.Uint64(data[24:32])).Lsh(64))
	pool.Metrics.TotalProtocolAFee = binary.LittleEndian.Uint64(data[32:40])
	pool.Metrics.TotalProtocolBFee = binary.LittleEndian.Uint64(data[40:48])
	pool.Metrics.TotalPartnerAFee = binary.LittleEndian.Uint64(data[48:56])
	pool.Metrics.TotalPartnerBFee = binary.LittleEndian.Uint64(data[56:64])
	pool.Metrics.TotalPosition = binary.LittleEndian.Uint64(data[64:72])
	pool.Metrics.Padding = binary.LittleEndian.Uint64(data[72:80])
	data = data[80:]

	// Read padding1
	for i := 0; i < 10; i++ {
		pool.Padding1[i] = binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
	}
	data = data[80:]

	// Read reward infos
	for i := 0; i < 2; i++ {
		pool.RewardInfos[i].Initialized = data[0]
		pool.RewardInfos[i].RewardTokenFlag = data[1]
		copy(pool.RewardInfos[i].Padding0[:], data[2:8])
		copy(pool.RewardInfos[i].Padding1[:], data[8:16])
		pool.RewardInfos[i].Mint = solana.PublicKeyFromBytes(data[16:48])
		pool.RewardInfos[i].Vault = solana.PublicKeyFromBytes(data[48:80])
		pool.RewardInfos[i].Funder = solana.PublicKeyFromBytes(data[80:112])
		pool.RewardInfos[i].RewardDuration = binary.LittleEndian.Uint64(data[112:120])
		pool.RewardInfos[i].RewardDurationEnd = binary.LittleEndian.Uint64(data[120:128])
		pool.RewardInfos[i].RewardRate = uint128.From64(binary.LittleEndian.Uint64(data[128:136])).Add(uint128.From64(binary.LittleEndian.Uint64(data[136:144])).Lsh(64))
		copy(pool.RewardInfos[i].RewardPerTokenStored[:], data[144:176])
		pool.RewardInfos[i].LastUpdateTime = binary.LittleEndian.Uint64(data[176:184])
		pool.RewardInfos[i].CumulativeSecondsWithEmptyLiquidity = binary.LittleEndian.Uint64(data[184:192])
		data = data[192:]
	}

	return pool, nil
}

// Deserializes position data
func DeserializePosition(data []byte) (*common.PositionState, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("data too short")
	}

	// Skip discriminator
	data = data[8:]

	position := &common.PositionState{}

	// Read public keys
	position.Pool = solana.PublicKeyFromBytes(data[0:32])
	position.NftMint = solana.PublicKeyFromBytes(data[32:64])
	data = data[64:]

	// Read fee checkpoints
	copy(position.FeeAPerTokenCheckpoint[:], data[0:32])
	copy(position.FeeBPerTokenCheckpoint[:], data[32:64])
	data = data[64:]

	// Read fee pending
	position.FeeAPending = binary.LittleEndian.Uint64(data[0:8])
	position.FeeBPending = binary.LittleEndian.Uint64(data[8:16])
	data = data[16:]

	// Read liquidity values
	position.UnlockedLiquidity = uint128.From64(binary.LittleEndian.Uint64(data[0:8])).Add(uint128.From64(binary.LittleEndian.Uint64(data[8:16])).Lsh(64))
	position.VestedLiquidity = uint128.From64(binary.LittleEndian.Uint64(data[16:24])).Add(uint128.From64(binary.LittleEndian.Uint64(data[24:32])).Lsh(64))
	position.PermanentLockedLiquidity = uint128.From64(binary.LittleEndian.Uint64(data[32:40])).Add(uint128.From64(binary.LittleEndian.Uint64(data[40:48])).Lsh(64))
	data = data[48:]

	// Read metrics
	position.Metrics.TotalClaimedAFee = binary.LittleEndian.Uint64(data[0:8])
	position.Metrics.TotalClaimedBFee = binary.LittleEndian.Uint64(data[8:16])
	data = data[16:]

	// Read reward infos
	for i := 0; i < 2; i++ {
		copy(position.RewardInfos[i].RewardPerTokenCheckpoint[:], data[0:32])
		position.RewardInfos[i].RewardPendings = binary.LittleEndian.Uint64(data[32:40])
		position.RewardInfos[i].TotalClaimedRewards = binary.LittleEndian.Uint64(data[40:48])
		data = data[48:]
	}

	// Read padding
	for i := 0; i < 6; i++ {
		position.Padding[i] = uint128.From64(binary.LittleEndian.Uint64(data[i*16 : i*16+8])).Add(uint128.From64(binary.LittleEndian.Uint64(data[i*16+8 : i*16+16])).Lsh(64))
	}

	return position, nil
}
