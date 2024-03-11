package withdraw_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/precompiles/withdraw"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	depositparams "github.com/ExocoreNetwork/exocore/x/deposit/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

func (s *WithdrawPrecompileTestSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			withdraw.MethodWithdraw,
			s.precompile.Methods[withdraw.MethodWithdraw].Name,
			true,
		},
		{
			"invalid",
			"invalid",
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(s.precompile.IsTransaction(tc.method), tc.isTx)
		})
	}
}

func paddingClientChainAddress(input []byte, outputLength int) []byte {
	if len(input) < outputLength {
		padding := make([]byte, outputLength-len(input))
		return append(input, padding...)
	}
	return input
}

func (s *WithdrawPrecompileTestSuite) TestRequiredGas() {
	clientChainLzID := 101
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), types.GeneralClientChainAddrLength)
	withdrawAddr := paddingClientChainAddress(s.Address.Bytes(), types.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress

	testcases := []struct {
		name     string
		malleate func() []byte
		expGas   uint64
	}{
		{
			"success - withdraw transaction with correct gas estimation",
			func() []byte {
				input, err := s.precompile.Pack(
					withdraw.MethodWithdraw,
					uint16(clientChainLzID),
					assetAddr,
					withdrawAddr,
					opAmount,
				)
				s.Require().NoError(err)
				return input
			},
			9680,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()

			// malleate contract input
			input := tc.malleate()
			gas := s.precompile.RequiredGas(input)

			s.Require().Equal(gas, tc.expGas)
		})
	}
}

// TestRun tests the precompiled Run method withdraw.
func (s *WithdrawPrecompileTestSuite) TestRunWithdrawThroughClientChain() {
	// deposit params for test
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzID := 101
	withdrawAmount := big.NewInt(10)
	depositAmount := big.NewInt(100)
	assetAddr := paddingClientChainAddress(usdtAddress, types.GeneralClientChainAddrLength)
	depositAsset := func(staker []byte, depositAmount sdkmath.Int) {
		// deposit asset for withdraw test
		params := &keeper.DepositParams{
			ClientChainLzID: 101,
			Action:          types.Deposit,
			StakerAddress:   staker,
			AssetsAddress:   usdtAddress,
			OpAmount:        depositAmount,
		}
		err := s.App.DepositKeeper.Deposit(s.Ctx, params)
		s.Require().NoError(err)
	}

	commonMalleate := func() (common.Address, []byte) {
		// Prepare the call input for withdraw test
		input, err := s.precompile.Pack(
			withdraw.MethodWithdraw,
			uint16(clientChainLzID),
			assetAddr,
			paddingClientChainAddress(s.Address.Bytes(), types.GeneralClientChainAddrLength),
			withdrawAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[withdraw.MethodWithdraw].Outputs.Pack(true, new(big.Int).Sub(depositAmount, withdrawAmount))
	s.Require().NoError(err)
	testcases := []struct {
		name        string
		malleate    func() (common.Address, []byte)
		readOnly    bool
		expPass     bool
		errContains string
		returnBytes []byte
	}{
		{
			name: "pass - withdraw via pre-compiles",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &depositparams.Params{
					ExoCoreLzAppAddress:    s.Address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.App.DepositKeeper.SetParams(s.Ctx, depositModuleParam)
				s.Require().NoError(err)
				depositAsset(s.Address.Bytes(), sdkmath.NewIntFromBigInt(depositAmount))
				return commonMalleate()
			},
			returnBytes: successRet,
			readOnly:    false,
			expPass:     true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			// setup basic test suite
			s.SetupTest()

			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   s.App.EvmKeeper.ChainID(),
				Nonce:     0,
				To:        &contractAddr,
				Amount:    nil,
				GasLimit:  100000,
				GasPrice:  app.MainnetMinGasPrices.BigInt(),
				GasFeeCap: baseFee,
				GasTipCap: big.NewInt(1),
				Accesses:  &ethtypes.AccessList{},
			}
			msgEthereumTx := evmtypes.NewTx(&txArgs)

			msgEthereumTx.From = s.Address.String()
			err := msgEthereumTx.Sign(s.EthSigner, s.Signer)
			s.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := s.Ctx.BlockHeader().ProposerAddress
			cfg, err := s.App.EvmKeeper.EVMConfig(s.Ctx, proposerAddress, s.App.EvmKeeper.ChainID())
			s.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(s.EthSigner, baseFee)
			s.Require().NoError(err, "failed to instantiate Ethereum message")

			// Create StateDB
			s.StateDB = statedb.New(s.Ctx, s.App.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))
			// Instantiate EVM
			evm := s.App.EvmKeeper.NewEVM(
				s.Ctx, msg, cfg, nil, s.StateDB,
			)
			params := s.App.EvmKeeper.GetParams(s.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := s.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			s.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := s.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				s.Require().NoError(err, "expected no error when running the precompile")
				s.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				s.Require().Error(err, "expected error to be returned when running the precompile")
				s.Require().Nil(bz, "expected returned bytes to be nil")
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}
