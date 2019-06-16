package delegation

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdGetFeeAllowances(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "fee-allowances [address]",
		Short: "get fee allowances granted to this address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, id), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
