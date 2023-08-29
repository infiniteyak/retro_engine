package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/games/shape_courier/entity"
	"strings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
)

func (this *Game) LoadAttractModeScene() {
    println("LoadAttractModeScene")
    this.curScene.SetId(Attract_sceneId)

    adjustDots := func(w donburi.World, e event.AdjustDots) {
    }
    event.AdjustDotsEvent.Subscribe(this.ecs.World, adjustDots)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.AdjustDotsEvent.Unsubscribe(this.ecs.World, adjustDots)
            },
        },
    )
    shape_courier_entity.AddDecorativeMaze(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2), //TODO these do nothing?? 
        this.screenView)

    entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2) - 1, 
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

