package withdraw

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"github.com/exocore/x/restaking_assets_manage/types"
	"github.com/exocore/x/withdraw/keeper"
	"math/big"
	"reflect"
)

func (p Precompile) GetWithdrawParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.WithdrawParams, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	withdrawParams := &keeper.WithdrawParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	withdrawParams.ClientChainLzId = uint64(clientChainLzID)

	info, err := p.stakingStateKeeper.GetClientChainInfoByIndex(ctx, withdrawParams.ClientChainLzId)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	//the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[0]), assetAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	withdrawParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[0]), stakerAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	withdrawParams.WithdrawAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 3, reflect.TypeOf(args[0]), opAmount)
	}

	withdrawParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return withdrawParams, nil
}