package token_convert

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sdksim "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tabilabs/tabi/x/token-convert/client/cli"
	"github.com/tabilabs/tabi/x/token-convert/keeper"
	"github.com/tabilabs/tabi/x/token-convert/simulation"
	"github.com/tabilabs/tabi/x/token-convert/types"
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

// Name returns the token-convert module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the token-convert module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the token-convert module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the token-convert module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, cfg client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}

	return types.ValidateGenesis(data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the token-convert module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root tx command for the token-convert module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for the token-convert module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

// AppModule implements an application module for the token-convert module.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// RegisterInvariants registers the token-convert module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO: uncomment me!
	// keeper.RegisterInvariants(ir, am.keeper)
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
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(&am.keeper))

	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerierImpl(&am.keeper))
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 {
	return consensusVersion
}

// InitGenesis performs genesis initialization for the token-convert module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) []abci.ValidatorUpdate {
	var data types.GenesisState
	cdc.MustUnmarshalJSON(bz, &data)
	am.keeper.InitGenesis(ctx, data)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the token-convert module
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	data := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(data)
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents returns all the token-convert content functions that simulate proposals
func (AppModule) ProposalContents(simState module.SimulationState) []sdksim.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized token-convert param changes for simulation
func (AppModule) RandomizedParams(r *rand.Rand) []sdksim.ParamChange {
	panic("not implemented")
}

// RegisterStoreDecoder registers a decoder for the module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the module operations with their respective weights
func (AppModule) WeightedOperations(simState module.SimulationState) []sdksim.WeightedOperation {
	panic("not implemented")
}
