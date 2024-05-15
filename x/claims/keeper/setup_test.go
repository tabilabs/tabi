package keeper_test

import (
	"testing"
	"time"

	tabitypes "github.com/tabilabs/tabi/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tabilabs/tabi/crypto/ethsecp256k1"
	utiltx "github.com/tabilabs/tabi/testutil/tx"
	feemarkettypes "github.com/tabilabs/tabi/x/feemarket/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/require"
	"github.com/tabilabs/tabi/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/suite"
	"github.com/tabilabs/tabi/app"
	claimskeeper "github.com/tabilabs/tabi/x/claims/keeper"
	"github.com/tabilabs/tabi/x/claims/types"
)

type ClaimsTestSuite struct {
	suite.Suite

	app           *app.Tabi
	ctx           sdk.Context
	denom         string
	address       common.Address
	cosmosAddress sdk.AccAddress

	keeper      *claimskeeper.Keeper
	msgServer   types.MsgServer
	queryClient types.QueryClient

	signer keyring.Signer

	appCodec codec.Codec
}

var s *ClaimsTestSuite

func TestClaimsTestSuite(t *testing.T) {
	s = new(ClaimsTestSuite)
	suite.Run(t, s)
	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest runs before each test in the suite.
func (suite *ClaimsTestSuite) SetupTest() {
	suite.execSetupTest(false, suite.T())
}

// SetupSubTest runs before each subtest in the test.
func (suite *ClaimsTestSuite) SetupSubTest() {
}

func (suite *ClaimsTestSuite) execSetupTest(checkTx bool, t require.TestingT) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.cosmosAddress = priv.PubKey().Address().Bytes()
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = utiltx.NewSigner(priv)

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	// setup new app
	suite.app = app.Setup(checkTx, feemarkettypes.DefaultGenesisState())

	header := testutil.NewHeader(
		1, time.Now().UTC(), "tabi_9788-1", consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, header)

	// setup captains keeper

	// setup keeper & msg server
	suite.keeper = &suite.app.ClaimsKeeper
	suite.msgServer = claimskeeper.NewMsgServerImpl(&suite.app.ClaimsKeeper)
	err = testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, types.ModuleName, sdk.NewCoins(sdk.NewCoin(tabitypes.AttoVeTabi, sdk.NewInt(1000000000000000000))))
	require.NoError(t, err)

	// setup query client
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, claimskeeper.NewQuerierImpl(&suite.app.ClaimsKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)

	// setup module params & default authorized member
	params := types.DefaultParams()
	err = suite.app.ClaimsKeeper.SetParams(suite.ctx, params)
	require.NoError(t, err)

	// setup validators
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	err = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	require.NoError(t, err)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)
}

// Commit commits and starts a new block with an updated context.
func (suite *ClaimsTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

// Commit commits a block at a given time.
func (suite *ClaimsTestSuite) CommitAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())

	types.RegisterQueryServer(queryHelper, claimskeeper.NewQuerierImpl(&suite.app.ClaimsKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)
}
