package keeper_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"
	"github.com/tabilabs/tabi/app"
	"github.com/tabilabs/tabi/crypto/ethsecp256k1"
	"github.com/tabilabs/tabi/testutil"
	utiltx "github.com/tabilabs/tabi/testutil/tx"
	feemarkettypes "github.com/tabilabs/tabi/x/feemarket/types"
	limiterkeeper "github.com/tabilabs/tabi/x/limiter/keeper"
	"github.com/tabilabs/tabi/x/limiter/types"
)

var (
	accounts = []sdk.AccAddress{
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()), // default member
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
	}

	s *IntegrationTestSuite
)

type IntegrationTestSuite struct {
	suite.Suite

	App      *app.Tabi
	Ctx      sdk.Context
	AppCodec codec.Codec

	MsgServer   types.MsgServer
	QueryClient types.QueryClient
}

func TestIntegrationTestSuite(t *testing.T) {
	s = new(IntegrationTestSuite)
	suite.Run(t, s)
	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.execSetup(false, suite.T())
}

func (suite *IntegrationTestSuite) execSetup(checkTx bool, t require.TestingT) {
	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	suite.App = app.Setup(checkTx, feemarkettypes.DefaultGenesisState())
	header := testutil.NewHeader(
		1, time.Now().UTC(), "tabi_9788-1", consAddress, nil, nil,
	)
	suite.Ctx = suite.App.NewContext(checkTx, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, limiterkeeper.NewQuerierImpl(&suite.App.LimiterKeeper))
	suite.QueryClient = types.NewQueryClient(queryHelper)
}
