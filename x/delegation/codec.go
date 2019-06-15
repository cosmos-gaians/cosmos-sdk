package delegation

import "github.com/cosmos/cosmos-sdk/codec"

var moduleCodec = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgExecDelegatedAction{}, "delegation/MsgExecDelegatedAction", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "delegation/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgRevoke{}, "delegation/MsgRevoke", nil)
	cdc.RegisterConcrete(capabilityGrant{}, "delegation/capabilityGrant", nil)
}