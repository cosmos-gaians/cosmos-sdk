package delegation

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
		"/delegation/capabilities/{granteeAddr}/{granterAddr}/{route}/{type}",
		getCapabilitiesHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/delegation/allowfees/{granteeAddr}",
		getAllowFeesHandlerFn(cliCtx),
	).Methods("GET")
}

type QueryCapabilityParams struct {
	Grantee sdk.AccAddress
	Granter sdk.AccAddress
	Route   string
	Typ     string
}

func getCapabilitiesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		grantee := vars["granteeAddr"]
		granter := vars["granterAddr"]
		rt := vars["route"]
		typ := vars["type"]
		route := fmt.Sprintf("custom/group/%s", QueryGetCaps)

		granteeAddr, err := sdk.AccAddressFromBech32(grantee)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		granterAddr, err := sdk.AccAddressFromBech32(granter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		params := QueryCapabilityParams{
			Grantee: granteeAddr,
			Granter: granterAddr,
			Route:   rt,
			Typ:     typ,
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

func getAllowFeesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		grantee := vars["granteeAddr"]
		route := fmt.Sprintf("custom/group/%s/%s", QueryGetFeeAllowances, grantee)

		res, err := cliCtx.QueryWithData(route, []byte{})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
