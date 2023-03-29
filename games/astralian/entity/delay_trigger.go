package astra_entity 

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/component"
)

func AddDelayTrigger(ecs *ecs.ECS, delay int, foo func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.Actions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Actions
    ad := component.NewActions()
    ad.AddCooldownAction(component.TriggerFunction_actionid, delay, foo)
    ad.TriggerMap[component.TriggerFunction_actionid] = true

    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}
