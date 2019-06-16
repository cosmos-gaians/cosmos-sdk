package delegation

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/cobra"
	"time"
)

func GetCmdExecDelefgated(cdc *codec.Codec) *cobra.Command {
	var exec bool

	cmd := &cobra.Command{
		Use:   "exec [msg-json]",
		Short: "execute a delegated action, with the message including its delegated signer in the encoded JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			account := cliCtx.GetFromAddress()

			var action sdk.Msg
			err := cdc.UnmarshalJSON([]byte(args[0]), &action)
			if err != nil {
				return err
			}

			msg := MsgExecDelegatedAction{
				Signer: account,
				Msgs:     []sdk.Msg{action},
			}
			err = action.ValidateBasic()
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

func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate [grantee] [capability]",
		Short: "Delegate a capability to a grantee",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var capability Capability
			err = cdc.UnmarshalJSON([]byte(args[1]), &capability)
			if err != nil {
				return err
			}

			//var expiration time.Time
			//expirationStr := cmd.Flags().GetString("expiration")

			msg := NewMsgDelegate(account, grantee, capability, time.Time{})

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String("expiration", "", "The expiration data of the delegation")
	return cmd
}

func GetCmdDelegateFees(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegate-fees [grantee] [fee-allowance]",
		Short: "delegate-fees",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			grantee, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var allowance FeeAllowance
			err = cdc.UnmarshalJSON([]byte(args[1]), &allowance)
			if err != nil {
				return err
			}

			msg := NewMsgDelegateFeeAllowance(account, grantee, allowance)

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
