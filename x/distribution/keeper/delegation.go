package keeper

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/types"
	stakingexported "github.com/okex/exchain/x/staking/exported"
)

// initialize starting info for a new delegation
func (k Keeper) initializeDelegation(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	if !tmtypes.HigherThanVenus2(ctx.BlockHeight()) {
		return
	}

	logger := k.Logger(ctx)
	// period has already been incremented - we want to store the period ended by this delegation action
	previousPeriod := k.GetValidatorCurrentRewards(ctx, val).Period - 1

	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val, previousPeriod)

	delegation := k.stakingKeeper.Delegator(ctx, del)
	logger.Debug(fmt.Sprintf("initialize delegation, val:%s, del:%s, shares:%s", val.String(), del.String(), delegation.GetLastAddedShares().String()))
	k.SetDelegatorStartingInfo(ctx, val, del, types.NewDelegatorStartingInfo(previousPeriod, delegation.GetLastAddedShares(), uint64(ctx.BlockHeight())))
}

// calculate the rewards accrued by a delegation between two periods
func (k Keeper) calculateDelegationRewardsBetween(ctx sdk.Context, val stakingexported.ValidatorI,
	startingPeriod, endingPeriod uint64, stake sdk.Dec) (rewards sdk.DecCoins) {
	logger := k.Logger(ctx)
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), startingPeriod)
	ending := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(stake)
	logger.Debug(fmt.Sprintf("starting ratio:%s, ending:%s, diff:%s, stake:%s, rewards:%s",
		starting.CumulativeRewardRatio.String(), ending.CumulativeRewardRatio.String(), difference.String(),
		stake.String(), rewards.String()))
	return
}

// calculate the total rewards accrued by a delegation
func (k Keeper) calculateDelegationRewards(ctx sdk.Context, val stakingexported.ValidatorI, delAddr sdk.AccAddress, endingPeriod uint64) (rewards sdk.DecCoins) {
	logger := k.Logger(ctx)
	del := k.stakingKeeper.Delegator(ctx, delAddr)

	// fetch starting info for delegation
	startingInfo := k.GetDelegatorStartingInfo(ctx, val.GetOperator(), del.GetDelegatorAddress())

	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		// started this height, no rewards yet
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake
	if stake.GT(del.GetLastAddedShares()) {
		panic(fmt.Sprintf("calculated final stake for delegator %s greater than current stake"+
			"\n\tfinal stake:\t%s"+
			"\n\tcurrent stake:\t%s",
			del.GetDelegatorAddress(), stake, del.GetLastAddedShares()))
	}

	// calculate rewards for final period
	rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)

	logger.Debug(fmt.Sprintf("calculate delegation rewards, val:%s, del:%s, period:[%d,%d], stake:%s,reward:%s",
		val.GetOperator().String(), delAddr.String(), startingPeriod, endingPeriod, stake.String(), rewards.String()))

	return rewards
}

//withdraw rewards according to the specified validator by delegator
func (k Keeper) withdrawDelegationRewards(ctx sdk.Context, val stakingexported.ValidatorI, delAddress sdk.AccAddress) (sdk.Coins, error) {
	if !tmtypes.HigherThanVenus2(ctx.BlockHeight()) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidVersion, "not support")
	}

	// check existence of delegator starting info
	if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), delAddress) {
		del := k.stakingKeeper.Delegator(ctx, delAddress)
		if del.GetLastAddedShares().IsZero() {
			return nil, types.ErrCodeEmptyDelegationDistInfo()
		}
		k.initDelegationStartInfo(ctx, val, del)
	}

	// end current period and calculate rewards
	endingPeriod := k.incrementValidatorPeriod(ctx, val)
	rewardsRaw := k.calculateDelegationRewards(ctx, val, delAddress, endingPeriod)
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())

	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(fmt.Sprintf("missing rewards rounding error, delegator %v"+
			"withdrawing rewards from validator %v, should have received %v, got %v",
			val.GetOperator(), delAddress, rewardsRaw, rewards))
	}

	// truncate coins, return remainder to community pool
	coins, remainder := rewards.TruncateDecimal()

	// add coins to user account
	if !coins.IsZero() {
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, delAddress)
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding.Sub(rewards))
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
	k.SetFeePool(ctx, feePool)

	// decrement reference count of starting period
	startingInfo := k.GetDelegatorStartingInfo(ctx, val.GetOperator(), delAddress)
	startingPeriod := startingInfo.PreviousPeriod
	k.decrementReferenceCount(ctx, val.GetOperator(), startingPeriod)

	// remove delegator starting info
	k.DeleteDelegatorStartingInfo(ctx, val.GetOperator(), delAddress)

	return coins, nil
}

func (k Keeper) initDelegationStartInfo(ctx sdk.Context, val stakingexported.ValidatorI, del stakingexported.DelegatorI) {
	if !tmtypes.HigherThanVenus2(ctx.BlockHeight()) {
		return
	}

	logger := k.Logger(ctx)
	//If the delegator a shares but no start info,
	//it add shares before distribution proposal, and need to set a new start info

	//set previous validator period 0
	previousPeriod := uint64(0)
	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val.GetOperator(), previousPeriod)

	logger.Debug(fmt.Sprintf("initialize delegation, val:%s, del:%s, shares:%s",
		val.GetOperator().String(), del.GetDelegatorAddress().String(), del.GetLastAddedShares().String()))
	k.SetDelegatorStartingInfo(ctx, val.GetOperator(), del.GetDelegatorAddress(), types.NewDelegatorStartingInfo(previousPeriod, del.GetLastAddedShares(), uint64(ctx.BlockHeight())))
	return
}