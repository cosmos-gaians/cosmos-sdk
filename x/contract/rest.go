package contract

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/contracts/state/{addr}",
		contractStateHandlerFn(cliCtx),
	).Methods("GET")
}

func contractStateHandlerFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		memberAddr := vars["addr"]
		route := fmt.Sprintf("custom/%s/%s", "contract", "state")

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

	rest.PostProcessResponse(w, cliCtx, res)
}
