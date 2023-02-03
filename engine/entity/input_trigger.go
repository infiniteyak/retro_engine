package entity

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/component"
)

func AddInputTrigger(ecs *ecs.ECS, key ebiten.Key, foo func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.Inputs,
        component.Actions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Inputs
    im := make(map[component.ActionId]ebiten.Key)
    im[component.TriggerFunction_actionid] = key //ebiten.KeySpace //should this be configurable?
    donburi.SetValue(entry, component.Inputs, component.InputData{Mapping: im})

    // Actions
    tm := make(map[component.ActionId]bool)
    cdm := make(map[component.ActionId]component.Cooldown)
    cdm[component.TriggerFunction_actionid] = component.Cooldown{Cur:100, Max:50}
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
