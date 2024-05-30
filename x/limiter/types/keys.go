package types

const (
	// ModuleName defines the module name
	ModuleName = "limiter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

const (
	prefixParams = 1
)

var ParamsKey = []byte{prefixParams}
