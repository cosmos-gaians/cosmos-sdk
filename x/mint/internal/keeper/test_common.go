package keeper

import (
	"os"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tendermint/libs/db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mint/internal/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type testInput struct {
	ctx        sdk.Context
	cdc        *codec.Codec
	mintKeeper Keeper
}

func newTestInput(t *testing.T) testInput {
	db := dbm.NewMemDB()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyMint := sdk.NewKVStoreKey(types.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))

	paramsKeeper := params.NewKeeper(types.ModuleCdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(types.ModuleCdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)
	supplyKeeper := supply.NewKeeper(types.ModuleCdc, keySupply, accountKeeper, bankKeeper, supply.DefaultCodespace)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	stakingKeeper := staking.NewKeeper(
		types.ModuleCdc, keyStaking, tkeyStaking, bankKeeper, supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace,
	)
	mintKeeper := NewKeeper(types.ModuleCdc, keyMint, paramsKeeper.Subspace(types.DefaultParamspace), &stakingKeeper, supplyKeeper)

	// set module accounts
	feeCollectorAcc := accountKeeper.NewAccountWithAddress(ctx, auth.FeeCollectorAddr)
	moduleAcc := supply.NewModuleMinterAccount(types.ModuleName)
	notBondedPool := supply.NewModuleHolderAccount(staking.NotBondedTokensName)
	bondPool := supply.NewModuleHolderAccount(staking.BondedTokensName)

	accountKeeper.SetAccount(ctx, feeCollectorAcc)
	mintKeeper.SetMinterAccount(ctx, moduleAcc)
	stakingKeeper.SetNotBondedPool(ctx, notBondedPool)
	stakingKeeper.SetBondedPool(ctx, bondPool)

	mintKeeper.SetParams(ctx, types.DefaultParams())
	mintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return testInput{ctx, types.ModuleCdc, mintKeeper}
}
