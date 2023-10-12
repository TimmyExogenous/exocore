package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"strings"
)

func (k Keeper) GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]*types2.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, operatorAddr.Bytes())
	defer iterator.Close()

	ret := make(map[string]*types2.OperatorSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo types2.OperatorSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		stringList := strings.SplitAfter(string(iterator.Key()), "_")
		assetId := stringList[len(stringList)-1]
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (info *types2.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorAssetInfos)
	key := operatorAddr.Bytes()
	ifExist := store.Has(key)
	if !ifExist {
		return nil, types2.ErrNoStakerAssetKey
	}

	value := store.Get(key)

	ret := types2.OperatorSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) UpdateOperatorAssetsState(ctx sdk.Context, operatorAddr sdk.Address, assetsUpdate map[string]types2.OperatorSingleAssetOrChangeInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorAssetInfos)
	for assetId, changeAmount := range assetsUpdate {
		key := types2.GetAssetStateKey(operatorAddr.String(), assetId)
		isExit := store.Has(key)
		assetState := types2.OperatorSingleAssetOrChangeInfo{
			TotalAmountOrWantChangeValue:       math.NewInt(0),
			OperatorOwnAmountOrWantChangeValue: math.NewInt(0),
		}
		if isExit {
			value := store.Get(key)
			k.cdc.MustUnmarshal(value, &assetState)
		}

		if changeAmount.TotalAmountOrWantChangeValue.IsZero() && changeAmount.OperatorOwnAmountOrWantChangeValue.IsZero() {
			return types2.ErrInputUpdateStateIsZero
		}

		if changeAmount.TotalAmountOrWantChangeValue.IsNegative() {
			if assetState.TotalAmountOrWantChangeValue.LT(changeAmount.TotalAmountOrWantChangeValue.Abs()) {
				return types2.ErrSubDepositAmountIsMoreThanOrigin
			}
		}
		if changeAmount.OperatorOwnAmountOrWantChangeValue.IsNegative() {
			if assetState.OperatorOwnAmountOrWantChangeValue.LT(changeAmount.OperatorOwnAmountOrWantChangeValue.Abs()) {
				return types2.ErrSubCanWithdrawAmountIsMoreThanOrigin
			}
		}

		if !changeAmount.TotalAmountOrWantChangeValue.IsZero() {
			assetState.TotalAmountOrWantChangeValue.Add(changeAmount.TotalAmountOrWantChangeValue)
		}

		if !changeAmount.OperatorOwnAmountOrWantChangeValue.IsZero() {
			assetState.OperatorOwnAmountOrWantChangeValue.Add(changeAmount.OperatorOwnAmountOrWantChangeValue)
		}

		bz := k.cdc.MustMarshal(&assetState)
		store.Set(key, bz)
	}
	return nil
}

func (k Keeper) GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}
