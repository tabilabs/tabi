package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tabitypes "github.com/tabilabs/tabi/types"
	claimskeeper "github.com/tabilabs/tabi/x/claims/keeper"
	claimstypes "github.com/tabilabs/tabi/x/claims/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tabilabs/tabi/app"
	"github.com/tabilabs/tabi/crypto/ethsecp256k1"
	"github.com/tabilabs/tabi/testutil"
	utiltx "github.com/tabilabs/tabi/testutil/tx"
	captainskeeper "github.com/tabilabs/tabi/x/captains/keeper"
	"github.com/tabilabs/tabi/x/captains/types"
	evmtypes "github.com/tabilabs/tabi/x/evm/types"
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

type CaptainsTestSuite struct {
	suite.Suite

	app      *app.Tabi
	ctx      sdk.Context
	appCodec codec.Codec

	valOps     []ValOperator
	validators []stakingtypes.Validator

	keeper         *captainskeeper.Keeper
	msgServer      types.MsgServer
	claimsServer   claimstypes.MsgServer
	queryClient    types.QueryClient
	queryClientEvm evmtypes.QueryClient
}

var s *CaptainsTestSuite

func TestCaptainsTestSuite(t *testing.T) {
	s = new(CaptainsTestSuite)
	suite.Run(t, s)
	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest runs before each test in the suite.
func (suite *CaptainsTestSuite) SetupTest() {
	suite.execSetupTest(false, suite.T())
}

// SetupSubTest runs before each subtest in the test.
func (suite *CaptainsTestSuite) SetupSubTest() {
}

func (suite *CaptainsTestSuite) execSetupTest(checkTx bool, t require.TestingT) {
	// setup val ops
	suite.setupValOps(10)

	// setup new app
	suite.app = app.Setup(checkTx, feemarkettypes.DefaultGenesisState())
	header := testutil.NewHeader(
		1, time.Now().UTC(), "tabi_9788-1", suite.valOps[0].consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, header)

	// setup keeper & msg server
	suite.keeper = &suite.app.CaptainsKeeper
	suite.msgServer = captainskeeper.NewMsgServerImpl(&suite.app.CaptainsKeeper)
	suite.claimsServer = claimskeeper.NewMsgServerImpl(&suite.app.ClaimsKeeper)
	err := testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, claimstypes.ModuleName,
		sdk.NewCoins(sdk.NewCoin(tabitypes.AttoVeTabi, sdk.NewInt(4_000_000_000_000_000).Mul(sdk.NewInt(1e18)))))
	require.NoError(t, err)

	// setup query client
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, captainskeeper.NewQuerierImpl(&suite.app.CaptainsKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)

	// setup module params & default authorized member
	params := types.DefaultParams()
	params.AuthorizedMembers = []string{accounts[0].String()}
	err = suite.app.CaptainsKeeper.SetParams(suite.ctx, params)
	suite.Require().NoError(err)

	// setup validators
	suite.setupValidators()
}

// Commit commits and starts a new block with an updated context.
func (suite *CaptainsTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

// Commit commits a block at a given time.
func (suite *CaptainsTestSuite) CommitAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())

	types.RegisterQueryServer(queryHelper, captainskeeper.NewQuerierImpl(&suite.app.CaptainsKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)
}

// genValOperator generates a new validator with a new account and consensus key.
func (suite *CaptainsTestSuite) genValOperator() (valOps ValOperator) {
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
func (suite *CaptainsTestSuite) setupValOps(number int) {
	suite.valOps = make([]ValOperator, number)
	for i := 0; i < number; i++ {
		suite.valOps[i] = suite.genValOperator()
	}
}

// setupValidators sets up a number of validators.
func (suite *CaptainsTestSuite) setupValidators() {
	number := len(suite.valOps)
	suite.validators = make([]stakingtypes.Validator, number)
	for i := 0; i < number; i++ {
		validator, err := stakingtypes.NewValidator(suite.valOps[i].valAddress, suite.valOps[i].valPriKey.PubKey(), stakingtypes.Description{})
		suite.Require().NoError(err)
		validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
		err = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
		require.NoError(suite.T(), err)
		err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
		require.NoError(suite.T(), err)
		suite.validators[i] = validator
	}
}
