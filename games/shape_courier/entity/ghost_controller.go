package shape_courier_entity

import (
	//gMath "math"
	//"math/rand"
	//"strconv"

	"github.com/infiniteyak/retro_engine/engine/component"
	//"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	//"github.com/yohamta/donburi/features/math"
    //"math"
)

const (
    scatterTime = 2000
    chaseTime = 2000
)

type ghostControllerData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    actions component.ActionsData
    curTime int
    curMode AiMode
}

func AddGhostController(ecs *ecs.ECS) {
    this := &ghostControllerData{}
    this.ecs = ecs

    entity := this.ecs.World.Create(
        component.Actions,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)

    this.curMode = Scatter_aimode //Start in scatter

    // Actions
    this.actions = component.NewActions()
    this.actions.AddUpkeepAction(func(){
        this.curTime++
        if this.curMode == Scatter_aimode && this.curTime >= scatterTime {
            this.curTime = 0
            this.curMode = Chase_aimode
            event.SetAiModeEvent.Publish(this.ecs.World, event.AiMode{Value:int(Chase_aimode)})
            println("controller switch to chase")
        } else if this.curMode == Chase_aimode && this.curTime >= chaseTime {
            this.curTime = 0
            this.curMode = Scatter_aimode
            event.SetAiModeEvent.Publish(this.ecs.World, event.AiMode{Value:int(Scatter_aimode)})
            println("controller switch to scatter")
        }
    })

    donburi.SetValue(this.entry, component.Actions, this.actions)
}
