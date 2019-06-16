package delegation

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgExecDelegatedAction struct {
	Signer sdk.AccAddress `json:"signer"`
	Msgs   []sdk.Msg      `json:"msg"`
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
	Capability Capability     `json:"capability"`
	Expiration time.Time      `json:"expiration"`
}

func NewMsgDelegate(granter sdk.AccAddress, grantee sdk.AccAddress, capability Capability, expiration time.Time) MsgDelegate {
	return MsgDelegate{Granter: granter, Grantee: grantee, Capability: capability, Expiration: expiration}
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
	MsgType sdk.Msg        `json:"msg_type"`
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

type MsgDelegateFeeAllowance struct {
	Granter   sdk.AccAddress `json:"granter"`
	Grantee   sdk.AccAddress `json:"grantee"`
	Allowance FeeAllowance   `json:"allowance"`
}

func (msg MsgDelegateFeeAllowance) Route() string {
	return "delegation"
}

func (msg MsgDelegateFeeAllowance) Type() string {
	return "delegate-fee-allowance"
}

func (msg MsgDelegateFeeAllowance) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgDelegateFeeAllowance) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDelegateFeeAllowance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Granter}
}

type MsgRevokeFeeAllowance struct {
	Granter sdk.AccAddress `json:"granter"`
	Grantee sdk.AccAddress `json:"grantee"`
}

func (msg MsgRevokeFeeAllowance) Route() string {
	return "delegation"
}

func (msg MsgRevokeFeeAllowance) Type() string {
	return "revoke-fee-allowance"
}

func (msg MsgRevokeFeeAllowance) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgRevokeFeeAllowance) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgRevokeFeeAllowance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Granter}
}
