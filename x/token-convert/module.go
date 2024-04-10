package token_convert

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	consensusVersion uint64 = 1
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the token-convert module.
type AppModuleBasic struct {
	cdc codec.Codec
}

func (AppModuleBasic) Name() string {
	panic("not implemented")
}

// RegisterLegacyAminoCodec registers the token-convert module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {
	panic("not implemented")
}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(codectypes.InterfaceRegistry) {
	panic("not implemented")
}

// DefaultGenesis returns default genesis state as raw bytes for the token-convert module.
func (AppModuleBasic) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	panic("not implemented")
}

// ValidateGenesis performs genesis state validation for the token-convert module.
func (AppModuleBasic) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error {
	panic("not implemented")
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the token-convert module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {
	panic("not implemented")
}

// GetTxCmd returns the root tx command for the token-convert module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	panic("not implemented")
}

// GetQueryCmd returns the root query command for the token-convert module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	panic("not implemented")
}

// AppModule implements an application module for the token-convert module.
type AppModule struct {
	AppModuleBasic

	// TODO: add keepers
}

// RegisterInvariants registers the token-convert module invariants.
func (AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("not implemented")
}

// Route returns the message routing key for the token-convert module.
// Deprecated: use RegisterServices
func (AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute returns the module querier route name.
// Deprecated: use RegisterServices
func (AppModule) QuerierRoute() string {
	return ""
}

// LegacyQuerierHandler returns the token-convert module sdk.Querier.
// Deprecated: use RegisterServices
func (AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers module services
func (AppModule) RegisterServices(cfg module.Configurator) {
	panic("not implemented")
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 {
	panic("not implemented")
}

// InitGenesis performs genesis initialization for the token-convert module
func (AppModule) InitGenesis(sdk.Context, codec.JSONCodec, json.RawMessage) []abci.ValidatorUpdate {
	panic("not implemented")
}

// ExportGenesis returns the exported genesis state as raw bytes for the token-convert module
func (AppModule) ExportGenesis(sdk.Context, codec.JSONCodec) json.RawMessage {
	panic("not implemented")
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	panic("not implemented")
}

// ProposalContents returns all the token-convert content functions that simulate proposals
func (AppModule) ProposalContents(simState module.SimulationState) []simulation.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized token-convert param changes for simulation
func (AppModule) RandomizedParams(r *rand.Rand) []simulation.ParamChange {
	panic("not implemented")
}

// RegisterStoreDecoder registers a decoder for the module's types
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	panic("not implemented")
}

// WeightedOperations returns the all the module operations with their respective weights
func (AppModule) WeightedOperations(simState module.SimulationState) []simulation.WeightedOperation {
	panic("not implemented")
}
