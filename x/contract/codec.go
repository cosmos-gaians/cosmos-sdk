package contract

import "github.com/cosmos/cosmos-sdk/codec"

var moduleCodec = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgStoreCode{}, "contract/MsgStoreCode", nil)
	cdc.RegisterConcrete(MsgCreateContract{}, "contract/MsgCreateContract", nil)
	cdc.RegisterConcrete(MsgSendContract{}, "contract/MsgSendContract", nil)
}
