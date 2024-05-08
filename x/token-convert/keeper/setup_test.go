package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/tabilabs/tabi/app"
	"github.com/tabilabs/tabi/crypto/ethsecp256k1"
	"github.com/tabilabs/tabi/testutil"
	utiltx "github.com/tabilabs/tabi/testutil/tx"
	"github.com/tabilabs/tabi/utils"

	feemarkettypes "github.com/tabilabs/tabi/x/feemarket/types"
	tokenconvertkeeper "github.com/tabilabs/tabi/x/token-convert/keeper"
	tokenconverttypes "github.com/tabilabs/tabi/x/token-convert/types"
)

var (
	accounts = []sdk.AccAddress{
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
	}
)

type TokenConvertTestSuite struct {
	suite.Suite

	app         *app.Tabi
	ctx         sdk.Context
	denom       string
	address     common.Address
	consAddress sdk.ConsAddress

	bankKeeper  tokenconverttypes.BankKeeper
	keeper      tokenconvertkeeper.Keeper
	msgServer   tokenconverttypes.MsgServer
	queryClient tokenconverttypes.QueryClient

	signer    keyring.Signer
	ethSigner ethtypes.Signer
	validator stakingtypes.Validator

	appCodec codec.Codec
}

var s *TokenConvertTestSuite

func TestCaptainsTestSuite(t *testing.T) {
	s = new(TokenConvertTestSuite)
	suite.Run(t, s)
	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest runs before each test in the suite.
func (suite *TokenConvertTestSuite) SetupTest() {
	suite.execSetupTest(false, suite.T())
}

// SetupSubTest runs before each subtest in the test.
func (suite *TokenConvertTestSuite) SetupSubTest() {
}

func (suite *TokenConvertTestSuite) execSetupTest(checkTx bool, t require.TestingT) {
	// setup fee denom
	suite.denom = utils.BaseDenom

	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
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

	// setup keeper & msg server
	suite.bankKeeper = suite.app.BankKeeper
	suite.keeper = suite.app.TokenConvertKeeper
	suite.msgServer = tokenconvertkeeper.NewMsgServerImpl(&suite.app.TokenConvertKeeper)

	// setup query client
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	tokenconverttypes.RegisterQueryServer(queryHelper, tokenconvertkeeper.NewQuerierImpl(&suite.app.TokenConvertKeeper))
	suite.queryClient = tokenconverttypes.NewQueryClient(queryHelper)

	// setup staking
	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = suite.denom
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	evmParams := suite.app.EvmKeeper.GetParams(suite.ctx)
	evmParams.EvmDenom = suite.denom
	err = suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
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
	validators := s.app.StakingKeeper.GetValidators(s.ctx, 1)
	suite.validator = validators[0]

	suite.ethSigner = ethtypes.LatestSignerForChainID(s.app.EvmKeeper.ChainID())
}

// Commit commits and starts a new block with an updated context.
func (suite *TokenConvertTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

// Commit commits a block at a given time.
func (suite *TokenConvertTestSuite) CommitAfter(t time.Duration) {
	var err error
	suite.ctx, err = testutil.Commit(suite.ctx, suite.app, t, nil)
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	tokenconverttypes.RegisterQueryServer(queryHelper, tokenconvertkeeper.NewQuerierImpl(&suite.app.TokenConvertKeeper))
	suite.queryClient = tokenconverttypes.NewQueryClient(queryHelper)
}
