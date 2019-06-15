package delegation

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	router sdk.Router
}

var _ Dispatcher = keeper{}

type capabilityGrant struct {
	capability Capability

	expiration time.Time
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, router sdk.Router) Keeper {
	return &keeper{storeKey, cdc, router}
}

func ActorCapabilityKey(grantee sdk.AccAddress, granter sdk.AccAddress, msg sdk.Msg) []byte {
	return []byte(fmt.Sprintf("c/%s/%s/%s/%s", grantee, granter, msg.Route(), msg.Type()))
}

func (k keeper) getCapabilityGrant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) (grant capabilityGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(ActorCapabilityKey(grantee, granter, msgType))
	if bz == nil {
		return grant, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &grant)
	return grant, true
}

func (k keeper) Delegate(ctx sdk.Context, grantee sdk.AccAddress, grantor sdk.AccAddress, capability Capability, expiration time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(capabilityGrant{capability, expiration})
	store.Set(ActorCapabilityKey(grantee, grantor, capability.MsgType()), bz)
}

func (k keeper) update(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, updated Capability) {
	grant, found := k.getCapabilityGrant(ctx, grantee, granter, updated.MsgType())
	if !found {
		return
	}
	grant.capability = updated
}

func (k keeper) Undelegate(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ActorCapabilityKey(grantee, granter, msgType))
}

func (k keeper) GetCapability(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) Capability {
	grant, found := k.getCapabilityGrant(ctx, grantee, granter, msgType)
	if !found {
		return nil
	}
	if !grant.expiration.IsZero() && grant.expiration.Before(ctx.BlockHeader().Time) {
		k.Undelegate(ctx, grantee, granter, msgType)
		return nil
	}
	return grant.capability
}

func (k keeper) DispatchAction(ctx sdk.Context, sender sdk.AccAddress, msg sdk.Msg) sdk.Result {
	signers := msg.GetSigners()
	if len(signers) != 1 {
		return sdk.ErrUnknownRequest("can only dispatch a delegated msg with 1 signer").Result()
	}
	actor := signers[0]
	if !bytes.Equal(actor, sender) {
		cap := k.GetCapability(ctx, sender, actor, msg)
		if cap == nil {
			return sdk.ErrUnauthorized("unauthorized").Result()
		}
		allow, updated, delete := cap.Accept(msg, ctx.BlockHeader())
		if !allow {
			return sdk.ErrUnauthorized("unauthorized").Result()
		}
		if delete {
			k.Undelegate(ctx, sender, actor, msg)
		} else if updated != nil {
			k.update(ctx, sender, actor, updated)
		}
	}
	return k.router.Route(msg.Route())(ctx, msg)
}

