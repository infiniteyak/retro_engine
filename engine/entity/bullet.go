package entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/yohamta/donburi/features/math"
	"github.com/hajimehoshi/ebiten/v2/audio"
    "log"
)

func AddBullet(ecs *ecs.ECS, pX, pY, wDist float64, velocity math.Vec2, view *utility.View, audioContext *audio.Context) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Velocity,
        component.Wrap,
        component.Actions,
        component.Collider,
        component.Factions,
        component.Damage,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Damage
    donburi.SetValue(entry, component.Damage, component.DamageData{
        Value: 1.0,
        DestroyOnDamage: true,
    })

    // Factions
    factions := []component.FactionId{component.Player_factionid} //TODO should be arg
    donburi.SetValue(entry, component.Factions, factions)

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(2, 0, 0))
    donburi.SetValue(entry, component.Collider, collider)

    // Position
    pd := component.NewPositionData(pX, pY)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load("SimpleBullet", nil)
    nsd.SetPlaySpeed(0.2) //TODO should be constant
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    vd := component.VelocityData{Velocity: &velocity}
    donburi.SetValue(entry, component.Velocity, vd)

    // Wrap
    wrap := component.WrapData{Distance: new(float64)}
    *wrap.Distance = wDist
    donburi.SetValue(entry, component.Wrap, wrap)

    // Actions
    tm := make(map[component.ActionId]bool)
    cdm := make(map[component.ActionId]component.Cooldown)
    am := make(map[component.ActionId]func())

    tm[component.SelfDestruct_actionid] = true
    cdm[component.SelfDestruct_actionid] = component.Cooldown{Cur:100, Max:100}
    am[component.SelfDestruct_actionid] = func() {
        tm[component.SelfDestruct_actionid] = false
        tm[component.DestroySilent_actionid] = true
    }
    am[component.DestroySilent_actionid] = func() {
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    }

    am[component.Destroy_actionid] = func() {
        AddSmallExplosion(
            ecs, 
            pd.Point.X, 
            pd.Point.Y, 
            view,
        )
        hDcopy := *asset.HitD
        hitPlayer, err := audioContext.NewPlayer(&hDcopy)
        if err != nil {
            log.Fatal(err)
        }

        hitPlayer.Rewind()
        hitPlayer.Play()
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    fDcopy := *asset.FireD
    firePlayer, err := audioContext.NewPlayer(&fDcopy)
    if err != nil {
        log.Fatal(err)
    }

    firePlayer.Rewind()
    firePlayer.Play()

    return &entity
}
