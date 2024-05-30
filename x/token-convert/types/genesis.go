package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/errors"
)

const (
	StrategyInstant = "instant"
	Strategy90Days  = "90days"
	Strategy180Days = "180days"
)

// NewGenesisState create a module's genesis state.
func NewGenesisState(nextSeq uint64, strategies []Strategy, vouchers []Voucher) *GenesisState {
	return &GenesisState{
		VoucherSequence: nextSeq,
		Strategies:      strategies,
		Vouchers:        vouchers,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(1, defaultStrategies(), []Voucher{})
}

// defaultStrategies returns a list of default strategies
func defaultStrategies() []Strategy {
	return []Strategy{
		{
			Name:           StrategyInstant,
			Period:         0,
			ConversionRate: sdk.NewDecWithPrec(25, 2),
		},
		{
			Name:           Strategy90Days,
			Period:         90 * 24 * 60 * 60,
			ConversionRate: sdk.NewDecWithPrec(5, 1),
		},
		{
			Name:           Strategy180Days,
			Period:         180 * 24 * 60 * 60,
			ConversionRate: sdk.NewDec(1),
		},
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
func ValidateGenesis(data GenesisState) error {
	mStrategy := make(map[string]bool, len(data.Strategies))
	for _, strategy := range data.Strategies {
		if err := validateStrategy(strategy); err != nil {
			return err
		}

		_, ok := mStrategy[strategy.Name]
		if ok {
			return errors.Wrapf(ErrInvalidStrategy, "strategy name duplicated")
		}
		mStrategy[strategy.Name] = true
	}

	mVoucher := make(map[string]bool, len(data.Vouchers))
	for _, voucher := range data.Vouchers {
		if err := validateVoucher(voucher); err != nil {
			return err
		}

		// no duplicated id
		_, ok := mVoucher[voucher.Id]
		if ok {
			return errors.Wrapf(ErrInvalidVoucher, "voucher id duplicated")
		}

		// no unknown strategy
		_, ok = mStrategy[voucher.Strategy]
		if !ok {
			return errors.Wrapf(ErrInvalidStrategy, "unknown strategy in voucher %s", voucher.Id)
		}

		mVoucher[voucher.Id] = true
	}

	// NOTE: as the module deletes voucher everytime it is cancelled or withdrew, the next
	// sequence would be greater than actual vouchers leaved in the app state. So here we
	// can only make sure the seq is bigger than the items count.
	if data.VoucherSequence <= uint64(len(data.Vouchers)) {
		return errors.Wrapf(ErrInvalidStrategy, "the next usable sequence is too small")
	}

	return nil
}

func validateStrategy(strategy Strategy) error {
	if len(strategy.Name) == 0 {
		return errors.Wrapf(ErrInvalidStrategy, "strategy name is empty")
	}

	if strategy.Period < 0 {
		return errors.Wrapf(ErrInvalidStrategy, "strategy period is negative")
	}

	if strategy.ConversionRate.IsNegative() {
		return errors.Wrapf(ErrInvalidStrategy, "conversion rate is negative")
	}

	return nil
}

func validateVoucher(voucher Voucher) error {
	if len(voucher.Id) == 0 {
		return errors.Wrapf(ErrInvalidVoucher, "voucher id is empty")
	}

	if len(voucher.Strategy) == 0 {
		return errors.Wrapf(ErrInvalidVoucher, "voucher strategy is empty")
	}

	if voucher.CreatedTime <= 0 {
		return errors.Wrapf(ErrInvalidVoucher, "created time is non-positive")
	}

	if _, err := sdk.AccAddressFromBech32(voucher.Owner); err != nil {
		return err
	}

	if !voucher.Amount.IsValid() {
		return errors.Wrapf(ErrInvalidCoin, "invalid coin in voucher")
	}

	return nil
}
