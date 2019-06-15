package delegate

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

type Capability interface {
	// MsgType returns the type of Msg's that this capability can accept
	MsgType() sdk.Msg
	// Accept determines whether this grant allows the provided action, and if
	// so provides an upgraded capability grant
	Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Capability, delete bool)
}

type Keeper interface {
	Delegate(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, capability Capability, expiration time.Time) bool
	Undelegate(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg)
	GetCapability(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType sdk.Msg) Capability
}

type Dispatcher interface {
	DispatchAction(ctx sdk.Context, actor sdk.AccAddress, msg sdk.Msg) sdk.Result
}
