package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagOwner = "owner"
)

var (
	FlagSetVouchers = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FlagSetVouchers.String(FlagOwner, "", "The owner of the vouchers")
}
