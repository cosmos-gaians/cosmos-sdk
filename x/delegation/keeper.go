package delegation

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	router sdk.Router
}

type capabilityGrant struct {
	capability Capability

	expiration time.Time
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, router sdk.Router) Keeper {
	return Keeper{storeKey, cdc, router}
}

func ActorCapabilityKey(grantee sdk.AccAddress, granter sdk.AccAddress, msg sdk.Msg) []byte {
	return []byte(fmt.Sprintf("c/%s/%s/%s/%s", grantee, granter, msg.Route(), msg.Type()))
}

func (k Keeper) getCapabilityGrant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) (grant capabilityGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(ActorCapabilityKey(grantee, granter, msgType))
	if bz == nil {
		return grant, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &grant)
	return grant, true
}

func (k Keeper) Delegate(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, capability Capability, expiration time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(capabilityGrant{capability, expiration})
	store.Set(ActorCapabilityKey(grantee, granter, capability.MsgType()), bz)
}

func (k Keeper) update(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, updated Capability) {
	grant, found := k.getCapabilityGrant(ctx, grantee, granter, updated.MsgType())
	if !found {
		return
	}
	grant.capability = updated
}

func (k Keeper) Revoke(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ActorCapabilityKey(grantee, granter, msgType))
}

func (k Keeper) GetCapability(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) Capability {
	grant, found := k.getCapabilityGrant(ctx, grantee, granter, msgType)
	if !found {
		return nil
	}
	if !grant.expiration.IsZero() && grant.expiration.Before(ctx.BlockHeader().Time) {
		k.Revoke(ctx, grantee, granter, msgType)
		return nil
	}
	return grant.capability
}

func (k Keeper) DispatchAction(ctx sdk.Context, sender sdk.AccAddress, msg sdk.Msg) sdk.Result {
	signers := msg.GetSigners()
	if len(signers) != 1 {
		return sdk.ErrUnknownRequest("can only dispatch a delegated msg with 1 signer").Result()
	}
	actor := signers[0]
	if !bytes.Equal(actor, sender) {
		capability := k.GetCapability(ctx, sender, actor, msg)
		if capability == nil {
			return sdk.ErrUnauthorized("unauthorized").Result()
		}
		allow, updated, del := capability.Accept(msg, ctx.BlockHeader())
		if !allow {
			return sdk.ErrUnauthorized("unauthorized").Result()
		}
		if del {
			k.Revoke(ctx, sender, actor, msg)
		} else if updated != nil {
			k.update(ctx, sender, actor, updated)
		}
	}
	return k.router.Route(msg.Route())(ctx, msg)
}

