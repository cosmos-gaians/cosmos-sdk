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
			return k.DispatchAction(ctx, msg.Signer, msg.Msg)
		case MsgRevoke:
			k.Revoke(ctx, msg.Grantee, msg.Granter, msg.MsgType)
			return sdk.Result{}
		default:
			errMsg := fmt.Sprintf("Unrecognized data Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
