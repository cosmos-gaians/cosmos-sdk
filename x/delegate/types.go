package delegate

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/types"
)

type Action interface {
	sdk.Msg
	Actor() sdk.AccAddress
	RequiredCapabilities() []Capability
}

type Capability interface {
	// Every type of capability should be have a system wide unique key
	CapabilityKey() string
	// Accept determines whether this grant allows the provided action, and if
	// so provides an upgraded capability grant
	Accept(action Action, block abci.Header) (allow bool, updated Capability, delete bool)
}

type ActorCapability struct {
	Capability Capability
	Actor      sdk.AccAddress
}

type Keeper interface {
	// Store capabilities under the key actor-id/capability-id
	// Grant stores a root flag, and delegate
	//GrantRootCapability(ctx sdk.Context, actor sdk.AccAddress, capability Capability)
	//RevokeRootCapability(ctx sdk.Context, actor sdk.AccAddress, capability Capability)
	Delegate(ctx sdk.Context, grantor sdk.AccAddress, grantee sdk.AccAddress, capability ActorCapability) bool
	Undelegate(ctx sdk.Context, grantor sdk.AccAddress, grantee sdk.AccAddress, capability ActorCapability)
	HasCapability(ctx sdk.Context, actor sdk.AccAddress, capability ActorCapability) bool
}

type Dispatcher interface {
	DispatchAction(ctx sdk.Context, actor sdk.AccAddress, action Action) sdk.Result
}
