
package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/tabilabs/ledger-go/ledger"
	"github.com/tabilabs/tabi/crypto/ethsecp256k1"
	"github.com/tabilabs/tabi/crypto/hd"
)

// AppName defines the Ledger app used for signing. Tabi uses the Ethereum app
const AppName = "Ethereum"

var (
	// SupportedAlgorithms defines the list of signing algorithms used on Tabi:
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithms = keyring.SigningAlgoList{hd.EthSecp256k1}
	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Tabi for the Ledger device:
	//  - secp256k1 (in order to comply with Cosmos SDK)
	// The Ledger derivation function is responsible for all signing and address generation.
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{hd.EthSecp256k1}
	// LedgerDerivation defines the Tabi Ledger Go derivation (Ethereum app with EIP-712 signing)
	LedgerDerivation = ledger.TabiLedgerDerivation()
	// CreatePubkey uses the ethsecp256k1 pubkey with Ethereum address generation and keccak hashing
	CreatePubkey = func(key []byte) types.PubKey { return &ethsecp256k1.PubKey{Key: key} }
	// SkipDERConversion represents whether the signed Ledger output should skip conversion from DER to BER.
	// This is set to true for signing performed by the Ledger Ethereum app.
	SkipDERConversion = true
)

// EthSecp256k1Option defines a function keys options for the ethereum Secp256k1 curve.
// It supports eth_secp256k1 keys for accounts.
func Option() keyring.Option {
	return func(options *keyring.Options) {
		options.SupportedAlgos = SupportedAlgorithms
		options.SupportedAlgosLedger = SupportedAlgorithmsLedger
		//options.LedgerDerivation = func() (cosmosLedger.SECP256K1, error) { return LedgerDerivation() }
		//options.LedgerCreateKey = CreatePubkey
		//options.LedgerAppName = AppName
		//options.LedgerSigSkipDERConv = SkipDERConversion
	}
}
