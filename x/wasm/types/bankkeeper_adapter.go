package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

type BankKeeperAdapter struct {
	bankKeeper BankKeeperOKC
}

func NewBankKeeperAdapter(bankKeeper BankKeeperOKC) *BankKeeperAdapter {
	return &BankKeeperAdapter{bankKeeper: bankKeeper}
}
func (adapter BankKeeperAdapter) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return adapter.bankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

func (adapter BankKeeperAdapter) BlockedAddr(addr sdk.AccAddress) bool {
	return adapter.bankKeeper.BlacklistedAddr(addr)
}

func (adapter BankKeeperAdapter) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	if !adapter.bankKeeper.GetSendEnabled(ctx) {
		return sdkerrors.Wrapf(ErrSendDisabled, "transfers are currently disabled")
	}
	// todo weather allow different form okt coin send enable
	return nil
}

func (adapter BankKeeperAdapter) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return adapter.bankKeeper.GetCoins(ctx, addr)
}

func (adapter BankKeeperAdapter) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	coins := adapter.bankKeeper.GetCoins(ctx, addr)
	return sdk.Coin{
		Amount: coins.AmountOf(denom),
		Denom:  denom,
	}
}