package main

import (
    "github.com/infiniteyak/retro_engine/games/demo/game"
    "github.com/infiniteyak/retro_engine/engine/constants"
    "github.com/infiniteyak/retro_engine/engine"
)

const (
    Title = "Demo Game"
)

func main() {
    g := game.NewGame(constants.ArcadeATateWidth, constants.ArcadeATateHeight)
    engine.RunGame(g, Title)
}
