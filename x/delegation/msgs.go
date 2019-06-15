package delegation

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type MsgExecDelegatedAction struct {
	Signer sdk.AccAddress `json:"signer"`
	Msg    sdk.Msg        `json:"msg"`
}

func (msg MsgExecDelegatedAction) Route() string {
	return "delegation"
}

func (msg MsgExecDelegatedAction) Type() string {
    return "exec_delegated"
}

func (msg MsgExecDelegatedAction) ValidateBasic() sdk.Error {
    return nil
}

func (msg MsgExecDelegatedAction) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgExecDelegatedAction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgDelegate struct {
	Granter    sdk.AccAddress `json:"granter"`
	Grantee    sdk.AccAddress `json:"grantee"`
	Capability Capability `json:"capability"`
	Expiration time.Time `json:"expiration"`
}

func (msg MsgDelegate) Route() string {
    return "delegation"
}

func (msg MsgDelegate) Type() string {
    return "delegate"
}

func (msg MsgDelegate) ValidateBasic() sdk.Error {
    return nil
}

func (msg MsgDelegate) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
    return []sdk.AccAddress{msg.Granter}
}

type MsgRevoke struct {
	Granter sdk.AccAddress `json:"granter"`
	Grantee sdk.AccAddress `json:"grantee"`
	MsgType sdk.Msg `json:"msg_type"`
}

func (msg MsgRevoke) Route() string {
	return "delegation"
}

func (msg MsgRevoke) Type() string {
    return "revoke"
}

func (msg MsgRevoke) ValidateBasic() sdk.Error {
    return nil
}

func (msg MsgRevoke) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgRevoke) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Granter}
}

