package keeper_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tabilabs/tabi/app"
	"github.com/tabilabs/tabi/crypto/ethsecp256k1"
	"github.com/tabilabs/tabi/testutil"
	utiltx "github.com/tabilabs/tabi/testutil/tx"
	tabitypes "github.com/tabilabs/tabi/types"
	captainskeeper "github.com/tabilabs/tabi/x/captains/keeper"
	"github.com/tabilabs/tabi/x/captains/types"
	claimskeeper "github.com/tabilabs/tabi/x/claims/keeper"
	claimstypes "github.com/tabilabs/tabi/x/claims/types"
	feemarkettypes "github.com/tabilabs/tabi/x/feemarket/types"
)

var (
	accounts = []sdk.AccAddress{
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()), // default member
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
	}
)

type ValOperator struct {
	accAddress  common.Address
	valAddress  sdk.ValAddress
	consAddress sdk.ConsAddress

	accPriKey *ethsecp256k1.PrivKey
	valPriKey *ethsecp256k1.PrivKey
}

type IntegrationTestSuite struct {
	suite.Suite

	App      *app.Tabi
	Ctx      sdk.Context
	AppCodec codec.Codec

	ValOps     []ValOperator
	Validators []stakingtypes.Validator

	Keeper       *captainskeeper.Keeper
	MsgServer    types.MsgServer
	ClaimsServer claimstypes.MsgServer
	QueryClient  types.QueryClient
}

var s *IntegrationTestSuite

func TestCaptainsTestSuite(t *testing.T) {
	s = new(IntegrationTestSuite)
	suite.Run(t, s)
	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest runs before each test in the suite.
func (suite *IntegrationTestSuite) SetupTest() {
	suite.execSetupTest(false, suite.T())
}

// SetupSubTest runs before each subtest in the test.
func (suite *IntegrationTestSuite) SetupSubTest() {
}

func (suite *IntegrationTestSuite) execSetupTest(checkTx bool, t require.TestingT) {
	// setup val ops
	suite.setupValOps(10)

	// setup new app
	suite.App = app.Setup(checkTx, feemarkettypes.DefaultGenesisState())
	header := testutil.NewHeader(
		1, time.Now().UTC(), "tabi_9788-1", suite.ValOps[0].consAddress, nil, nil,
	)
	suite.Ctx = suite.App.BaseApp.NewContext(checkTx, header)

	// setup keeper & msg server
	suite.Keeper = &suite.App.CaptainsKeeper
	suite.MsgServer = captainskeeper.NewMsgServerImpl(&suite.App.CaptainsKeeper)
	suite.ClaimsServer = claimskeeper.NewMsgServerImpl(&suite.App.ClaimsKeeper)
	err := testutil.FundModuleAccount(suite.Ctx, suite.App.BankKeeper, claimstypes.ModuleName,
		sdk.NewCoins(sdk.NewCoin(tabitypes.AttoVeTabi, sdk.NewInt(4_000_000_000_000_000).Mul(sdk.NewInt(1e18)))))
	require.NoError(t, err)

	// setup query client
	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, captainskeeper.NewQuerierImpl(&suite.App.CaptainsKeeper))
	suite.QueryClient = types.NewQueryClient(queryHelper)

	// setup module params & default authorized member
	params := types.DefaultParams()
	params.AuthorizedMembers = []string{accounts[0].String()}
	err = suite.App.CaptainsKeeper.SetParams(suite.Ctx, params)
	suite.Require().NoError(err)

	// setup validators
	suite.setupValidators()
}

// Commit commits and starts a new block with an updated context.
func (suite *IntegrationTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

// Commit commits a block at a given time.
func (suite *IntegrationTestSuite) CommitAfter(t time.Duration) {
	var err error
	suite.Ctx, err = testutil.Commit(suite.Ctx, suite.App, t, nil)
	suite.Require().NoError(err)
	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())

	types.RegisterQueryServer(queryHelper, captainskeeper.NewQuerierImpl(&suite.App.CaptainsKeeper))
	suite.QueryClient = types.NewQueryClient(queryHelper)
}

// genValOperator generates a new validator with a new account and consensus key.
func (suite *IntegrationTestSuite) genValOperator() (valOps ValOperator) {
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	privCons, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	valOps.accAddress = common.BytesToAddress(priv.PubKey().Address().Bytes())
	valOps.valAddress = valOps.accAddress.Bytes()
	valOps.consAddress = sdk.ConsAddress(privCons.PubKey().Address())
	valOps.valPriKey = priv
	valOps.accPriKey = privCons

	return
}

// setupValOps sets up a number of validators.
func (suite *IntegrationTestSuite) setupValOps(number int) {
	suite.ValOps = make([]ValOperator, number)
	for i := 0; i < number; i++ {
		suite.ValOps[i] = suite.genValOperator()
	}
}

// setupValidators sets up a number of validators.
func (suite *IntegrationTestSuite) setupValidators() {
	number := len(suite.ValOps)
	suite.Validators = make([]stakingtypes.Validator, number)
	for i := 0; i < number; i++ {
		validator, err := stakingtypes.NewValidator(suite.ValOps[i].valAddress, suite.ValOps[i].valPriKey.PubKey(), stakingtypes.Description{})
		suite.Require().NoError(err)
		validator = stakingkeeper.TestingUpdateValidator(suite.App.StakingKeeper, suite.Ctx, validator, true)
		err = suite.App.StakingKeeper.AfterValidatorCreated(suite.Ctx, validator.GetOperator())
		require.NoError(suite.T(), err)
		err = suite.App.StakingKeeper.SetValidatorByConsAddr(suite.Ctx, validator)
		require.NoError(suite.T(), err)
		suite.Validators[i] = validator
	}
}
