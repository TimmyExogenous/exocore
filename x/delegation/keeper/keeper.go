package keeper

import (
	"context"

	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	// isGeneralInit indicates whether the genesis state is initialized from
	// the bootStrap contract or a general exporting genesis file.
	// It's better to pass this flag through the context, but the context is
	// constructed from the base app. Seems like there isn't a good way to pass
	// the start commandline flag to the context.
	isGeneralInit bool

	// other keepers
	assetsKeeper   delegationtype.AssetsKeeper
	slashKeeper    delegationtype.SlashKeeper
	operatorKeeper delegationtype.OperatorKeeper
	hooks          delegationtype.DelegationHooks
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	isGeneralInit bool,
	assetsKeeper delegationtype.AssetsKeeper,
	slashKeeper delegationtype.SlashKeeper,
	operatorKeeper delegationtype.OperatorKeeper,
) Keeper {
	return Keeper{
		storeKey:       storeKey,
		cdc:            cdc,
		isGeneralInit:  isGeneralInit,
		assetsKeeper:   assetsKeeper,
		slashKeeper:    slashKeeper,
		operatorKeeper: operatorKeeper,
	}
}

// SetHooks stores the given hooks implementations.
// Note that the Keeper is changed into a pointer to prevent an ineffective assignment.
func (k *Keeper) SetHooks(hooks delegationtype.DelegationHooks) {
	if hooks == nil {
		panic("cannot set nil hooks")
	}
	if k.hooks != nil {
		panic("cannot set hooks twice")
	}
	k.hooks = hooks
}

func (k *Keeper) Hooks() delegationtype.DelegationHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set to prevent calling nil functions
		return delegationtype.MultiDelegationHooks{}
	}
	return k.hooks
}

// IDelegation interface will be implemented by delegation keeper
type IDelegation interface {
	// DelegateAssetToOperator handle the DelegateAssetToOperator txs from msg service
	DelegateAssetToOperator(ctx context.Context, delegation *delegationtype.MsgDelegation) (*delegationtype.DelegationResponse, error)
	// UndelegateAssetFromOperator handle the UndelegateAssetFromOperator txs from msg service
	UndelegateAssetFromOperator(ctx context.Context, delegation *delegationtype.MsgUndelegation) (*delegationtype.UndelegationResponse, error)

	GetSingleDelegationInfo(ctx sdk.Context, stakerID, assetID, operatorAddr string) (*delegationtype.DelegationAmounts, error)

	GetDelegationInfo(ctx sdk.Context, stakerID, assetID string) (*delegationtype.QueryDelegationInfoResponse, error)
}
