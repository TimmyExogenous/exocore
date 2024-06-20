package keeper

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAVSSupportedAssets(ctx sdk.Context, avsAddr string) ([]string, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("GetAVSSupportedAssets: key is %s", avsAddr))
	}

	return avsInfo.Info.AssetId, nil
}

func (k *Keeper) GetAVSSlashContract(ctx sdk.Context, avsAddr string) (string, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return "", errorsmod.Wrap(err, fmt.Sprintf("GetAVSSlashContract: key is %s", avsAddr))
	}

	return avsInfo.Info.SlashAddr, nil
}

// GetAVSMinimumSelfDelegation returns the USD value of minimum self delegation, which
// is set for operator
func (k *Keeper) GetAVSMinimumSelfDelegation(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return sdkmath.LegacyNewDec(0), errorsmod.Wrap(err, fmt.Sprintf("GetAVSMinimumSelfDelegation: key is %s", avsAddr))
	}

	return sdkmath.LegacyNewDec(avsInfo.Info.MinSelfDelegation.Int64()), nil
}

// GetEpochEndAVSs returns the AVS list where the current block marks the end of their epoch.
func (k *Keeper) GetEpochEndAVSs(ctx sdk.Context) ([]string, error) {
	var avsList []types.AVSInfo
	k.IteratAVSInfo(ctx, func(_ int64, epochEndAVSInfo types.AVSInfo) (stop bool) {
		avsList = append(avsList, epochEndAVSInfo)
		return false
	})

	// Check if avsList is not empty
	if len(avsList) == 0 {
		return nil, fmt.Errorf("avsList is empty")
	}

	// Convert avsList to []string
	avsAddrList := make([]string, len(avsList))
	for i, avsInfo := range avsList {
		avsAddrList[i] = avsInfo.AvsAddress
	}

	return avsAddrList, nil
}
