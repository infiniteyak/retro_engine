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

func AddMediumAsteroid(ecs *ecs.ECS, x, y float64, view *utility.View, doCount bool) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Velocity,
        component.Wrap,
        component.Collider,
        component.Health,
        component.Factions,
        component.Actions,
        component.Damage,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Damage
    donburi.SetValue(entry, component.Damage, component.DamageData{
        Value: 1.0,
    })

    // Factions
    factions := []component.FactionId{component.Enemy_factionid}
    donburi.SetValue(entry, component.Factions, factions)

    // Health
    donburi.SetValue(entry, component.Health, component.HealthData{Value:1.0})

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(8, 0, 0))
    donburi.SetValue(entry, component.Collider, collider)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load("MediumAsteroid", nil)
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    const (
        minVelocity = -0.5
        maxVelocity = 0.5
    )
    //rand.Seed(time.Now().UnixNano())
    xVel := minVelocity + rand.Float64() * (maxVelocity - minVelocity)
    yVel := minVelocity + rand.Float64() * (maxVelocity - minVelocity)
    vd := component.VelocityData{Velocity: &math.Vec2{X:xVel, Y:yVel}}
    donburi.SetValue(entry, component.Velocity, vd)

    // Wrap
    wrap := component.WrapData{Distance: new(float64)}
    *wrap.Distance = 10.0
    donburi.SetValue(entry, component.Wrap, wrap)

    tm := make(map[component.ActionId]bool)
    cdm := make(map[component.ActionId]component.Cooldown)
    am := make(map[component.ActionId]func())
    am[component.Destroy_actionid] = func() {
        event.ScoreEvent.Publish(ecs.World, event.Score{Value:5})
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
        AddSmallAsteroid(
            ecs, 
            pd.Point.X, 
            pd.Point.Y, 
            view,
            doCount,
        )
        AddSmallAsteroid(
            ecs, 
            pd.Point.X, 
            pd.Point.Y, 
            view,
            doCount,
        )
        event.AsteroidsCountUpdateEvent.Publish(ecs.World, event.AsteroidsCountUpdate{Value:-1})
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    if doCount {
        event.AsteroidsCountUpdateEvent.Publish(ecs.World, event.AsteroidsCountUpdate{Value:1})
    }

    return &entity
}
