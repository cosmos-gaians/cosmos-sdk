package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/group"
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
		params := group.QueryGroupsByMemberParams{
			Address: decodedAddr,
		}

		// bz := cliCtx.Codec.MarshalJSON()
		res, err := cliCtx.QueryWithData(route, )
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
