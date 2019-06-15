package delegation

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDelegate:
			k.Delegate(ctx, msg.Grantee, msg.Granter, msg.Capability, msg.Expiration)
			return sdk.Result{}
		case MsgExecDelegatedAction:
			return k.DispatchActions(ctx, msg.Signer, msg.Msgs)
		case MsgRevoke:
			k.Revoke(ctx, msg.Grantee, msg.Granter, msg.MsgType)
			return sdk.Result{}
		case MsgDelegateFeeAllowance:
			k.DelegateFeeAllowance(ctx, msg.Grantee, msg.Granter, msg.Allowance)
			return sdk.Result{}
		case MsgRevokeFeeAllowance:
			k.RevokeFeeAllowance(ctx, msg.Grantee, msg.Granter)
			return sdk.Result{}
		default:
			errMsg := fmt.Sprintf("Unrecognized data Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
