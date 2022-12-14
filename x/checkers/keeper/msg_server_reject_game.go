package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jcompagni10/checkers/x/checkers/types"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)


	refund := uint64(types.RejectGameRefundGas)
	if consumed := ctx.GasMeter().GasConsumed(); consumed < refund {
    refund = consumed
	}
	ctx.GasMeter().RefundGas(refund, "Reject game")
	
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
    return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	if storedGame.Black == msg.Creator {
    if 0 < storedGame.MoveCount { 
			return nil, types.ErrBlackAlreadyPlayed
    }
	} else if storedGame.Red == msg.Creator {
    if 1 < storedGame.MoveCount {
			return nil, types.ErrRedAlreadyPlayed
    }
	} else {
    return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	}

	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	k.Keeper.RemoveFromFifo(ctx, &storedGame, &systemInfo)

	k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)

	k.Keeper.SetSystemInfo(ctx, systemInfo)

	k.Keeper.MustRefundWager(ctx, &storedGame)


	ctx.EventManager().EmitEvent(
    sdk.NewEvent(types.GameRejectedEventType,
			sdk.NewAttribute(types.GameRejectedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameRejectedEventGameIndex, msg.GameIndex),
    ),
	)

	return &types.MsgRejectGameResponse{}, nil
}
