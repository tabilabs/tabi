package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	utiltx "github.com/tabilabs/tabi/testutil/tx"
	"github.com/tabilabs/tabi/x/captains/types"
)

func TestSplitEpochPrefixStoreKey(t *testing.T) {
	owner := sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	epochId := uint64(1)
	t.Logf("expect: node-id %d, owner %s", epochId, owner.String())

	key := types.OwnerPledgeOnEpochStoreKey(owner, epochId)
	actualEpochId, actualOwner := types.SplitEpochAndStrFromStoreKey(types.OwnerPledgeOnEpochKey, key)
	t.Logf("actual: node-id %d, owner %s", actualEpochId, sdk.AccAddress(actualOwner).String())

	require.Equal(t, epochId, actualEpochId)
	require.Equal(t, owner.String(), sdk.AccAddress(actualOwner).String())
}
