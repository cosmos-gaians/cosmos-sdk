package group

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryGet                = "get"
	QueryGroups             = "groups"
	QueryGroupsByMember     = "groups_by_member"
	QueryProposalsByGroupID = "proposals_by_group_id"
)

type QueryGroupsByMemberParams struct {
	Address sdk.AccAddress
}

type QueryProposalsByGroupIDrParams struct {
	Address sdk.AccAddress
}

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGroup(ctx, path[1:], req, keeper)
		case QueryGroups:
			return queryGroups(ctx, path[1:], req, keeper)
		case QueryGroupsByMember:
			return queryGroupsByMemberAddress(ctx, path[1:], req, keeper)
		case QueryProposalsByGroupID:
			return queryProposalsByGroupID(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown data query endpoint")
		}
	}
}

func queryGroup(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	idStr := path[0]

	decodedId, e := sdk.AccAddressFromBech32(idStr)

	if e != nil {
		return []byte{}, sdk.ErrUnknownRequest("could not decode group ID")
	}

	info, err := keeper.GetGroupInfo(ctx, decodedId)

	if err != nil {
		return []byte{}, err
	}

	res, jsonErr := codec.MarshalJSONIndent(keeper.cdc, info)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return res, nil
}

func queryGroups(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	groups := keeper.GetGroups(ctx)

	res, jsonErr := codec.MarshalJSONIndent(keeper.cdc, groups)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return res, nil
}

func queryGroupsByMemberAddress(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	var params QueryGroupsByMemberParams
	parseErr := moduleCodec.UnmarshalJSON(req.Data, &params)
	if parseErr != nil {
		err = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", parseErr.Error()))
		return
	}

	groups := keeper.GetGroupsByMemberAddress(ctx, params.Address)

	res, jsonErr := codec.MarshalJSONIndent(keeper.cdc, groups)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return res, nil
}

func queryProposalsByGroupID(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {

	var params QueryProposalsByGroupIDrParams
	parseErr := moduleCodec.UnmarshalJSON(req.Data, &params)
	if parseErr != nil {
		err = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", parseErr.Error()))
		return
	}

	proposals := keeper.GetProposalsByGroupID(ctx, params.Address)

	res, jsonErr := codec.MarshalJSONIndent(keeper.cdc, proposals)
	if jsonErr != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", jsonErr.Error()))
	}
	return res, nil
}
