//go:build gofuzz || go1.18

package types_test

import (
	"testing"

	utiltx "github.com/tabilabs/tabi/testutil/tx"
	"github.com/tabilabs/tabi/x/incentives/types"
)

func FuzzSplitGasMeterKey(f *testing.F) {
	contract := utiltx.GenerateAddress()
	user := utiltx.GenerateAddress()

	key := types.KeyPrefixGasMeter
	key = append(key, contract.Bytes()...)
	key = append(key, user.Bytes()...)
	f.Add(key)
	f.Fuzz(func(t *testing.T, key []byte) {
		types.SplitGasMeterKey(key)
	})
}
