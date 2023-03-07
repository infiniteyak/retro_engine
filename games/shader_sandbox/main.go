package main

import (
    "github.com/infiniteyak/retro_engine/games/shader_sandbox/game"
    "github.com/infiniteyak/retro_engine/engine/constants"
    "github.com/infiniteyak/retro_engine/engine"
)

func main() {
    g := game.NewGame(constants.ArcadeBTateWidth, constants.ArcadeBTateHeight)
    engine.RunGame(g, game.Title)
}
