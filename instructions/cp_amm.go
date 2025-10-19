package instructions

import (
	"github.com/MeteoraAg/damm-v2-go/common"
	"github.com/MeteoraAg/damm-v2-go/helpers"
	"github.com/gagliardetto/solana-go"
)

func ClaimPositionFee(
	pool solana.PublicKey,
	position solana.PublicKey,
	tokenAAccount solana.PublicKey,
	tokenBAccount solana.PublicKey,
	tokenAVault solana.PublicKey,
	tokenBVault solana.PublicKey,
	tokenAMint solana.PublicKey,
	tokenBMint solana.PublicKey,
	tokenAProgram solana.PublicKey,
	tokenBProgram solana.PublicKey,
	positionNftAccount solana.PublicKey,
	owner solana.PublicKey,
) solana.Instruction {
	disc := []byte{180, 38, 154, 17, 133, 33, 162, 211}

	// Derive PDAs
	poolAuthority := helpers.DerivePoolAuthorityPDA()
	eventAuthority := helpers.DeriveEventAuthorityPDA()

	acctMeta := solana.AccountMetaSlice{
		// 1. pool_authority
		{PublicKey: poolAuthority, IsSigner: false, IsWritable: false},
		// 2. pool
		{PublicKey: pool, IsSigner: false, IsWritable: false},
		// 3. position
		{PublicKey: position, IsSigner: false, IsWritable: true},
		// 4. token_a_account
		{PublicKey: tokenAAccount, IsSigner: false, IsWritable: true},
		// 5. token_b_account
		{PublicKey: tokenBAccount, IsSigner: false, IsWritable: true},
		// 6. token_a_vault
		{PublicKey: tokenAVault, IsSigner: false, IsWritable: true},
		// 7. token_b_vault
		{PublicKey: tokenBVault, IsSigner: false, IsWritable: true},
		// 8. token_a_mint
		{PublicKey: tokenAMint, IsSigner: false, IsWritable: false},
		// 9. token_b_mint
		{PublicKey: tokenBMint, IsSigner: false, IsWritable: false},
		// 10. position_nft_account
		{PublicKey: positionNftAccount, IsSigner: false, IsWritable: false},
		// 11. owner (signer)
		{PublicKey: owner, IsSigner: true, IsWritable: false},
		// 12. token_a_program
		{PublicKey: tokenAProgram, IsSigner: false, IsWritable: false},
		// 13. token_b_program
		{PublicKey: tokenBProgram, IsSigner: false, IsWritable: false},
		// 14. event_authority (PDA)
		{PublicKey: eventAuthority, IsSigner: false, IsWritable: false},
		// 15. program
		{PublicKey: solana.MustPublicKeyFromBase58(common.DammV2ProgramID), IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		solana.MustPublicKeyFromBase58(common.DammV2ProgramID),
		acctMeta,
		disc,
	)
}
