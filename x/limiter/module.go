package limiter

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/tabilabs/tabi/x/limiter/client/cli"
	"github.com/tabilabs/tabi/x/limiter/keeper"
	"github.com/tabilabs/tabi/x/limiter/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the captains module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the limiter module's name.
func (a AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the limiter module's types on the LegacyAmino codec.
func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types.
func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the limiter module.
func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the limiter module.
func (a AppModuleBasic) ValidateGenesis(codec codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the captains module.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the transaction commands for the limiter module
func (a AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

// GetQueryCmd returns the root query command for the limiter module.
func (a AppModuleBasic) GetQueryCmd() *cobra.Command { return cli.GetQueryCmd() }

// ____________________________________________________________________________

// AppModule implements an application module for the captains module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
	}
}

// InitGenesis performs genesis initialization for the limiter module. It returns
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)
	am.keeper.InitGenesis(ctx, &gs)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the limiter module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// RegisterInvariants registers the limiter module invariants.
func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

// Route returns the message routing key for the limiter module.
func (am AppModule) Route() sdk.Route { return sdk.Route{} }

// QuerierRoute returns the limiter module's Router.
func (am AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler returns the limiter module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier { return nil }

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(&am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerierImpl(&am.keeper))
}

// ConsensusVersion return the module consensus version.
func (am AppModule) ConsensusVersion() uint64 { return 1 }
