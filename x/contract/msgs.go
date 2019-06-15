package contract

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgStoreCode struct {
	Sender       sdk.AccAddress `json:"sender"`
	WASMByteCode []byte `json:"wasm_byte_code"`
}

func (msg MsgStoreCode) Route() string {
	return "contract"
}

func (msg MsgStoreCode) Type() string {
	return "store-code"
}

func (msg MsgStoreCode) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgStoreCode) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Sender}
}

type MsgCreateContract struct {
	Sender    sdk.AccAddress
	Code      CodeID
	InitMsg   []byte
	InitFunds sdk.Coins
}

func (msg MsgCreateContract) Route() string {
	return "contract"
}

func (msg MsgCreateContract) Type() string {
    return "create"
}

func (msg MsgCreateContract) ValidateBasic() sdk.Error {
    return nil
}

func (msg MsgCreateContract) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCreateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgSendContract struct {
	Sender   sdk.AccAddress
	Contract sdk.AccAddress
	Msg      [] byte
	Payment  sdk.Coins
}

func (msg MsgSendContract) Route() string {
    return "contract"
}

func (msg MsgSendContract) Type() string {
	return "send"
}

func (msg MsgSendContract) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgSendContract) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSendContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

