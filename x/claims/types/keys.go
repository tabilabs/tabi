package types

// nolint
const (
	// ModuleName defines the module name
	ModuleName = "claims"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// Query endpoints supported by the minting querier
	QueryParameters = "parameters"
)

var (
	// use for the keeper store
	ParamsKey = []byte{0x00}
)
