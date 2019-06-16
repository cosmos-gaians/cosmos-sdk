package contract

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/delegation"
)

// Keeper is the model object for the package contract module
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	accountKeeper    auth.AccountKeeper
	bankKeeper       bank.Keeper
	delegationKeeper delegation.Keeper
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, delegationKeeper delegation.Keeper) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc, accountKeeper: accountKeeper, bankKeeper: bankKeeper, delegationKeeper: delegationKeeper}
}

var (
	keyNextCodeID     = []byte("nextCodeId")
	keyNextContractID = []byte("nextContractId")
)

type CodeID uint64

func KeyCode(id CodeID) []byte {
	return []byte(fmt.Sprintf("d/%x", id))
}

func KeyContractCode(id sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("n/%x", id))
}

func KeyContractState(id sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("s/%x", id))
}

func KeyCodeHasContract(id CodeID, contract sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("cc/%x/%x", id, contract))
}

func (k Keeper) autoIncrementID(ctx sdk.Context, nextIdKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nextIdKey)
	var id uint64 = 0
	if bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &id)
	}
	bz = k.cdc.MustMarshalBinaryBare(id + 1)
	store.Set(nextIdKey, bz)
	return id
}

func (k Keeper) getNewCodeID(ctx sdk.Context) CodeID {
	return CodeID(k.autoIncrementID(ctx, keyNextCodeID))
}

func (k Keeper) StoreCode(ctx sdk.Context, byteCode []byte) (CodeID, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	id := k.getNewCodeID(ctx)
	store.Set(KeyCode(id), byteCode)
	return id, nil
}

func (k Keeper) getNewContractId(ctx sdk.Context) sdk.AccAddress {
	id := k.autoIncrementID(ctx, keyNextContractID)
	return addrFromUint64(id)
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return addr
}

type contractMsg struct {
	ContractAddress sdk.AccAddress  `json:"contract_address"`
	Sender          sdk.AccAddress  `json:"sender"`
	Msg             json.RawMessage `json:"msg"`
	SentFunds       int64           `json:"sent_funds"`
}

func (k Keeper) CreateContract(ctx sdk.Context, creator sdk.AccAddress, codeId CodeID, initData []byte, coins sdk.Coins) (sdk.AccAddress, sdk.Result) {
	// Create a contract address
	addr := k.getNewContractId(ctx)

	// Create a contract account
	existingAcc := k.accountKeeper.GetAccount(ctx, addr)
	if existingAcc != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("account with address %s already exists", addr.String())).Result()
	}

	// Deposit initial contract funds
	k.accountKeeper.SetAccount(ctx, &auth.BaseAccount{Address: addr})
	err := k.bankKeeper.SendCoins(ctx, creator, addr, coins)
	if err != nil {
		return nil, err.Result()
	}

	// Retrieve contract code
	store := ctx.KVStore(k.storeKey)
	codeBz := store.Get(KeyCode(codeId))
	if len(codeBz) == 0 {
		return nil, sdk.ErrUnknownRequest("can't find contract code").Result()
	}

	// Store contract code ID
	store.Set(KeyContractCode(addr), k.cdc.MustMarshalBinaryBare(codeId))
	// Store secondary index to look up contracts using a specific CodeID
	store.Set(KeyCodeHasContract(codeId, addr), []byte{0})

	// TODO: we really need to handle coins, not just one int
	amt := coins[0].Amount.Int64()
	msg := contractMsg{
		ContractAddress: addr,
		Sender:          creator,
		Msg:             initData,
		SentFunds:       amt,
	}
	txtMsg, stdErr := json.Marshal(msg)
	if stdErr != nil {
		return nil, sdk.ErrUnknownRequest(stdErr.Error()).Result()
	}

	// TODO: setup proper db key to expose for Read/Write
	res, err := Run(k.cdc, store, KeyContractState(addr), codeBz, "init_wrapper", []interface{}{txtMsg})
	if err != nil {
		return nil, err.Result()
	}

	out := sdk.Result{}
	for _, msg := range res.Msgs {
		out = k.delegationKeeper.DispatchAction(ctx, addr, msg)
		if !out.IsOK() {
			return nil, out
		}
	}

	return addr, out
}

func (k Keeper) SendContract(ctx sdk.Context, sender sdk.AccAddress, contract sdk.AccAddress, msg []byte, coins sdk.Coins) sdk.Result {
	// Send coins
	err := k.bankKeeper.SendCoins(ctx, sender, contract, coins)
	if err != nil {
		return err.Result()
	}

	// Retrieve code ID
	store := ctx.KVStore(k.storeKey)
	codeIdBz := store.Get(KeyContractCode(contract))
	var codeId CodeID
	k.cdc.MustUnmarshalBinaryBare(codeIdBz, &codeId)

	// Retrieve code
	codeBz := store.Get(KeyCode(codeId))
	if len(codeBz) == 0 {
		return sdk.ErrUnknownRequest("can't find contract code").Result()
	}

	// TODO: we really need to handle coins, not just one int
	amt := coins[0].Amount.Int64()
	cmsg := contractMsg{
		ContractAddress: contract,
		Sender:          sender,
		Msg:             msg,
		SentFunds:       amt,
	}
	txtMsg, stdErr := json.Marshal(cmsg)
	if stdErr != nil {
		return sdk.ErrUnknownRequest(stdErr.Error()).Result()
	}

	res, err := Run(k.cdc, store, KeyContractState(contract), codeBz, "send_wrapper", []interface{}{txtMsg})
	if err != nil {
		return err.Result()
	}

	out := sdk.Result{}
	for _, msg := range res.Msgs {
		fmt.Printf("msg: %#v\n", msg)
		out = k.delegationKeeper.DispatchAction(ctx, contract, msg)
		if !out.IsOK() {
			return out
		}
	}
	return out
}
