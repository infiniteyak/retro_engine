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
    tm := make(map[component.ActionId]bool)
    tm[component.TriggerFunction_actionid] = true
    cdm := make(map[component.ActionId]component.Cooldown)
    cdm[component.TriggerFunction_actionid] = component.Cooldown{Cur:delay, Max:delay}
    am := make(map[component.ActionId]func())
    
    // Advance to next screen
    am[component.TriggerFunction_actionid] = foo

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}
