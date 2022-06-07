package proxy

import (
	apptypes "github.com/okex/exchain/app/types"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	supplyexported "github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/ammswap"
	dex "github.com/okex/exchain/x/dex/types"
	distr "github.com/okex/exchain/x/distribution"
	"github.com/okex/exchain/x/farm"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/order"
	"github.com/okex/exchain/x/staking"
	token "github.com/okex/exchain/x/token/types"
	"github.com/okex/exchain/x/wasm/types"
	"github.com/okex/exchain/x/wasm/watcher"
)

const (
	accountBytesLen = 80
)

var gasConfig = types2.KVGasConfig()

// AccountKeeperProxy defines the expected account keeper interface
type AccountKeeperProxy struct {
	cachedAcc map[string]*apptypes.EthAccount
}

func NewAccountKeeperProxy() AccountKeeperProxy {
	return AccountKeeperProxy{}
}

func (a AccountKeeperProxy) SetObserverKeeper(observer auth.ObserverI) {}

func (a AccountKeeperProxy) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	ctx.GasMeter().ConsumeGas(3066, "AccountKeeperProxy NewAccountWithAddress")
	acc := apptypes.EthAccount{
		BaseAccount: &auth.BaseAccount{
			Address: addr,
		},
	}
	return &acc
}

func (a AccountKeeperProxy) GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account) {
	return nil
}

func (a AccountKeeperProxy) IterateAccounts(ctx sdk.Context, cb func(account authexported.Account) bool) {
}

func (a AccountKeeperProxy) GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	ctx.GasMeter().ConsumeGas(gasConfig.ReadCostFlat, types2.GasReadCostFlatDesc)
	ctx.GasMeter().ConsumeGas(gasConfig.ReadCostPerByte*accountBytesLen, types2.GasReadPerByteDesc)
	acc, ok := a.cachedAcc[addr.String()]
	if ok {
		return acc
	}
	return nil
}

func (a AccountKeeperProxy) SetAccount(ctx sdk.Context, account authexported.Account, updateState ...bool) {
	acc, ok := account.(*apptypes.EthAccount)
	if !ok {
		return
	}
	// delay to make
	if a.cachedAcc == nil {
		a.cachedAcc = make(map[string]*apptypes.EthAccount)
	}
	a.cachedAcc[account.GetAddress().String()] = acc
	ctx.GasMeter().ConsumeGas(gasConfig.WriteCostFlat, types2.GasWriteCostFlatDesc)
	ctx.GasMeter().ConsumeGas(gasConfig.WriteCostPerByte*accountBytesLen, types2.GasWritePerByteDesc)
	return
}

func (a AccountKeeperProxy) RemoveAccount(ctx sdk.Context, account authexported.Account) {
	delete(a.cachedAcc, account.GetAddress().String())
	ctx.GasMeter().ConsumeGas(gasConfig.DeleteCost, types2.GasDeleteDesc)
}

type SubspaceProxy struct{}

func (s SubspaceProxy) GetParamSet(ctx sdk.Context, ps params.ParamSet) {
	ctx.GasMeter().ConsumeGas(2111, "SubspaceProxy GetParamSet")
	if wasmParams, ok := ps.(*types.Params); ok {
		wasmParams.CodeUploadAccess = watcher.Params.CodeUploadAccess
		wasmParams.InstantiateDefaultPermission = watcher.Params.InstantiateDefaultPermission
	}
}
func (s SubspaceProxy) SetParamSet(ctx sdk.Context, ps params.ParamSet) {}

type BankKeeperProxy struct {
	blacklistedAddrs map[string]bool
	akp              AccountKeeperProxy
}

func NewBankKeeperProxy(akp AccountKeeperProxy) BankKeeperProxy {
	modAccAddrs := make(map[string]bool)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            nil,
		token.ModuleName:          {supply.Minter, supply.Burner},
		dex.ModuleName:            nil,
		order.ModuleName:          nil,
		ammswap.ModuleName:        {supply.Minter, supply.Burner},
		farm.ModuleName:           nil,
		farm.YieldFarmingAccount:  nil,
		farm.MintFarmingAccount:   {supply.Burner},
	}

	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return BankKeeperProxy{
		blacklistedAddrs: modAccAddrs,
		akp:              akp,
	}
}

func (b BankKeeperProxy) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	//acc := b.akp.GetAccount(ctx, addr)
	//return acc.GetCoins()
	return global.GetSupply()
}

func (b BankKeeperProxy) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	//acc := b.akp.GetAccount(ctx, addr)
	//return sdk.Coin{
	//	Denom:  denom,
	//	Amount: acc.GetCoins().AmountOf(denom),
	//}
	s := global.GetSupply()
	return sdk.Coin{
		Denom:  denom,
		Amount: s.AmountOf(denom),
	}
}

func (b BankKeeperProxy) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	if b.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled
	}
	return nil
}

func (b BankKeeperProxy) GetSendEnabled(ctx sdk.Context) bool {
	ctx.GasMeter().ConsumeGas(1012, "BankKeeperProxy GetSendEnabled")
	return global.GetSendEnabled()
}

func (b BankKeeperProxy) BlockedAddr(addr sdk.AccAddress) bool {
	return b.BlacklistedAddr(addr)
}

func (b BankKeeperProxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return b.blacklistedAddrs[addr.String()]
}

func (b BankKeeperProxy) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	ctx.GasMeter().ConsumeGas(16748, "BankKeeperProxy SendCoins")
	return nil
}

type SupplyKeeperProxy struct{}

func (s SupplyKeeperProxy) GetSupply(ctx sdk.Context) supplyexported.SupplyI {
	return supply.Supply{
		Total: global.GetSupply(),
	}
}

type CapabilityKeeperProxy struct{}

func (c CapabilityKeeperProxy) GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool) {
	return nil, false
}

func (c CapabilityKeeperProxy) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return nil
}

func (c CapabilityKeeperProxy) AuthenticateCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) bool {
	return false
}

type PortKeeperProxy struct{}

func (p PortKeeperProxy) BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability {
	return nil
}
