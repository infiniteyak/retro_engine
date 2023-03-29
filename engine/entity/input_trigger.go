package entity

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/component"
)

func AddInputTrigger(ecs *ecs.ECS, key ebiten.Key, triggerFunction func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.Inputs,
        component.Actions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Inputs
    input := component.NewInput()
    input.AddContinuousInput(component.TriggerFunction_actionid, key)
    donburi.SetValue(entry, component.Inputs, input)

    // Actions
    ad := component.NewActions()
    ad.AddCooldownAction(component.TriggerFunction_actionid, 50, func() {
        ad.ResetCooldown(component.TriggerFunction_actionid)
        triggerFunction()
    })
    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}

func AddLimitedInputTrigger(ecs *ecs.ECS, key ebiten.Key, triggerFunction func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.Inputs,
        component.Actions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Inputs
    input := component.NewInput()
    input.AddLimitedInput(component.TriggerFunction_actionid, key)
    donburi.SetValue(entry, component.Inputs, input)

    // Actions
    ad := component.NewActions()
    ad.AddNormalAction(component.TriggerFunction_actionid, func() {
        ad.TriggerMap[component.TriggerFunction_actionid] = false
        triggerFunction()
    })
    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}

func AddHybridInputTrigger(ecs *ecs.ECS, key ebiten.Key, delay int, freq int, triggerFunction func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.Inputs,
        component.Actions,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Inputs
    input := component.NewInput()
    input.AddHybridInput(component.TriggerFunction_actionid, key, delay, freq)
    donburi.SetValue(entry, component.Inputs, input)

    // Actions
    ad := component.NewActions()
    ad.AddNormalAction(component.TriggerFunction_actionid, func() {
        ad.TriggerMap[component.TriggerFunction_actionid] = false
        triggerFunction()
    })
    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}
