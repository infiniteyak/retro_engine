package entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/tanema/gween"
)

func AddTweenSprite(ecs *ecs.ECS, 
                    layer ecs.LayerID, 
                    x, y float64, 
                    xTween, yTween *gween.Tween,
                    delay int, 
                    spriteName string, 
                    view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer, 
        component.Position, 
        component.GraphicObject,
        component.View,
        component.PosTween,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(spriteName, nil)
    nsd.Play("Idle")
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // PosTween
    comp := component.PosTweenData{
        XTween: xTween,
        YTween: yTween,
        Delay: delay,
    }
    donburi.SetValue(entry, component.PosTween, comp)

    return &entity
}
