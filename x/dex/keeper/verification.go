package keeper

import (
	"context"

	"github.com/NicholasDotSol/duality/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) SwapVerification(goCtx context.Context, msg types.MsgSwap) (string, string, sdk.AccAddress, sdk.AccAddress, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// lexographically sort token0, token1
	token0, token1, err := SortTokens(ctx, msg.TokenA, msg.TokenB)

	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(types.ErrInvalidTradingPair, "Not a valid Token Pair: tokenA and tokenB cannot be the same")
	}

	// Converts input address (string) to sdk.AccAddress
	callerAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	// Error checking for the calling address
	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	receiverAddr, err := sdk.AccAddressFromBech32(msg.Receiver)
	// Error Checking for receiver address
	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	if msg.TokenIn != token0 && msg.TokenIn != token1 {
		return "", "", nil, nil, sdkerrors.Wrapf(types.ErrInvalidTradingPair, "TokenIn must be either Tokne0 or Token1")
	}
	// Error checking for valid sdk.Int
	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Not a valid decimal type: %s", err)
	}

	AccountsAmountInBalance := k.bankKeeper.GetBalance(ctx, callerAddr, msg.TokenIn).Amount

	// Error handling to verify the amount wished to deposit is NOT more then the msg.creator holds in their accounts
	if AccountsAmountInBalance.LT(msg.AmountIn) {
		return "", "", nil, nil, sdkerrors.Wrapf(types.ErrNotEnoughCoins, "Address %s  does not have enough of token 0", callerAddr)
	}

	return token0, token1, callerAddr, receiverAddr, nil
}

func (k Keeper) WithdrawLimitOrderVerification(goCtx context.Context, msg types.MsgWithdrawFilledLimitOrder) (string, string, sdk.AccAddress, sdk.AccAddress, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	token0, token1, err := SortTokens(ctx, msg.TokenA, msg.TokenB)

	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(types.ErrInvalidTradingPair, "Not a valid Token Pair: tokenA and tokenB cannot be the same")
	}

	// Converts input address (string) to sdk.AccAddress
	callerAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	// Error checking for the calling address
	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	ReceiverAddr, err := sdk.AccAddressFromBech32(msg.Receiver)
	// Error Checking for receiver address
	// Note we do not actually need to save the sdk.AccAddress here but we do want the address to be checked to determine if it valid
	if err != nil {
		return "", "", nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address (%s)", err)
	}

	pairId := CreatePairId(token0, token1)

	shares, sharesFound := k.GetLimitOrderTrancheUser(ctx, pairId, msg.TickIndex, msg.KeyToken, msg.Key, msg.Receiver)
	if !sharesFound {
		return "", "", nil, nil, sdkerrors.Wrapf(types.ErrNotEnoughShares, "Not enough shares were found")
	}

	// checks that the user has some number of limit order shares wished to withdraw
	if shares.SharesOwned.LTE(sdk.ZeroInt()) {
		return "", "", nil, nil, sdkerrors.Wrapf(types.ErrNotEnoughShares, "Not enough shares were found")
	}

	_ = ctx
	return token0, token1, callerAddr, ReceiverAddr, nil
}
