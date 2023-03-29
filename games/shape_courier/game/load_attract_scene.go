package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
	"strings"
	"github.com/hajimehoshi/ebiten/v2"
)

func (this *Game) LoadAttractModeScene() {
    println("LoadAttractModeScene")
    this.curScene.SetId(Attract_sceneId)

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
}

