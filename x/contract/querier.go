package contract

import (
	//"github.com/cosmos/cosmos-sdk/codec"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryGetState  = "state"
	QueryListState = "list"
)

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGetState:
			return queryContractState(ctx, path[1], req, keeper)
		case QueryListState:
			return queryContractList(ctx, req, keeper)
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

func queryContractList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var addrs []string

	var i uint64
	for true {
		addr := addrFromUint64(i)
		i++
		res = keeper.GetContractState(ctx, addr)
		if res == nil {
			break
		}
		addrs = append(addrs, addr.String())
	}

	bz, e := json.MarshalIndent(addrs, "", "  ")
	if e != nil {
		return nil, sdk.ErrInvalidAddress(e.Error())
	}

	return bz, nil
}
