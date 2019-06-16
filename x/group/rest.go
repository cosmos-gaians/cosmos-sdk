package group

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/group/groups_by_member/{memberAddr}",
		memberGroupsHandlerFn(cliCtx),
	).Methods("GET")
}

func memberGroupsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		memberAddr := vars["memberAddr"]
		route := fmt.Sprintf("custom/%s/%s", "group", "groups_by_member")

		decodedAddr, _ := sdk.AccAddressFromBech32(memberAddr)
		params := QueryGroupsByMemberParams{
			Address: decodedAddr,
		}

		bz, _ := cliCtx.Codec.MarshalJSON(params)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
