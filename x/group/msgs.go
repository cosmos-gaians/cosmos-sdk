package group

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/delegation"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Creates a group on the blockchain
// Should return a tag "group.id" with the bech32 address of the group
type MsgCreateGroup struct {
	Data   Group          `json:"data"`
	Signer sdk.AccAddress `json:"signer"`
}

type MsgUpdateGroup struct {
	GroupID sdk.AccAddress `json:"group_id"`
	Data    Group          `json:"data"`
}

type CapabilityUpdateGroup struct {
	GroupIDs []sdk.AccAddress `json:"group_ids"`
}

var _ delegation.Capability = CapabilityUpdateGroup{}

type MsgCreateProposal struct {
	Proposer sdk.AccAddress `json:"proposer"`
	Group    sdk.AccAddress `json:"group"`
	Msgs     []sdk.Msg      `json:"msgs"`
	// Whether to try to execute this propose right away upon creation
	Exec bool `json:"exec,omitempty"`
}

type MsgVote struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Voter      sdk.AccAddress `json:"voter"`
	Vote       bool           `json:"vote"`
}

type MsgTryExecuteProposal struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Signer     sdk.AccAddress `json:"signer"`
}

type MsgWithdrawProposal struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Proposer   sdk.AccAddress `json:"proposer"`
}

func NewMsgCreateGroup(group Group, signer sdk.AccAddress) MsgCreateGroup {
	return MsgCreateGroup{
		Data:   group,
		Signer: signer,
	}
}

func (msg MsgCreateGroup) Route() string { return "group" }

func (msg MsgCreateGroup) Type() string { return "group.create" }

func (info Group) ValidateBasic() sdk.Error {
	if len(info.Members) <= 0 {
		return sdk.ErrUnknownRequest("Group must reference a non-empty set of members")
	}
	if !info.DecisionThreshold.IsPositive() {
		return sdk.ErrUnknownRequest(fmt.Sprintf("DecisionThreshold must be a positive integer, got %s", info.DecisionThreshold.String()))
	}
	return nil
}

func (msg MsgCreateGroup) ValidateBasic() sdk.Error {
	return msg.Data.ValidateBasic()
}

func (msg MsgCreateGroup) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCreateGroup) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

func (msg MsgCreateProposal) Route() string { return "proposal" }

func (msg MsgCreateProposal) Type() string { return "proposal.create" }

func (msg MsgCreateProposal) ValidateBasic() sdk.Error {
	return msg.Action.ValidateBasic()
}

func (msg MsgCreateProposal) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCreateProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

func (msg MsgVote) Route() string { return "group" }

func (msg MsgVote) Type() string { return "proposal.vote" }

func (msg MsgVote) ValidateBasic() sdk.Error { return nil }

func (msg MsgVote) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

func (msg MsgTryExecuteProposal) Route() string { return "group" }

func (msg MsgTryExecuteProposal) Type() string { return "group.exec-proposal" }

func (msg MsgTryExecuteProposal) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgTryExecuteProposal) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgTryExecuteProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

func (msg MsgWithdrawProposal) Route() string { return "group" }

func (msg MsgWithdrawProposal) Type() string { return "group.withdraw-proposal" }

func (msg MsgWithdrawProposal) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgWithdrawProposal) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgWithdrawProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

func (msg MsgUpdateGroup) Route() string {
	return "group"
}

func (msg MsgUpdateGroup) Type() string {
	return "group.update"
}

func (msg MsgUpdateGroup) ValidateBasic() sdk.Error {
	return msg.Data.ValidateBasic()
}

func (msg MsgUpdateGroup) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgUpdateGroup) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.GroupID}
}

func (cap CapabilityUpdateGroup) MsgType() sdk.Msg {
	return MsgUpdateGroup{}
}

func (cap CapabilityUpdateGroup) Accept(msg sdk.Msg, block abci.Header) (allow bool, updated delegation.Capability, delete bool) {
	switch msg := msg.(type) {
	case MsgUpdateGroup:
		fmt.Println(msg.Route())
		return true, nil, false
	default:
		panic("Unexpected")
	}
}
