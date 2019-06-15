package contract

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for contract module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgStoreCode:
			return handleMsgStoreCode(ctx, keeper, msg)
		case MsgCreateContract:
			return handleMsgCreateContract(ctx, keeper, msg)
		case MsgSendContract:
			return handleMsgSendContract(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized contract message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

const (
	Bech32CodePrefix = "code"
)

func EncodeBech32CodeID(id CodeID) (string, error) {
	addr := make([]byte, binary.MaxVarintLen64+1)
	n := binary.PutUvarint(addr, uint64(id))
	bch, err := bech32.ConvertAndEncode(Bech32CodePrefix, addr[:n])
	if err != nil {
		return "", err
	}
	return bch, nil
}

func DecodeBech32CodeID(bch string) (CodeID, error) {
	hrp, bz, err := bech32.DecodeAndConvert(bch)
	if err != nil {
		return 0, err
	}
	if hrp != Bech32CodePrefix {
		return 0, fmt.Errorf("expected bech32 prefix %s, got %s", Bech32CodePrefix, hrp)
	}
	n, err := binary.ReadUvarint(bytes.NewBuffer(bz))
	if err != nil {
		return 0, err
	}
	return CodeID(n), nil
}

func handleMsgStoreCode(ctx sdk.Context, keeper Keeper, msg MsgStoreCode) sdk.Result {
	id, err := keeper.StoreCode(ctx, msg.WASMByteCode)
	if err != nil {
		return err.Result()
	}
	res := sdk.Result{}
	bch, e := EncodeBech32CodeID(id)
	if e != nil {
		return sdk.ErrUnknownRequest(e.Error()).Result()
	}
	res.Tags = res.Tags.AppendTag("contract.code-id", bch)
	return res
}

func handleMsgCreateContract(ctx sdk.Context, keeper Keeper, msg MsgCreateContract) sdk.Result {
	id, err := keeper.CreateContract(ctx, msg.Sender, msg.Code, msg.InitMsg, msg.InitFunds)
	if err != nil {
		return err.Result()
	}
	res := sdk.Result{}
	res.Tags = res.Tags.AppendTag("contract.address", id.String())
	return res
}

func handleMsgSendContract(ctx sdk.Context, keeper Keeper, msg MsgSendContract) sdk.Result {
	return keeper.SendContract(ctx, msg.Sender, msg.Contract, msg.Msg, msg.Payment)
}
