package keeper

import (
	"github.com/jcompagni10/checkers/x/checkers/types"
	"github.com/jcompagni10/checkers/x/checkers/rules"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"context"
	"fmt"


)

func (k Keeper) ForfeitExpiredGames(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	opponents := map[string]string{
    rules.PieceStrings[rules.BLACK_PLAYER]: rules.PieceStrings[rules.RED_PLAYER],
    rules.PieceStrings[rules.RED_PLAYER]:   rules.PieceStrings[rules.BLACK_PLAYER],
	}

	systemInfo, found := k.GetSystemInfo(ctx)
	if !found {
    panic("SystemInfo not found")
	}

	gameIndex := systemInfo.FifoHeadIndex

	var storedGame types.StoredGame



	for {

		if gameIndex == types.NoFifoIndex {
			break
		}

		storedGame, found = k.GetStoredGame(ctx, gameIndex)

		if !found {
			panic("Fifo head game not found " + systemInfo.FifoHeadIndex)
		}
		deadline, err := storedGame.GetDeadlineAsTime()

		if err != nil {
			panic(err)
		}

		if deadline.Before(ctx.BlockTime()) {
			k.RemoveFromFifo(ctx, &storedGame, &systemInfo)

			// handle if worth keeping game
			lastBoard := storedGame.Board

			if storedGame.MoveCount <= 1 {
				k.RemoveStoredGame(ctx, gameIndex)
			} else {
				storedGame.Winner, found = opponents[storedGame.Turn]
				if !found {
					panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Turn))
				}
				storedGame.Board = ""
				k.SetStoredGame(ctx, storedGame)
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.GameForfeitedEventType,
					sdk.NewAttribute(types.GameForfeitedEventGameIndex, gameIndex),
					sdk.NewAttribute(types.GameForfeitedEventWinner, storedGame.Winner),
					sdk.NewAttribute(types.GameForfeitedEventBoard, lastBoard),
				),
			)

			gameIndex = systemInfo.FifoHeadIndex
		} else {
			// All further games active
			break
		}

	}

	k.SetSystemInfo(ctx, systemInfo)



}
