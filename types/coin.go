package types

import (
	"math/big"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AttoTabi defines the default coin denomination used in Tabi in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in Tabi.
	AttoTabi string = "atabi"

	// BaseDenomUnit defines the base denomination unit for Tabi.
	// 1 tabi = 1x10^{BaseDenomUnit} atabi
	BaseDenomUnit = 18

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20

	AttoVeTabi string = "avetabi"
)

// PowerReduction defines the default power reduction value for staking
var PowerReduction = sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(BaseDenomUnit), nil))

// NewTabiCoin is a utility function that returns an "atabi" coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewTabiCoin(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(AttoTabi, amount)
}

// NewTabiDecCoin is a utility function that returns an "atabi" decimal coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewTabiDecCoin(amount sdkmath.Int) sdk.DecCoin {
	return sdk.NewDecCoin(AttoTabi, amount)
}

// NewTabiCoinInt64 is a utility function that returns an "atabi" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewTabiCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(AttoTabi, amount)
}

// NewVeTabiCoinInt64 is a utility function that returns an "avetabi" coin with the given sdkmath.Int amount.
// The function will panic if the provided amount is negative.
func NewVeTabiCoinInt64(amount int64) sdk.Coin { return sdk.NewInt64Coin(AttoVeTabi, amount) }
