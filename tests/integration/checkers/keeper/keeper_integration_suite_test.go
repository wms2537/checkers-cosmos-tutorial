package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	checkersapp "github.com/wms2537/checkers/app"
	"github.com/wms2537/checkers/x/checkers/keeper"
	"github.com/wms2537/checkers/x/checkers/testutil"
	"github.com/wms2537/checkers/x/checkers/types"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
	carol = testutil.Carol
)
const (
	balAlice = 50000000
	balBob   = 20000000
	balCarol = 10000000
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *checkersapp.App
	msgServer   types.MsgServer
	ctx         sdk.Context
	queryClient types.QueryClient
}

var (
	checkersModuleAddress string
)

func TestCheckersKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {
	t := &testing.T{}
	app := checkersapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	checkersModuleAddress = app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.CheckersKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.CheckersKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func makeBalance(address string, balance int64) banktypes.Balance {
	return banktypes.Balance{
		Address: address,
		Coins: sdk.Coins{
			sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(balance),
			},
		},
	}
}

func getBankGenesis() *banktypes.GenesisState {
	coins := []banktypes.Balance{
		makeBalance(alice, balAlice),
		makeBalance(bob, balBob),
		makeBalance(carol, balCarol),
	}

	state := banktypes.NewGenesisState(
		banktypes.DefaultParams(),
		coins,
		coins[0].Coins.Add(coins[1].Coins...).Add(coins[2].Coins...),
		[]banktypes.Metadata{}, []banktypes.SendEnabled{})

	return state
}

func (suite *IntegrationTestSuite) setupSuiteWithBalances() {
	suite.app.BankKeeper.InitGenesis(suite.ctx, getBankGenesis())
}

func (suite *IntegrationTestSuite) RequireBankBalance(expected int, atAddress string) {
	sdkAdd, err := sdk.AccAddressFromBech32(atAddress)
	suite.Require().Nil(err, "Failed to parse address: %s", atAddress)
	suite.Require().Equal(
		int64(expected),
		suite.app.BankKeeper.GetBalance(suite.ctx, sdkAdd, sdk.DefaultBondDenom).Amount.Int64())
}
