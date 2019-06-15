package group

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	groupTxCmd := &cobra.Command{
		Use:   ModuleName,
		Short: "Group transactions subcommands",
	}

	// agentTxCmd.AddCommand(client.PostCommands(
	// 	agentcmd.GetCmdCreateGroup(mc.cdc),
	// 	agentcmd.GetCmdApprove(mc.cdc),
	// 	agentcmd.GetCmdUnapprove(mc.cdc),
	// 	agentcmd.GetCmdTryExec(mc.cdc),
	// 	agentcmd.GetCmdWithdraw(mc.cdc),
	// )...)

	groupTxCmd.AddCommand(
		GetCmdApprove(cdc),
		GetCmdCreateGroup(cdc),
		// GetCmdPropose(cdc),
		GetCmdTryExec(cdc),
		GetCmdUnapprove(cdc),
		// GetCmdUnjail(cdc),
		GetCmdWithdraw(cdc),
	)

	return groupTxCmd
}

func membersFromArray(arr []string) []Member {
	n := len(arr)
	res := make([]Member, n)
	for i := 0; i < n; i++ {
		strs := strings.Split(arr[i], "=")
		if len(strs) <= 0 {
			panic("empty array")
		}
		acc, err := sdk.AccAddressFromBech32(strs[0])
		if err != nil {
			panic(err)
		}
		mem := Member{
			Address: acc,
		}
		if len(strs) == 2 {
			var ok bool
			mem.Weight, ok = sdk.NewIntFromString(strs[1])
			if !ok {
				panic(fmt.Errorf("invalid weight: %s", strs[i]))
			}
		} else {
			mem.Weight = sdk.NewInt(1)
		}
		res[i] = mem
	}
	return res
}

func GetCmdCreateGroup(cdc *codec.Codec) *cobra.Command {
	var threshold int64
	var members []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create an group",
		//Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			info := Group{
				Members:           membersFromArray(members),
				DecisionThreshold: sdk.NewInt(threshold),
			}

			msg := NewMsgCreateGroup(info, account)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int64Var(&threshold, "decision-threshold", 1, "Decision threshold")
	cmd.Flags().StringArrayVar(&members, "members", []string{}, "Members")

	return cmd
}

type ActionCreator func(cmd *cobra.Command, args []string) (sdk.Msg, error)

func GetCmdPropose(cdc *codec.Codec, actionCreator ActionCreator) *cobra.Command {
	var exec bool

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			action, err := actionCreator(cmd, args)

			if err != nil {
				return err
			}

			msg := MsgCreateProposal{
				Proposer: account,
				Action:   action,
				Exec:     exec,
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().BoolVar(&exec, "exec", false, "try to execute the proposal immediately")
	return cmd
}

func getRunVote(cdc *codec.Codec, approve bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

		txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

		if err := cliCtx.EnsureAccountExists(); err != nil {
			return err
		}

		account := cliCtx.GetFromAddress()

		id := MustDecodeProposalIDBech32(args[0])

		msg := MsgVote{
			ProposalID: id,
			Voter:      account,
			Vote:       approve,
		}
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}

		cliCtx.PrintResponse = true

		return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
	}
}

func GetCmdApprove(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "approve [ID]",
		Short: "vote to approve a proposal",
		Args:  cobra.ExactArgs(1),
		RunE:  getRunVote(cdc, true),
	}
}

func GetCmdUnapprove(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unapprove [ID]",
		Short: "vote to un-approve a proposal that you have previously approved",
		Args:  cobra.ExactArgs(1),
		RunE:  getRunVote(cdc, false),
	}
}

func GetCmdTryExec(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "try-exec [ID]",
		Short: "try to execute the proposal (will fail if not enough signers have approved it)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			id := MustDecodeProposalIDBech32(args[0])

			msg := MsgTryExecuteProposal{
				ProposalID: id,
				Signer:     account,
			}
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [ID]",
		Short: "withdraw a proposer that you previously proposed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			id := MustDecodeProposalIDBech32(args[0])

			msg := MsgWithdrawProposal{
				ProposalID: id,
				Proposer:   account,
			}
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
