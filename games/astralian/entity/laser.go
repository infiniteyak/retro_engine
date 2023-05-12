package astra_entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/yohamta/donburi/features/math"
)

func AddLaser( ecs *ecs.ECS, 
               pX, pY float64, 
               velocity math.Vec2, 
               view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Velocity,
        component.Actions,
        component.Collider,
        component.Factions,
        component.Damage,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Factions
    factions := []component.FactionId{component.Player_factionid} //TODO should be arg
    donburi.SetValue(entry, component.Factions, factions)

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(2, 0, -2))
    donburi.SetValue(entry, component.Collider, collider)

    // Position
    pd := component.NewPositionData(pX, pY)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load("Laser", nil)
    nsd.Play("")
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    vd := component.VelocityData{Velocity: &velocity}
    donburi.SetValue(entry, component.Velocity, vd)

    // Actions
    ad := component.NewActions()

    ad.AddCooldownAction(component.SelfDestruct_actionid, 50, func(){
        ad.TriggerMap[component.SelfDestruct_actionid] = false
        ad.TriggerMap[component.DestroySilent_actionid] = true
    })
    ad.TriggerMap[component.SelfDestruct_actionid] = true

    ad.AddNormalAction(component.DestroySilent_actionid, func(){
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    })

    ad.AddNormalAction(component.Destroy_actionid, func(){
        asset.PlaySound("GenericHit")

        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    })

    donburi.SetValue(entry, component.Actions, ad)

    // Damage
    dd := component.NewDamageData(1.0)
    *dd.DestroyOnDamage = true
    donburi.SetValue(entry, component.Damage, dd)

    asset.PlaySound("SciFiProjectile")

    return &entity
}
