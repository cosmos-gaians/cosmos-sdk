package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Get all delegations from a delegator
	// r.HandleFunc(
	// 	"/staking/delegators/{delegatorAddr}/delegations",
	// 	delegatorDelegationsHandlerFn(cliCtx),
	// ).Methods("GET")

}

// func delegatorDelegationsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
// 	return queryDelegator(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDelegatorDelegations))
// }
