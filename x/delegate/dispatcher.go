package delegate

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type dispatcher struct {
	Keeper
	Router sdk.Router
}

func NewDispatcher(k Keeper, r sdk.Router) Dispatcher {
	return &dispatcher{k, r}
}

func (dispatcher dispatcher) DispatchAction(ctx sdk.Context, sender sdk.AccAddress, msg sdk.Msg) sdk.Result {
	signers := msg.GetSigners()
	if len(signers) != 1 {
		return sdk.ErrUnknownRequest("can only dispatch a delegated msg with 1 signer").Result()
	}
	actor := signers[0]
	if !bytes.Equal(actor, sender) {
		cap := dispatcher.GetCapability(ctx, sender, actor, msg)
		if cap == nil {
			return sdk.ErrUnauthorized("unauthorized").Result()
		}
		allow, updated, delete := cap.Accept(msg, ctx.BlockHeader())
		if !allow {
			return sdk.ErrUnauthorized("unauthorized").Result()
		}
		if delete {
			dispatcher.Undelegate(ctx, sender, actor, msg)
		} else if updated != nil {
			//dispatcher.update(ctx, sender, actor, updated)
		}
	}
	return dispatcher.Router.Route(msg.Route())(ctx, msg)
}

