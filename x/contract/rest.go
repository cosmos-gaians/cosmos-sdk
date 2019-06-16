package contract

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/contracts/state/{addr}",
		contractStateHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/contracts/list",
		contractListHandlerFn(cliCtx),
	).Methods("GET")
}

func contractStateHandlerFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["addr"]
		route := fmt.Sprintf("custom/%s/%s/%s", "contract", "state", addr)
		fmt.Println(route)

		res, err := cliContext.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

func contractListHandlerFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/contract/list/abc")
		fmt.Println(route)

		res, err := cliContext.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}
