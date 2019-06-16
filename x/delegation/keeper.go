package delegation

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	router   sdk.Router
}

type capabilityGrant struct {
	Capability Capability

	Expiration time.Time
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, router sdk.Router) Keeper {
	return Keeper{storeKey, cdc, router}
}

func ActorCapabilityKey(grantee sdk.AccAddress, granter sdk.AccAddress, msg sdk.Msg) []byte {
	return actorCapabilityKey(grantee, granter, msg.Route(), msg.Type())
}

func actorCapabilityKey(grantee sdk.AccAddress, granter sdk.AccAddress, route, typ string) []byte {
	return []byte(fmt.Sprintf("c/%x/%x/%s/%s", grantee, granter, route, typ))
}

func FeeAllowanceKey(grantee sdk.AccAddress, granter sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("f/%x/%x", grantee, granter))
}

// grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg
func (k Keeper) getCapabilityGrant(ctx sdk.Context, actor []byte) (grant capabilityGrant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(actor)
	if bz == nil {
		return grant, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &grant)
	return grant, true
}

func (k Keeper) Delegate(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, capability Capability, expiration time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(capabilityGrant{capability, expiration})
	actor := ActorCapabilityKey(grantee, granter, capability.MsgType())
	store.Set(actor, bz)
}

func (k Keeper) update(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, updated Capability) {
	grant, found := k.getCapabilityGrant(ctx, ActorCapabilityKey(grantee, granter, updated.MsgType()))
	if !found {
		return
	}
	grant.Capability = updated
}

func (k Keeper) Revoke(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ActorCapabilityKey(grantee, granter, msgType))
}

func (k Keeper) GetCapability(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) Capability {
	grant, found := k.getCapabilityGrant(ctx, ActorCapabilityKey(grantee, granter, msgType))
	if !found {
		return nil
	}
	if !grant.Expiration.IsZero() && grant.Expiration.Before(ctx.BlockHeader().Time) {
		k.Revoke(ctx, grantee, granter, msgType)
		return nil
	}
	return grant.Capability
}

func (k Keeper) DispatchActions(ctx sdk.Context, sender sdk.AccAddress, msgs []sdk.Msg) sdk.Result {
	var res sdk.Result
	for _, msg := range msgs {
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
		res = k.router.Route(msg.Route())(ctx, msg)
		if !res.IsOK() {
			return res
		}
	}
	return res
}

func (k Keeper) DelegateFeeAllowance(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, allowance FeeAllowance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(allowance)
	store.Set(FeeAllowanceKey(grantee, granter), bz)
}

func (k Keeper) RevokeFeeAllowance(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(FeeAllowanceKey(grantee, granter))
}

type FeeAllowanceGrant struct {
	Allowance FeeAllowance   `json:"allowance"`
	Grantee   sdk.AccAddress `json:"grantee"`
	Granter   sdk.AccAddress `json:"granter"`
}

func (k Keeper) GetFeeAllowances(ctx sdk.Context, grantee sdk.AccAddress) []FeeAllowanceGrant {
	prefix := fmt.Sprintf("g/%x/", grantee)
	prefixBytes := []byte(prefix)
	store := ctx.KVStore(k.storeKey)
	var grants []FeeAllowanceGrant
	iter := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		granter, _ := sdk.AccAddressFromHex(string(iter.Key()[len(prefix):]))
		bz := iter.Value()
		var allowance FeeAllowance
		k.cdc.MustUnmarshalBinaryBare(bz, &allowance)
		grants = append(grants, FeeAllowanceGrant{
			Allowance: allowance,
			Grantee:   grantee,
			Granter:   granter,
		})
	}
	return grants
}

func (k Keeper) AllowDelegatedFees(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, fee sdk.Coins) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(FeeAllowanceKey(grantee, granter))
	if len(bz) == 0 {
		return false
	}
	var allowance FeeAllowance
	k.cdc.MustUnmarshalBinaryBare(bz, &allowance)
	if allowance == nil {
		return false
	}
	allow, updated, delete := allowance.Accept(fee, ctx.BlockHeader())
	if allow == false {
		return false
	}
	if delete {
		k.RevokeFeeAllowance(ctx, grantee, granter)
	} else if updated != nil {
		k.DelegateFeeAllowance(ctx, grantee, granter, updated)
	}
	return true
}
