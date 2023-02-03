package game

import (
	//"github.com/infiniteyak/retro_engine/games/astralian/entity"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"strings"
	"github.com/hajimehoshi/ebiten/v2"
)

func (this *Game) LoadTitleScene() {
    println("LoadTitleScene")
    this.curScene.SetId(Title_sceneId)

    this.curScore = utility.ScoreEntry{}
    this.curWave = StartingWave
    this.curShips = StartingShips

    this.GenerateStars(this.screenView) //TODO move these to shared entities?

    entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        this.screenView,
        strings.ToUpper(Title),
    )

    // Advance to the next state when you hit space
    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeySpace,
        func() {
            this.Transition(Advance_sceneEvent)
        },
    )

    // Start game when you hit Enter
    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeyEnter,
        func() {
            this.Transition(GameStart_sceneEvent)
        },
    )
}

