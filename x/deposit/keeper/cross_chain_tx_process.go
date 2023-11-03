package keeper

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/evmos/v14/rpc/namespaces/ethereum/eth/filters"
	types2 "github.com/exocore/x/deposit/types"
	"github.com/exocore/x/restaking_assets_manage/types"
	"log"
	"math/big"
)

type DepositParams struct {
	clientChainLzId uint64
	action          types.CrossChainOpType
	assetsAddress   []byte
	stakerAddress   []byte
	opAmount        sdkmath.Int
}

func (k Keeper) getDepositParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*DepositParams, error) {
	// check if action is deposit
	var action types.CrossChainOpType
	var err error
	readStart := uint32(0)
	readEnd := uint32(types.CrossChainActionLength)
	r := bytes.NewReader(log.Data[readStart:readEnd])
	err = binary.Read(r, binary.BigEndian, &action)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read action")
	}
	if action != types.Deposit {
		// not handle the actions that isn't deposit
		return nil, nil
	}

	var clientChainLzId uint64
	r = bytes.NewReader(log.Topics[types.ClientChainLzIdIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &clientChainLzId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read clientChainLzId from topic")
	}

	clientChainInfo, err := k.retakingStateKeeper.GetClientChainInfoByIndex(ctx, clientChainLzId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when get client chain info")
	}

	//decode the action parameters
	readStart = readEnd
	readEnd += clientChainInfo.AddressLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	assetsAddress := make([]byte, clientChainInfo.AddressLength)
	err = binary.Read(r, binary.BigEndian, assetsAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += clientChainInfo.AddressLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	depositAddress := make([]byte, clientChainInfo.AddressLength)
	err = binary.Read(r, binary.BigEndian, depositAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += types.CrossChainOpAmountLength
	amount := sdkmath.NewIntFromBigInt(big.NewInt(0).SetBytes(log.Data[readStart:readEnd]))

	return &DepositParams{
		clientChainLzId: clientChainLzId,
		action:          action,
		assetsAddress:   assetsAddress,
		stakerAddress:   depositAddress,
		opAmount:        amount,
	}, nil
}

func (k Keeper) FilterCrossChainEventLogs(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) ([]*ethtypes.Log, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	//filter needed logs
	addresses := []common.Address{common.HexToAddress(params.ExoCoreLzAppAddress)}
	topics := [][]common.Hash{
		{common.HexToHash(params.ExoCoreLzAppEventTopic)},
	}
	needLogs := filters.FilterLogs(receipt.Logs, nil, nil, addresses, topics)
	return needLogs, nil
}

func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	//TODO check if contract address is valid layerZero relayer address

	needLogs, err := k.FilterCrossChainEventLogs(ctx, msg, receipt)
	if err != nil {
		return err
	}
	if len(needLogs) == 0 {
		log.Println("the hook message doesn't have any event needed to handle")
		return nil
	}

	for _, log := range needLogs {
		depositParams, err := k.getDepositParamsFromEventLog(ctx, log)
		if err != nil {
			return err
		}
		if depositParams != nil {
			err = k.Deposit(ctx, depositParams)
			if err != nil {
				// todo: need to test if the changed storage state will be reverted when there is an error occurred
				return err
			}
		}
	}
	return nil
}

func (k Keeper) Deposit(ctx sdk.Context, event *DepositParams) error {
	//check event parameter before executing deposit operation
	if event.opAmount.IsNegative() {
		return errorsmod.Wrap(types2.ErrDepositAmountIsNegative, fmt.Sprintf("the amount is:%s", event.opAmount))
	}
	stakeId, assetId := types.GetStakeIDAndAssetId(event.clientChainLzId, event.stakerAddress, event.assetsAddress)
	//check if asset exist
	if !k.retakingStateKeeper.StakingAssetIsExist(ctx, assetId) {
		return errorsmod.Wrap(types2.ErrDepositAssetNotExist, fmt.Sprintf("the assetId is:%s", assetId))
	}
	changeAmount := types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: event.opAmount,
		CanWithdrawAmountOrWantChangeValue:  event.opAmount,
	}
	err := k.retakingStateKeeper.UpdateStakerAssetState(ctx, stakeId, assetId, changeAmount)
	if err != nil {
		return err
	}
	return nil
}
