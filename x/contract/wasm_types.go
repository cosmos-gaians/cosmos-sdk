package contract

import (
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type SendResponse struct {
	Error string `json:"error"`
	// Msgs  []sdk.Msg `json:"msgs"`
	Msgs []json.RawMessage `json:"msgs"`
}

func MockCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	return cdc
}

func ParseResponse(cdc *codec.Codec, raw string) (*SendResponse, error) {
	var out SendResponse
	err := cdc.UnmarshalJSON([]byte(raw), &out)
	if err != nil {
		return nil, err
	}
	if out.Error != "" {
		return nil, errors.New(out.Error)
	}
	return &out, nil
}
