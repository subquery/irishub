package mint_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/irisnet/irishub/v3/modules/mint"
	"github.com/irisnet/irishub/v3/modules/mint/types"
	"github.com/irisnet/irishub/v3/simapp"
)

func TestBeginBlocker(t *testing.T) {
	app, ctx := createTestApp(t, true)

	mint.BeginBlocker(ctx, app.MintKeeper)
	minter := app.MintKeeper.GetMinter(ctx)
	param := app.MintKeeper.GetParams(ctx)
	mintCoins := minter.BlockProvision(param)

	acc1 := app.AccountKeeper.GetModuleAccount(ctx, "fee_collector")
	mintedCoins := app.BankKeeper.GetAllBalances(ctx, acc1.GetAddress())
	require.Equal(t, mintedCoins, sdk.NewCoins(mintCoins))
}

// returns context and an app with updated mint keeper
func createTestApp(t *testing.T, isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(t, false)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: 2})
	app.MintKeeper.SetParams(ctx, types.NewParams(
		sdk.DefaultBondDenom,
		sdk.NewDecWithPrec(4, 2),
	))
	app.MintKeeper.SetMinter(ctx, types.DefaultMinter())
	app.DistrKeeper.SetFeePool(ctx, distributiontypes.InitialFeePool())
	return app, ctx
}
