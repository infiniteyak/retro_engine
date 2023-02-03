package entity

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
)

func AddTextInput(ecs *ecs.ECS, str *string, length int, foo func()) *donburi.Entity {
    entity := ecs.World.Create(
        component.TextInput,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Inputs
    donburi.SetValue(entry, component.TextInput, component.TextInputData{
        String: str,
        Length: length,
        Function: foo,
    })

    return &entity
}

