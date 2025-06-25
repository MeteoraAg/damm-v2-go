package common

import (
	"github.com/gagliardetto/solana-go"
	"lukechampine.com/uint128"
)

type BaseFeeStruct struct {
	CliffFeeNumerator uint64
	FeeSchedulerMode  uint8
	Padding0          [5]uint8
	NumberOfPeriod    uint16
	PeriodFrequency   uint64
	ReductionFactor   uint64
	Padding1          uint64
}

type DynamicFeeStruct struct {
	Initialized              uint8
	Padding                  [7]uint8
	MaxVolatilityAccumulator uint32
	VariableFeeControl       uint32
	BinStep                  uint16
	FilterPeriod             uint16
	DecayPeriod              uint16
	ReductionFactor          uint16
	LastUpdateTimestamp      uint64
	BinStepU128              uint128.Uint128
	SqrtPriceReference       uint128.Uint128
	VolatilityAccumulator    uint128.Uint128
	VolatilityReference      uint128.Uint128
}

type PoolFeesStruct struct {
	BaseFee            BaseFeeStruct
	ProtocolFeePercent uint8
	PartnerFeePercent  uint8
	ReferralFeePercent uint8
	Padding0           [5]uint8
	DynamicFee         DynamicFeeStruct
	Padding1           [2]uint64
}

type PoolMetrics struct {
	TotalLpAFee       uint128.Uint128
	TotalLpBFee       uint128.Uint128
	TotalProtocolAFee uint64
	TotalProtocolBFee uint64
	TotalPartnerAFee  uint64
	TotalPartnerBFee  uint64
	TotalPosition     uint64
	Padding           uint64
}

type RewardInfo struct {
	Initialized                         uint8
	RewardTokenFlag                     uint8
	Padding0                            [6]uint8
	Padding1                            [8]uint8
	Mint                                solana.PublicKey
	Vault                               solana.PublicKey
	Funder                              solana.PublicKey
	RewardDuration                      uint64
	RewardDurationEnd                   uint64
	RewardRate                          uint128.Uint128
	RewardPerTokenStored                [32]uint8
	LastUpdateTime                      uint64
	CumulativeSecondsWithEmptyLiquidity uint64
}

type Pool struct {
	PoolFees               PoolFeesStruct
	TokenAMint             solana.PublicKey
	TokenBMint             solana.PublicKey
	TokenAVault            solana.PublicKey
	TokenBVault            solana.PublicKey
	WhitelistedVault       solana.PublicKey
	Partner                solana.PublicKey
	Liquidity              uint128.Uint128
	Padding                uint128.Uint128
	ProtocolAFee           uint64
	ProtocolBFee           uint64
	PartnerAFee            uint64
	PartnerBFee            uint64
	SqrtMinPrice           uint128.Uint128
	SqrtMaxPrice           uint128.Uint128
	SqrtPrice              uint128.Uint128
	ActivationPoint        uint64
	ActivationType         uint8
	PoolStatus             uint8
	TokenAFlag             uint8
	TokenBFlag             uint8
	CollectFeeMode         uint8
	PoolType               uint8
	Padding0               [2]uint8
	FeeAPerLiquidity       [32]uint8
	FeeBPerLiquidity       [32]uint8
	PermanentLockLiquidity uint128.Uint128
	Metrics                PoolMetrics
	Padding1               [10]uint64
	RewardInfos            [2]RewardInfo
}

type PositionNftAccount struct {
	PositionNft        solana.PublicKey
	PositionNftAccount solana.PublicKey
}

type TokenAccountLayout struct {
	Mint            solana.PublicKey
	Owner           solana.PublicKey
	Amount          uint64
	Delegate        solana.PublicKey
	State           uint8
	IsNative        uint64
	DelegatedAmount uint64
	CloseAuthority  solana.PublicKey
}

type PositionMetrics struct {
	TotalClaimedAFee uint64
	TotalClaimedBFee uint64
}

type UserRewardInfo struct {
	RewardPerTokenCheckpoint [32]uint8
	RewardPendings           uint64
	TotalClaimedRewards      uint64
}

type PositionState struct {
	Pool                     solana.PublicKey
	NftMint                  solana.PublicKey
	FeeAPerTokenCheckpoint   [32]uint8
	FeeBPerTokenCheckpoint   [32]uint8
	FeeAPending              uint64
	FeeBPending              uint64
	UnlockedLiquidity        uint128.Uint128
	VestedLiquidity          uint128.Uint128
	PermanentLockedLiquidity uint128.Uint128
	Metrics                  PositionMetrics
	RewardInfos              [2]UserRewardInfo
	Padding                  [6]uint128.Uint128
}

type PositionResult struct {
	PositionNftAccount solana.PublicKey
	Position           solana.PublicKey
	PositionState      PositionState
}

const (
	LIQUIDITY_SCALE = 128
)

type UnclaimReward struct {
	FeeTokenA uint128.Uint128
	FeeTokenB uint128.Uint128
	Rewards   []uint128.Uint128
}
