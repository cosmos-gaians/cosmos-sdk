package contract

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	codec "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/delegation"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type testInput struct {
	cdc    *codec.Codec
	ctx    sdk.Context
	ak     auth.AccountKeeper
	pk     params.Keeper
	bk     bank.Keeper
	dk     delegation.Keeper
	ck     Keeper
	router sdk.Router
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	delegation.RegisterCodec(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	contCapKey := sdk.NewKVStoreKey("contKey")
	delCapKey := sdk.NewKVStoreKey("delKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(contCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(delCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	bk := bank.NewBaseKeeper(ak, pk.Subspace(banktypes.DefaultParamspace), banktypes.DefaultCodespace)
	bk.SetSendEnabled(ctx, true)

	router := baseapp.NewRouter()

	dk := delegation.NewKeeper(delCapKey, cdc, router)

	router.AddRoute("bank", bank.NewHandler(bk))

	ck := NewKeeper(contCapKey, cdc, ak, bk, dk)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{cdc: cdc, ctx: ctx, ak: ak, pk: pk, bk: bk, ck: ck, dk: dk, router: router}
}

const (
	// some valid cosmos keys....
	sender    = "cosmos157ez5zlaq0scm9aycwphhqhmg3kws4qusmekll"
	recipient = "cosmos1rjxwm0rwyuldsg00qf5lt26wxzzppjzxs2efdw"
)

func TestKeeperRegen(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	addr := sdk.AccAddress([]byte(sender))
	addr2 := sdk.AccAddress([]byte(recipient))
	input.bk.SetCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("tree", 10000)))

}
