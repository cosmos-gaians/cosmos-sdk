package group

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// // GetQueryCmd returns the cli query commands for this module
// func (mc ModuleClient) GetQueryCmd() *cobra.Command {
// 	agentQueryCmd := &cobra.Command{
// 		Use:   "group",
// 		Short: "Querying commands for the group module",
// 	}

// 	agentQueryCmd.AddCommand(client.GetCommands(
// 		agentcmd.GetCmdGetGroup(mc.storeKey, mc.cdc),
// 		agentcmd.GetCmdGetProposal(mc.storeKey, mc.cdc),
// 	)...)

// 	return agentQueryCmd
// }

// GetCmdGroup queries information about an group
func GetCmdGetGroup(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "get group by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("could not resolve group - %s \n", id)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// GetCmdProposal queries information about an proposal
func GetCmdGetProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [id]",
		Short: "get proposal by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("could not resolve proposal - %s \n", id)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
