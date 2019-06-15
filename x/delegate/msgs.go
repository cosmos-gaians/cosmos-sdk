package delegate

import sdk "github.com/cosmos/cosmos-sdk/types"

type MsgDelegatedAction struct {
	Actor  sdk.AccAddress
	Action sdk.Msg
}

func (msg MsgDelegatedAction) Route() string {
	panic("implement me")
}

func (msg MsgDelegatedAction) Type() string {
	panic("implement me")
}

func (msg MsgDelegatedAction) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (msg MsgDelegatedAction) GetSignBytes() []byte {
	panic("implement me")
}

func (msg MsgDelegatedAction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Actor}
}
