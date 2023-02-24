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
	"github.com/hajimehoshi/ebiten/v2/audio"
    "log"
    //gMath "math"
)

func AddAlienBullet( ecs *ecs.ECS, 
                   pX, pY float64, 
                   velocity math.Vec2, 
                   view *utility.View, 
                   audioContext *audio.Context) *donburi.Entity {
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
    factions := []component.FactionId{component.Enemy_factionid} //TODO should be arg
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
    nsd.Load("AlienBullet", nil)
    nsd.Play("")
    //nsd.SetPlaySpeed(2.0) //TODO should be constant
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    vd := component.VelocityData{Velocity: &velocity}
    donburi.SetValue(entry, component.Velocity, vd)

    // Actions
    tm := make(map[component.ActionId]bool)
    cdm := make(map[component.ActionId]component.Cooldown)
    am := make(map[component.ActionId]func())

    tm[component.SelfDestruct_actionid] = true
    cdm[component.SelfDestruct_actionid] = component.Cooldown{Cur:500, Max:500}
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
        //hDcopy := *asset.HitD
        hDcopy := *asset.AudioAssets["GenericHit"].DecodedAudio
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

    // Damage
    damageAmount := 1.0
    donburi.SetValue(entry, component.Damage, component.DamageData{
        Value: &damageAmount,
        //DestroyOnDamage: true,
        DestroyOnDamage: true,
    })

    //fDcopy := *asset.FireD
    //TODO alien fire noise?
    fDcopy := *asset.AudioAssets["SciFiProjectile"].DecodedAudio
    firePlayer, err := audioContext.NewPlayer(&fDcopy)
    if err != nil {
        log.Fatal(err)
    }

    firePlayer.Rewind()
    firePlayer.Play()

    return &entity
}