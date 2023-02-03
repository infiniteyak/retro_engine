package entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/yohamta/donburi/features/math"
    "math/rand"
)

func AddStar(ecs *ecs.ECS, pX, pY float64, view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Velocity,
        component.Wrap,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Position
    pd := component.NewPositionData(pX, pY)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load("Star", nil)
    //nsd.SetPlaySpeed(0.1) //TODO should be constant
    nsd.SetPlaySpeed(float32(rand.Float64()*0.1 + 0.03)) //TODO should be constant
    nsd.SetFrame(rand.Intn(4)) //TODO should be constant
    nsd.Play("")
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    vd := component.VelocityData{Velocity: &math.Vec2{X:0, Y:float64(rand.Float64()*0.5 + 0.3)}}
    donburi.SetValue(entry, component.Velocity, vd)

    // Wrap
    wrap := component.WrapData{Distance: new(float64)}
    *wrap.Distance = 2.0
    donburi.SetValue(entry, component.Wrap, wrap)

    return &entity
}
