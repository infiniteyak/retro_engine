package entity

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	//"github.com/hajimehoshi/ebiten/v2"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/component"
)

func AddTimer(ecs *ecs.ECS, delay int, foo func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.Actions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Actions
    curTime := 0
    ad := component.NewActions()
    ad.AddUpkeepAction(func(){
        curTime++
        if curTime >= delay {
            foo()
            ree := event.RemoveEntity{Entity:&entity}
            event.RemoveEntityEvent.Publish(ecs.World, ree)
        }
    })
    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}

