package keeper

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibctransferType "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	ibcclienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/x/erc20/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// OnMintVouchers after minting vouchers on this chain, convert these vouchers into evm tokens.
func (k Keeper) OnMintVouchers(ctx sdk.Context, vouchers sdk.SysCoins, receiver string) {
	cacheCtx, commit := ctx.CacheContext()
	err := k.ConvertVouchers(cacheCtx, receiver, vouchers)
	if err != nil {
		logrusplugin.Error(
			fmt.Sprintf("Failed to convert vouchers to evm tokens for receiver %s, coins %s. Receive error %s",
				receiver, vouchers.String(), err))
	}
	commit()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
}

// ConvertVouchers convert vouchers into native coins or evm tokens.
func (k Keeper) ConvertVouchers(ctx sdk.Context, from string, vouchers sdk.SysCoins) error {
	fromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	params := k.GetParams(ctx)
	for _, c := range vouchers {
		switch c.Denom {
		case params.IbcDenom:
			if len(params.IbcDenom) == 0 {
				return errors.New("ibc denom is empty")
			}
			// oec1:okt----->oec2:ibc/okt---->oec2:okt
			if err := k.ConvertVoucherToEvmDenom(ctx, fromAddr, c); err != nil {
				return err
			}
		default:
			// oec1:xxb----->oec2:ibc/xxb---->oec2:erc20/xxb
			if err := k.ConvertVoucherToERC20(ctx, fromAddr, c, params.EnableAutoDeployment); err != nil {
				return err
			}
		}
	}
	return nil
}

// ConvertVoucherToEvmDenom convert vouchers into evm denom.
func (k Keeper) ConvertVoucherToEvmDenom(ctx sdk.Context, from sdk.AccAddress, voucher sdk.SysCoin) error {
	logrusplugin.Info("convert voucher into evm token", "from", from.String(), "voucher", voucher.String())
	// 1. send voucher to escrow address
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(voucher)); err != nil {
		return err
	}

	evmCoin := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, voucher.Amount)
	// 2. Mint evm token
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(evmCoin)); err != nil {
		return err
	}
	// 3. Send evm token to receiver
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoins(evmCoin)); err != nil {
		return err
	}
	return nil
}

// ConvertVoucherToERC20 convert vouchers into evm tokens.
func (k Keeper) ConvertVoucherToERC20(ctx sdk.Context, from sdk.AccAddress, voucher sdk.SysCoin, autoDeploy bool) error {
	logrusplugin.Info("convert vouchers into evm tokens",
		"fromBech32", from.String(),
		"fromEth", common.BytesToAddress(from.Bytes()).String(),
		"voucher", voucher.String())

	if !types.IsValidIBCDenom(voucher.Denom) {
		return fmt.Errorf("coin %s is not supported for wrapping", voucher.Denom)
	}

	var err error
	contract, found := k.GetContractByDenom(ctx, voucher.Denom)
	if !found {
		// automated deployment contracts
		if !autoDeploy {
			return fmt.Errorf("no contract found for the denom %s", voucher.Denom)
		}
		contract, err = k.deployModuleERC20(ctx, voucher.Denom)
		if err != nil {
			return err
		}
		k.SetAutoContractForDenom(ctx, voucher.Denom, contract)
		logrusplugin.Info("contract created for coin", "contract", contract.String(), "denom", voucher.Denom)
	}
	// 1. transfer voucher from user address to contact address in bank
	if err := k.bankKeeper.SendCoins(ctx, from, sdk.AccAddress(contract.Bytes()), sdk.NewCoins(voucher)); err != nil {
		return err
	}
	// 2. call contract, mint token to user address in contract
	ac, err := sdk.ConvertDecCoinToAdapterCoin(voucher)
	if err != nil {
		return err
	}
	if _, err := k.callModuleERC20(
		ctx,
		contract,
		"mint_by_oec_module",
		common.BytesToAddress(from.Bytes()),
		ac.Amount.BigInt()); err != nil {
		return err
	}
	return nil
}

// deployModuleERC20 deploy an embed erc20 contract
func (k Keeper) deployModuleERC20(ctx sdk.Context, denom string) (common.Address, error) {
	byteCode := common.Hex2Bytes(types.ModuleERC20Contract.Bin)
	input, err := types.ModuleERC20Contract.ABI.Pack("", denom, uint8(0))
	if err != nil {
		return common.Address{}, err
	}

	data := append(byteCode, input...)
	_, res, err := k.callEvmByModule(ctx, nil, big.NewInt(0), data)
	if err != nil {
		return common.Address{}, err
	}
	return res.ContractAddress, nil
}

// callModuleERC20 call a method of ModuleERC20 contract
func (k Keeper) callModuleERC20(ctx sdk.Context, contract common.Address, method string, args ...interface{}) ([]byte, error) {
	data, err := types.ModuleERC20Contract.ABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	_, _, err = k.callEvmByModule(ctx, &contract, big.NewInt(0), data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// callEvmByModule execute an evm message from native module
func (k Keeper) callEvmByModule(ctx sdk.Context, to *common.Address, value *big.Int, data []byte) (*evmtypes.ExecutionResult, *evmtypes.ResultData, error) {
	config, found := k.evmKeeper.GetChainConfig(ctx)
	if !found {
		return nil, nil, types.ErrChainConfigNotFound
	}

	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, nil, err
	}

	nonce := uint64(0)
	acc := k.accountKeeper.GetAccount(ctx, types.EVMModuleBechAddr)
	if acc != nil {
		nonce = acc.GetSequence()
	}
	st := evmtypes.StateTransition{
		AccountNonce: nonce,
		Price:        big.NewInt(0),
		GasLimit:     evmtypes.DefaultMaxGasLimitPerTx,
		Recipient:    to,
		Amount:       value,
		Payload:      data,
		Csdb:         evmtypes.CreateEmptyCommitStateDB(k.evmKeeper.GenerateCSDBParams(), ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &common.Hash{},
		Sender:       types.EVMModuleETHAddr,
		Simulate:     ctx.IsCheckTx(),
		TraceTx:      false,
		TraceTxLog:   false,
	}

	executionResult, resultData, err, _, _ := st.TransitionDb(ctx, config)
	return executionResult, resultData, err
}

// IbcTransferVouchers transfer vouchers to other chain by ibc
func (k Keeper) IbcTransferVouchers(ctx sdk.Context, from, to string, vouchers sdk.SysCoins) error {
	fromAddr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	if len(to) == 0 {
		return errors.New("to address cannot be empty")
	}
	logrusplugin.Info("transfer vouchers to other chain by ibc", "from", from, "to", to)
	//params := k.GetParams(ctx)
	for _, c := range vouchers {
		switch c.Denom {
		case sdk.DefaultBondDenom:
			// oec2:okt----->oec2:ibc/okt---ibc--->oec1:okt
			if err := k.ibcSendEvmDenom(ctx, fromAddr, to, c); err != nil {
				return err
			}
		default:
			if _, found := k.GetContractByDenom(ctx, c.Denom); !found {
				return fmt.Errorf("coin %s id not support", c.Denom)
			}
			// oec2:erc20/xxb----->oec2:ibc/xxb---ibc--->oec1:xxb
			if err := k.ibcSendTransfer(ctx, fromAddr, to, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) ibcSendEvmDenom(ctx sdk.Context, sender sdk.AccAddress, to string, coin sdk.Coin) error {
	// TODO Not supported at the moment
	// 1. Send evm token to escrow address
	// 2. Burn the evm token
	// 3. Send ibc coin from module account to sender
	// 4. Send ibc coin to ibc

	return nil
}

func (k Keeper) ibcSendTransfer(ctx sdk.Context, sender sdk.AccAddress, to string, coin sdk.Coin) error {
	// Coin needs to be a voucher so that we can extract the channel id from the denom
	channelID, err := k.GetSourceChannelID(ctx, coin.Denom)
	if err != nil {
		return err
	}

	ac, err := sdk.ConvertDecCoinToAdapterCoin(coin)
	if err != nil {
		return err
	}
	// Transfer coins to receiver through IBC
	// We use current time for timeout timestamp and zero height for timeoutHeight
	// it means it can never fail by timeout
	params := k.GetParams(ctx)
	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + params.IbcTimeout
	timeoutHeight := ibcclienttypes.ZeroHeight()

	return k.transferKeeper.SendTransfer(
		ctx,
		ibctransferType.PortID,
		channelID,
		ac,
		sender,
		to,
		timeoutHeight,
		timeoutTimestamp,
	)
}
