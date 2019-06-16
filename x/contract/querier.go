package contract

import (
	//"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryGetState            = "state"
)

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGetState:
			return queryContractState(ctx, path[1], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown data query endpoint")
		}
	}
}

func queryContractState(ctx sdk.Context, bech string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	addr, e := sdk.AccAddressFromBech32(bech)
	if e != nil {
		return nil, sdk.ErrUnknownRequest(e.Error())
	}
	res = keeper.GetContractState(ctx, addr)
	return res, nil
}
