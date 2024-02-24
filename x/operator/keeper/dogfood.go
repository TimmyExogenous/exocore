package keeper

import (
	"fmt"
	"github.com/ExocoreNetwork/exocore/x/operator/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	restakingtypes "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
)

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetOperatorConsKeyForChainId sets the (consensus) public key for the given operator address
// and chain id. By doing this, an operator is consenting to be an operator on the given chain.
// If a key already exists, it will be overwritten and the change in voting power will flow
// through to the validator set.
func (k *Keeper) SetOperatorConsKeyForChainId(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainId string,
	// should be tm-ed25519
	consKey tmprotocrypto.PublicKey,
) error {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		return delegationtypes.ErrOperatorNotExist
	}
	// check for slashing
	if k.slashKeeper.IsOperatorFrozen(ctx, opAccAddr) {
		return delegationtypes.ErrOperatorIsFrozen
	}
	// check if the chain id is valid
	if !k.restakingStateKeeper.AppChainInfoIsExist(ctx, chainId) {
		return restakingtypes.ErrNoAppChainKey
	}
	// if opting out, do not allow key replacement
	if k.IsOperatorOptingOutFromChainId(ctx, opAccAddr, chainId) {
		return types.ErrAlreadyOptingOut
	}
	// convert to bytes
	bz, err := consKey.Marshal()
	if err != nil {
		return errorsmod.Wrap(
			err,
			"SetOperatorConsKeyForChainId: error occurred when marshal public key",
		)
	}
	// convert to address for reverse lookup
	consAddr, err := types.TMCryptoPublicKeyToConsAddr(consKey)
	if err != nil {
		return errorsmod.Wrap(
			err,
			"SetOperatorConsKeyForChainId: error occurred when convert public key to consensus address",
		)
	}
	// check if the key is already in use by another operator
	// operators may call this function with their own key
	// to unjail themselves, so we will allow that.
	keyInUse, existingAddr := k.GetOperatorAddressForChainIdAndConsAddr(ctx, chainId, consAddr)
	if keyInUse {
		if !existingAddr.Equals(opAccAddr) {
			return types.ErrConsKeyAlreadyInUse
		}
	}
	// check that such a key is already set. if yes, we will consider it as key replacement.
	found, prevKey, err := k.getOperatorConsKeyForChainId(ctx, opAccAddr, chainId)
	if err != nil {
		// this should not happen
		panic(err)
	}
	var alreadyRecorded bool
	if found {
		// ultimately performs bytes.Equal
		if prevKey.Equal(consKey) {
			// no-op
			return nil
		}
		// if this key is different, we will set the vote power of the old key to 0
		// in the validator update. but, we must only do so once in a block, since the
		// first existing key is the one to replace with 0 vote power and not any others.
		alreadyRecorded, _, err = k.getOperatorPrevConsKeyForChainId(ctx, opAccAddr, chainId)
		if err != nil {
			// this should not happen
			panic(err)
		}
		if !alreadyRecorded {
			if err := k.setOperatorPrevConsKeyForChainId(ctx, opAccAddr, chainId, prevKey); err != nil {
				// this should not happen
				panic(err)
			}
		}
	}
	// 	k.setOperatorConsKeyForChainId(ctx, opAccAddr, chainId, bz)
	// 	return nil
	// }

	// // setOperatorConsKeyForChainId is the internal private version. It performs
	// // no error checking of the input.
	// func (k Keeper) setOperatorConsKeyForChainId(
	// 	ctx sdk.Context,
	// 	opAccAddr sdk.AccAddress,
	// 	chainId string,
	// 	bz []byte,
	// ) {
	store := ctx.KVStore(k.storeKey)
	// forward lookup
	// given operator address and chain id, find the consensus key,
	// since it is sorted by operator address, it helps for faster indexing by operator
	// for example, when an operator is delegated to, we can find all impacted
	// chain ids and their respective consensus keys
	store.Set(types.KeyForOperatorAndChainIdToConsKey(opAccAddr, chainId), bz)
	// reverse lookups
	// 1. given chain id and operator address, find the consensus key,
	// at initial onboarding of an app chain, it will allow us to find all
	// operators that have opted in and their consensus keys
	store.Set(types.KeyForChainIdAndOperatorToConsKey(chainId, opAccAddr), bz)
	// 2. given a chain id and a consensus addr, find the operator address,
	// the slashing module asks for an operator to be slashed by their consensus
	// address, so this will allow us to find the operator address to slash.
	// however, we do not want to retain this information forever, so we will
	// prune it once the validator set update id matures (if key replacement).
	// this pruning will be triggered by the app chain module and will not be
	// recorded here.
	store.Set(types.KeyForChainIdAndConsKeyToOperator(chainId, consAddr), opAccAddr.Bytes())
	if found {
		if !alreadyRecorded {
			k.Hooks().AfterOperatorKeyReplacement(ctx, opAccAddr, prevKey, consKey, chainId)
		}
	} else {
		k.Hooks().AfterOperatorOptIn(ctx, opAccAddr, chainId, consKey)
	}
	return nil
}

// setOperatorPrevConsKeyForChainId sets the previous (consensus) public key for the given
// operator address and chain id. This is used to track the previous key when a key is replaced.
// It is internal-only because such a key must only be set upon key replacement. So it does
// not perform any meaningful error checking of the input.
func (k *Keeper) setOperatorPrevConsKeyForChainId(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainId string,
	prevKey tmprotocrypto.PublicKey,
) error {
	bz, err := prevKey.Marshal()
	if err != nil {
		return errorsmod.Wrap(
			err,
			"SetOperatorPrevConsKeyForChainId: error occurred when marshal public key",
		)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForOperatorAndChainIdToPrevConsKey(opAccAddr, chainId), bz)
	return nil
}

// GetOperatorPrevConsKeyForChainId gets the previous (consensus) public key for the given
// operator address and chain id. When such a key is returned, callers should set its vote power
// to 0 in the validator update.
func (k *Keeper) GetOperatorPrevConsKeyForChainId(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainId string,
) (found bool, key tmprotocrypto.PublicKey, err error) {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		err = delegationtypes.ErrOperatorNotExist
		return
	}
	if !k.restakingStateKeeper.AppChainInfoIsExist(ctx, chainId) {
		err = restakingtypes.ErrNoAppChainKey
		return
	}
	// do not check for slashing here
	found, key, err = k.getOperatorPrevConsKeyForChainId(ctx, opAccAddr, chainId)
	return
}

// getOperatorPrevConsKeyForChainId is the internal version of
// GetOperatorPrevConsKeyForChainId.
// It performs no error checking of the input.
func (k *Keeper) getOperatorPrevConsKeyForChainId(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainId string,
) (found bool, key tmprotocrypto.PublicKey, err error) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.KeyForOperatorAndChainIdToPrevConsKey(opAccAddr, chainId))
	if res == nil {
		return
	}
	if err = key.Unmarshal(res); err != nil {
		return
	}
	found = true
	return
}

// GetOperatorConsKeyForChainId gets the (consensus) public key for the given operator address
// and chain id. This should be exposed via the query surface.
func (k *Keeper) GetOperatorConsKeyForChainId(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainId string,
) (found bool, key tmprotocrypto.PublicKey, err error) {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		err = delegationtypes.ErrOperatorNotExist
		return
	}
	if !k.restakingStateKeeper.AppChainInfoIsExist(ctx, chainId) {
		err = restakingtypes.ErrNoAppChainKey
		return
	}
	// do not check for slashing, since this function will be used to update voting power even
	// when slashed
	found, key, err = k.getOperatorConsKeyForChainId(ctx, opAccAddr, chainId)
	return
}

// getOperatorConsKeyForChainId is the internal version of GetOperatorConsKeyForChainId. It
// performs no error checking of the input.
func (k *Keeper) getOperatorConsKeyForChainId(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainId string,
) (found bool, key tmprotocrypto.PublicKey, err error) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.KeyForOperatorAndChainIdToConsKey(opAccAddr, chainId))
	if res == nil {
		return
	}
	if err = key.Unmarshal(res); err != nil {
		return
	}
	return true, key, nil
}

// GetChainIdsAndKeysForOperator gets the chain ids for which the given operator address has set a
// (consensus) public key. TODO: would it be better to make this a key per operator?
// This is intentionally an array of strings because I don't see the utility for the vote power
// or the public key here. If we need it, we can add it later.
func (k *Keeper) GetChainIdsAndKeysForOperator(
	ctx sdk.Context, opAccAddr sdk.AccAddress,
) (chainIds []string, consKeys []tmprotocrypto.PublicKey) {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		return
	}
	// do not check for slashing here
	prefix := types.AppendMany(
		[]byte{types.BytePrefixForOperatorAndChainIdToConsKey},
		opAccAddr.Bytes(),
	)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(
		store, prefix,
	)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// the key returned is the full key, with the prefix. drop the prefix and the length.
		chainId := string(iterator.Key()[len(prefix)+8:])
		var key tmprotocrypto.PublicKey
		if err := key.Unmarshal(iterator.Value()); err != nil {
			// grave error because we are the ones who stored this information in the first
			// place
			panic(err)
		}
		chainIds = append(chainIds, chainId)
		consKeys = append(consKeys, key)
	}
	return
}

// GetOperatorsForChainId returns a list of {operatorAddr, pubKey} for the given
// chainId. This is used to create or update the validator set. It skips
// jailed or frozen operators.
func (k *Keeper) GetOperatorsForChainId(
	ctx sdk.Context, chainId string,
) (addrs []sdk.AccAddress, pubKeys []tmprotocrypto.PublicKey) {
	if !k.restakingStateKeeper.AppChainInfoIsExist(ctx, chainId) {
		return
	}
	// prefix is the byte prefix and then chainId with length
	prefix := types.ChainIdAndAddrKey(
		types.BytePrefixForChainIdAndOperatorToConsKey,
		chainId,
		nil,
	)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(
		store, prefix,
	)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// this key is of the format prefix | len | chainId | addr
		// and our prefix is of the format prefix | len | chainId
		// so just drop it and convert to sdk.AccAddress
		addr := iterator.Key()[len(prefix):]
		res := iterator.Value()
		var ret tmprotocrypto.PublicKey
		if err := ret.Unmarshal(res); err != nil {
			// grave error
			panic(err)
		}
		addrs = append(addrs, addr)
		pubKeys = append(pubKeys, ret)

	}
	return
}

func (k *Keeper) GetOperatorAddressForChainIdAndConsAddr(
	ctx sdk.Context, chainId string, consAddr sdk.ConsAddress,
) (found bool, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.KeyForChainIdAndConsKeyToOperator(chainId, consAddr))
	if res == nil {
		return
	}
	found = true
	addr = sdk.AccAddress(res)
	return
}

// DeleteOperatorAddressForChainIdAndConsAddr is a pruning method used to delete the
// mapping from chain id and consensus address to operator address. This mapping is used
// to obtain the operator address from its consensus public key, which is sent to the
// coordinator chain by a subscriber chain for slashing.
func (k *Keeper) DeleteOperatorAddressForChainIdAndConsAddr(
	ctx sdk.Context, chainId string, consAddr sdk.ConsAddress,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForChainIdAndConsKeyToOperator(chainId, consAddr))
}

// SetHooks stores the given hooks implementations.
// Note that the Keeper is changed into a pointer to prevent an ineffective assignment.
func (k *Keeper) SetHooks(hooks types.OperatorConsentHooks) {
	if hooks == nil {
		panic("cannot set nil hooks")
	}
	if k.hooks != nil {
		panic("cannot set hooks twice")
	}
	k.hooks = hooks
}

func (k *Keeper) Hooks() types.OperatorConsentHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set to prevent calling nil functions
		return types.MultiOperatorConsentHooks{}
	}
	return k.hooks
}

// InitiateOperatorOptOutFromChainId initiates an operator opting out from the given chain id.
// It validates whether the operator is registered, and that it is not frozen, and that the
// chain is present within the system. It also checks if the operator is already opting out.
func (k *Keeper) InitiateOperatorOptOutFromChainId(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainId string,
) error {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		return delegationtypes.ErrOperatorNotExist
	}
	// check for slashing
	if k.slashKeeper.IsOperatorFrozen(ctx, opAccAddr) {
		return delegationtypes.ErrOperatorIsFrozen
	}
	// check if the chain id is valid
	if !k.restakingStateKeeper.AppChainInfoIsExist(ctx, chainId) {
		return restakingtypes.ErrNoAppChainKey
	}
	found, key, err := k.getOperatorConsKeyForChainId(ctx, opAccAddr, chainId)
	if err != nil {
		return err
	}
	if !found {
		return types.ErrNotOptedIn
	}
	isAlreadyOptingOut := k.IsOperatorOptingOutFromChainId(ctx, opAccAddr, chainId)
	if isAlreadyOptingOut {
		return types.ErrAlreadyOptingOut
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForOperatorOptOutFromChainId(opAccAddr, chainId), []byte{})
	k.Hooks().AfterOperatorOptOutInitiated(ctx, opAccAddr, chainId, key)
	return nil
}

// IsOperatorOptingOutFromChainId returns true if the operator is opting out from the given
// chain id.
func (k *Keeper) IsOperatorOptingOutFromChainId(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainId string,
) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForOperatorOptOutFromChainId(opAccAddr, chainId))
	return bz != nil
}

// CompleteOperatorOptOutFromChainId completes the operator opting out from the given chain id.
// TODO(mm): would it be better to store as 3 states? (opted in, opting out, opted out)
func (k *Keeper) CompleteOperatorOptOutFromChainId(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainId string,
) {
	if !k.IsOperatorOptingOutFromChainId(ctx, opAccAddr, chainId) {
		panic("operator is not opting out")
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForOperatorOptOutFromChainId(opAccAddr, chainId))
}

// IsOperatorJailedForChainId add for dogfood
func (k *Keeper) IsOperatorJailedForChainId(sdk.Context, sdk.AccAddress, string) bool {
	return false
}
func (k *Keeper) Jail(sdk.Context, sdk.ConsAddress, string) {
	return
}