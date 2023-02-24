package astra_entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
)

func AddExplosion(ecs *ecs.ECS, x, y float64, sprite string, view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.GraphicObject,
        component.View,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(sprite, nil)
    nsd.Play("Explosion")
    nsd.SetLoopCallback(func() {
        //Destroy might not happen as quickly as I'd like so hide in order
        // to prevent the animation looping
        *nsd.RenderableData.GetTransInfo().Hide = true
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    })
    nsd.SetPlaySpeed(1)
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    return &entity
}
