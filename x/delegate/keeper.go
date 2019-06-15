package delegate

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

var _ Keeper = keeper{}

type capabilityGrant struct {
	//// all the actors that delegated this capability to the actor
	//// the capability should be cleared if root is false and this array is cleared
	//delegatedBy []sdk.AccAddress
	//
	//// whenever this capability is undelegated or revoked, these delegations
	//// need to be cleared recursively
	//delegatedTo []sdk.AccAddress

	capability Capability

	expiration time.Time
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return &keeper{storeKey: storeKey, cdc: cdc}
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

func (k keeper) Delegate(ctx sdk.Context, grantee sdk.AccAddress, grantor sdk.AccAddress, capability Capability, expiration time.Time) bool {
	//store := ctx.KVStore(k.storeKey)
	//grantorGrant, found := k.getCapabilityGrant(ctx, grantor, capability)
	//if !found {
	//	return false
	//}
	//if !bytes.Equal(grantor, actor) {
	//
	//}
	//grantorGrant.delegatedTo = append(grantorGrant.delegatedTo, grantee)
	//store.Set(ActorCapabilityKey(capability, grantor), k.cdc.MustMarshalBinaryBare(grantorGrant))
	//granteeGrant, _ := k.getCapabilityGrant(ctx, grantee, capability)
	//granteeGrant.delegatedBy = append(granteeGrant.delegatedBy, grantor)
	//store.Set(ActorCapabilityKey(capability, grantee), k.cdc.MustMarshalBinaryBare(granteeGrant))
	//return true
	panic("TODO")
}


func (k keeper) update(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, updated Capability) {
	grant, found := k.getCapabilityGrant(ctx, grantee, granter, updated.MsgType())
	if !found {
		return
	}
	grant.capability = updated

}

func (k keeper) Undelegate(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) {
	panic("implement me")
}

func (k keeper) GetCapability(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) Capability {
	grant, found := k.getCapabilityGrant(ctx, grantee, granter, msgType)
	if !found {
		return nil
	}
	if grant.expiration.Before(ctx.BlockHeader().Time) {
		k.Undelegate(ctx, grantee, granter, msgType)
		return nil
	}
	return grant.capability
}
