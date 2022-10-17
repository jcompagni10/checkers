package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkers module sentinel errors
var (
	ErrInvalidBlack     = sdkerrors.Register(ModuleName, 1100, "black address is invalid:  %s")
	ErrInvalidRed       = sdkerrors.Register(ModuleName, 1101, "red address is invalid:  %s")
	ErrGameNotParseable = sdkerrors.Register(ModuleName, 1102, "Game can't be parsed")
	ErrGameNotFound     = sdkerrors.Register(ModuleName, 1103, "game with id doesnt exist: %s")
	ErrCreatorNotPlayer = sdkerrors.Register(ModuleName, 1104, "message creator isnt a player")
	ErrNotPlayerTurn    = sdkerrors.Register(ModuleName, 1105, "its not your turn buddy: %s")
	ErrWrongMove        = sdkerrors.Register(ModuleName, 1106, "invalid move")

)
