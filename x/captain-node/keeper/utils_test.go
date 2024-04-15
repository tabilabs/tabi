package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCalculatePowerOnPeriod(t *testing.T) {
	testCases := []struct {
		proportion           uint64
		maximumPowerOnPeriod uint64
		expected             sdk.Dec
		shouldPanic          bool
	}{
		{
			proportion:           0,
			maximumPowerOnPeriod: 24,
			expected:             sdk.ZeroDec(),
		},
		{
			proportion:           24,
			maximumPowerOnPeriod: 24,
			expected:             sdk.OneDec(),
		},
		{
			proportion:           25,
			maximumPowerOnPeriod: 24,
			expected:             sdk.NewDec(1),
		},
		{
			proportion:           1,
			maximumPowerOnPeriod: 24,
			expected:             sdk.NewDecWithPrec(41666666666666667, 18),
		},
	}

	for _, tc := range testCases {
		if tc.shouldPanic {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			calculatePowerOnPeriod(tc.proportion, tc.maximumPowerOnPeriod)
		} else {
			result := calculatePowerOnPeriod(tc.proportion, tc.maximumPowerOnPeriod)
			if !result.Equal(tc.expected) {
				t.Errorf("expected: %s, got: %s", tc.expected, result)
			}
		}
	}
}
