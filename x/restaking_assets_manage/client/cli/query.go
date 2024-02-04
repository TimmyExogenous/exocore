package cli

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the parent command for all incentives CLI query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the restaking_assets_manage module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueClientChainInfoByIndex(),
		QueAllClientChainInfo(),
		QueStakingAssetInfo(),
		QueAllStakingAssetsInfo(),
		QueStakerAssetInfos(),
		QueStakerSpecifiedAssetAmount(),
		QueOperatorAssetInfos(),
		QueOperatorSpecifiedAssetAmount(),
		// QueStakerExoCoreAddr(),
	)
	return cmd
}

// QueClientChainInfoByIndex queries the client chain info by index
func QueClientChainInfoByIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueClientChainInfoByIndex clientChainLzId",
		Short: "Get client chain info by layerZero Id",
		Long:  "Get client chain info by layerZero Id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			clientChainLzId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errorsmod.Wrap(types.ErrCliCmdInputArg, err.Error())
			}
			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryClientChainInfo{
				ChainIndex: clientChainLzId,
			}
			res, err := queryClient.QueClientChainInfoByIndex(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueAllClientChainInfo queries all client chain info
func QueAllClientChainInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueAllClientChainInfo",
		Short: "Get all client chain info",
		Long:  "Get all client chain info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryAllClientChainInfo{}
			res, err := queryClient.QueAllClientChainInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueStakingAssetInfo queries staking asset info
func QueStakingAssetInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueStakingAssetInfo assetAddr clientChainLzId",
		Short: "Get staking asset info",
		Long:  "Get staking asset info",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			clientChainLzId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return errorsmod.Wrap(types.ErrCliCmdInputArg, fmt.Sprintf("error arg is:%v", args[1]))
			}

			_, assetId := types.GetStakeIDAndAssetIdFromStr(clientChainLzId, "", args[0])
			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryStakingAssetInfo{
				AssetId: assetId,
			}
			res, err := queryClient.QueStakingAssetInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueAllStakingAssetsInfo queries all staking asset info
func QueAllStakingAssetsInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueAllStakingAssetsInfo",
		Short: "Get all staking asset info",
		Long:  "Get all staking asset info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryAllStakingAssetsInfo{}
			res, err := queryClient.QueAllStakingAssetsInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueStakerAssetInfos queries staker asset info
func QueStakerAssetInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueStakerAssetInfos stakerId",
		Short: "Get staker asset state",
		Long:  "Get staker asset state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryStakerAssetInfo{
				StakerId: args[0],
			}
			res, err := queryClient.QueStakerAssetInfos(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueStakerSpecifiedAssetAmount queries staker specified asset info
func QueStakerSpecifiedAssetAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueStakerSpecifiedAssetAmount clientChainId stakerAddr assetAddr",
		Short: "Get staker specified asset state",
		Long:  "Get staker specified asset state",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			clientChainLzId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errorsmod.Wrap(types.ErrCliCmdInputArg, err.Error())
			}
			stakerId, assetId := types.GetStakeIDAndAssetIdFromStr(clientChainLzId, args[1], args[2])
			req := &types.QuerySpecifiedAssetAmountReq{
				StakerId: stakerId,
				AssetId:  assetId,
			}
			res, err := queryClient.QueStakerSpecifiedAssetAmount(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueOperatorAssetInfos queries operator asset info
func QueOperatorAssetInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueOperatorAssetInfos operatorAddr",
		Short: "Get operator asset state",
		Long:  "Get operator asset state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryOperatorAssetInfos{
				OperatorAddr: args[0],
			}
			res, err := queryClient.QueOperatorAssetInfos(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueOperatorSpecifiedAssetAmount queries specified operator asset info
func QueOperatorSpecifiedAssetAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueOperatorSpecifiedAssetAmount operatorAddr clientChainId assetAddr",
		Short: "Get operator specified asset state",
		Long:  "Get operator specified asset state",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			clientChainLzId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return errorsmod.Wrap(types.ErrCliCmdInputArg, err.Error())
			}
			_, assetId := types.GetStakeIDAndAssetIdFromStr(clientChainLzId, "", args[2])
			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryOperatorSpecifiedAssetAmountReq{
				OperatorAddr: args[0],
				AssetId:      assetId,
			}
			res, err := queryClient.QueOperatorSpecifiedAssetAmount(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueStakerExoCoreAddr queries staker ExoCore address
func QueStakerExoCoreAddr() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueStakerExoCoreAddr stakerId",
		Short: "Get staker ExoCore address",
		Long:  "Get staker ExoCore address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryStakerExCoreAddr{
				StakerId: args[0],
			}
			res, err := queryClient.QueStakerExoCoreAddr(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
