package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/tabilabs/tabi/v1/app"
	ibctesting "github.com/tabilabs/tabi/v1/ibc/testing"
	"github.com/tabilabs/tabi/v1/utils"
	claimstypes "github.com/tabilabs/tabi/v1/x/claims/types"
	inflationtypes "github.com/tabilabs/tabi/v1/x/inflation/types"
	"github.com/tabilabs/tabi/v1/x/recovery/types"
)

func CreatePacket(amount, denom, sender, receiver, srcPort, srcChannel, dstPort, dstChannel string, seq, timeout uint64) channeltypes.Packet {
	transfer := transfertypes.FungibleTokenPacketData{
		Amount:   amount,
		Denom:    denom,
		Receiver: sender,
		Sender:   receiver,
	}
	return channeltypes.NewPacket(
		transfer.GetBytes(),
		seq,
		srcPort,
		srcChannel,
		dstPort,
		dstChannel,
		clienttypes.ZeroHeight(), // timeout height disabled
		timeout,
	)
}

func (suite *IBCTestingSuite) SetupTest() {
	// initializes 3 test chains
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 1, 2)
	suite.TabiChain = suite.coordinator.GetChain(ibcgotesting.GetChainID(1))
	suite.IBCOsmosisChain = suite.coordinator.GetChain(ibcgotesting.GetChainID(2))
	suite.IBCCosmosChain = suite.coordinator.GetChain(ibcgotesting.GetChainID(3))
	suite.coordinator.CommitNBlocks(suite.TabiChain, 2)
	suite.coordinator.CommitNBlocks(suite.IBCOsmosisChain, 2)
	suite.coordinator.CommitNBlocks(suite.IBCCosmosChain, 2)

	// Mint coins locked on the tabi account generated with secp.
	amt, ok := sdk.NewIntFromString("1000000000000000000000")
	suite.Require().True(ok)
	coinTabi := sdk.NewCoin(utils.BaseDenom, amt)
	coins := sdk.NewCoins(coinTabi)
	err := suite.TabiChain.App.(*app.Tabi).BankKeeper.MintCoins(suite.TabiChain.GetContext(), inflationtypes.ModuleName, coins)
	suite.Require().NoError(err)

	// Fund sender address to pay fees
	err = suite.TabiChain.App.(*app.Tabi).BankKeeper.SendCoinsFromModuleToAccount(suite.TabiChain.GetContext(), inflationtypes.ModuleName, suite.TabiChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)

	coinTabi = sdk.NewCoin(utils.BaseDenom, sdk.NewInt(10000))
	coins = sdk.NewCoins(coinTabi)
	err = suite.TabiChain.App.(*app.Tabi).BankKeeper.MintCoins(suite.TabiChain.GetContext(), inflationtypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.TabiChain.App.(*app.Tabi).BankKeeper.SendCoinsFromModuleToAccount(suite.TabiChain.GetContext(), inflationtypes.ModuleName, suite.IBCOsmosisChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)

	// Mint coins on the osmosis side which we'll use to unlock our atabi
	coinOsmo := sdk.NewCoin("uosmo", sdk.NewInt(10))
	coins = sdk.NewCoins(coinOsmo)
	err = suite.IBCOsmosisChain.GetSimApp().BankKeeper.MintCoins(suite.IBCOsmosisChain.GetContext(), minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.IBCOsmosisChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.IBCOsmosisChain.GetContext(), minttypes.ModuleName, suite.IBCOsmosisChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)

	// Mint coins on the cosmos side which we'll use to unlock our atabi
	coinAtom := sdk.NewCoin("uatom", sdk.NewInt(10))
	coins = sdk.NewCoins(coinAtom)
	err = suite.IBCCosmosChain.GetSimApp().BankKeeper.MintCoins(suite.IBCCosmosChain.GetContext(), minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.IBCCosmosChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.IBCCosmosChain.GetContext(), minttypes.ModuleName, suite.IBCCosmosChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)

	// Mint coins for IBC tx fee on Osmosis and Cosmos chains
	stkCoin := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amt))

	err = suite.IBCOsmosisChain.GetSimApp().BankKeeper.MintCoins(suite.IBCOsmosisChain.GetContext(), minttypes.ModuleName, stkCoin)
	suite.Require().NoError(err)
	err = suite.IBCOsmosisChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.IBCOsmosisChain.GetContext(), minttypes.ModuleName, suite.IBCOsmosisChain.SenderAccount.GetAddress(), stkCoin)
	suite.Require().NoError(err)

	err = suite.IBCCosmosChain.GetSimApp().BankKeeper.MintCoins(suite.IBCCosmosChain.GetContext(), minttypes.ModuleName, stkCoin)
	suite.Require().NoError(err)
	err = suite.IBCCosmosChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.IBCCosmosChain.GetContext(), minttypes.ModuleName, suite.IBCCosmosChain.SenderAccount.GetAddress(), stkCoin)
	suite.Require().NoError(err)

	claimparams := claimstypes.DefaultParams()
	claimparams.AirdropStartTime = suite.TabiChain.GetContext().BlockTime()
	claimparams.EnableClaims = true
	err = suite.TabiChain.App.(*app.Tabi).ClaimsKeeper.SetParams(suite.TabiChain.GetContext(), claimparams)
	suite.Require().NoError(err)

	params := types.DefaultParams()
	params.EnableRecovery = true
	err = suite.TabiChain.App.(*app.Tabi).RecoveryKeeper.SetParams(suite.TabiChain.GetContext(), params)
	suite.Require().NoError(err)

	evmParams := suite.TabiChain.App.(*app.Tabi).EvmKeeper.GetParams(s.TabiChain.GetContext())
	evmParams.EvmDenom = utils.BaseDenom
	err = suite.TabiChain.App.(*app.Tabi).EvmKeeper.SetParams(s.TabiChain.GetContext(), evmParams)
	suite.Require().NoError(err)

	suite.pathOsmosisTabi = ibctesting.NewTransferPath(suite.IBCOsmosisChain, suite.TabiChain) // clientID, connectionID, channelID empty
	suite.pathCosmosTabi = ibctesting.NewTransferPath(suite.IBCCosmosChain, suite.TabiChain)
	suite.pathOsmosisCosmos = ibctesting.NewTransferPath(suite.IBCCosmosChain, suite.IBCOsmosisChain)
	ibctesting.SetupPath(suite.coordinator, suite.pathOsmosisTabi) // clientID, connectionID, channelID filled
	ibctesting.SetupPath(suite.coordinator, suite.pathCosmosTabi)
	ibctesting.SetupPath(suite.coordinator, suite.pathOsmosisCosmos)
	suite.Require().Equal("07-tendermint-0", suite.pathOsmosisTabi.EndpointA.ClientID)
	suite.Require().Equal("connection-0", suite.pathOsmosisTabi.EndpointA.ConnectionID)
	suite.Require().Equal("channel-0", suite.pathOsmosisTabi.EndpointA.ChannelID)
}

var timeoutHeight = clienttypes.NewHeight(1000, 1000)

func (suite *IBCTestingSuite) SendAndReceiveMessage(path *ibctesting.Path, origin *ibcgotesting.TestChain, coin string, amount int64, sender string, receiver string, seq uint64) {
	// Send coin from A to B
	transferMsg := transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.NewCoin(coin, sdk.NewInt(amount)), sender, receiver, timeoutHeight, 0, "")
	_, err := ibctesting.SendMsgs(origin, ibctesting.DefaultFeeAmt, transferMsg)
	suite.Require().NoError(err) // message committed
	// Recreate the packet that was sent
	transfer := transfertypes.NewFungibleTokenPacketData(coin, strconv.Itoa(int(amount)), sender, receiver, "")
	packet := channeltypes.NewPacket(transfer.GetBytes(), seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
	// Receive message on the counterparty side, and send ack
	err = path.RelayPacket(packet)
	suite.Require().NoError(err)
}
