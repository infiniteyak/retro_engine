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
)

const (
    AlienBulletDamage = 1.0
    AlienBulletHitRadius = 2
    AlienBulletHitOffsetX = 0
    AlienBulletHitOffsetY = 0
    AlienBulletSpriteName = "AlienBullet" 
    AlienBulletDestroyCooldown = 500
    AlienBulletFireSoundName = "SciFiProjectile" 
    AlienBulletDestroySoundName = "GenericHit" 
)

type alienBulletData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    audioContext *audio.Context

    factions component.FactionsData
    damage component.DamageData
    health component.HealthData
    collider component.ColliderData
    position component.PositionData
    view component.ViewData
    velocity component.VelocityData
    graphicObject component.GraphicObjectData
    actions component.ActionsData
}

func AddAlienBullet( ecs *ecs.ECS, 
                   pX, pY float64, 
                   velocity math.Vec2, 
                   view *utility.View, 
                   audioContext *audio.Context) *donburi.Entity {
    abd := &alienBulletData{}
    abd.ecs = ecs

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
    abd.entity = &entity

    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:abd.entity})

    abd.entry = ecs.World.Entry(*abd.entity)

    abd.audioContext = audioContext

    // Factions
    factions := []component.FactionId{component.Enemy_factionid}
    abd.factions = component.FactionsData{Values: factions}
    donburi.SetValue(abd.entry, component.Factions, abd.factions)

    // Collider
    abd.collider = component.NewColliderData()
    hb := component.NewHitbox(AlienBulletHitRadius, AlienBulletHitOffsetX, AlienBulletHitOffsetY)
    abd.collider.Hitboxes = append(abd.collider.Hitboxes, hb)
    donburi.SetValue(abd.entry, component.Collider, abd.collider)

    // Position
    abd.position = component.NewPositionData(pX, pY)
    donburi.SetValue(abd.entry, component.Position, abd.position)

    // Graphic Object
    abd.graphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(AlienBulletSpriteName, nil)
    nsd.Play("")
    abd.graphicObject.Renderables = append(abd.graphicObject.Renderables, &nsd)
    donburi.SetValue(abd.entry, component.GraphicObject, abd.graphicObject)

    // View
    abd.view = component.ViewData{View:view}
    donburi.SetValue(abd.entry, component.View, abd.view)

    // Velocity
    abd.velocity = component.VelocityData{Velocity: &velocity}
    donburi.SetValue(abd.entry, component.Velocity, abd.velocity)

    // Actions
    abd.actions = component.NewActions()

    abd.actions.TriggerMap[component.SelfDestruct_actionid] = true
    cd := component.Cooldown{
        Cur:AlienBulletDestroyCooldown, 
        Max:AlienBulletDestroyCooldown,
    }
    abd.actions.CooldownMap[component.SelfDestruct_actionid] = cd
    abd.actions.ActionMap[component.SelfDestruct_actionid] = func() {
        abd.actions.TriggerMap[component.SelfDestruct_actionid] = false
        abd.actions.TriggerMap[component.DestroySilent_actionid] = true
    }
    abd.actions.ActionMap[component.DestroySilent_actionid] = func() {
        event.RemoveEntityEvent.Publish(
            abd.ecs.World, 
            event.RemoveEntity{Entity:abd.entity},
        )
    }

    abd.actions.ActionMap[component.Destroy_actionid] = func() {
        /*
        hDcopy := *asset.AudioAssets[AlienBulletDestroySoundName].DecodedAudio
        hitPlayer, err := abd.audioContext.NewPlayer(&hDcopy)
        */
        //hitPlayer, err := abd.audioContext.NewPlayer(asset.AudioAssets[AlienBulletDestroySoundName].DecodedAudio)
        asset.PlaySound(audioContext, AlienBulletDestroySoundName)
        /*
        hDcopy := asset.AudioAssets[AlienBulletDestroySoundName].DecodedAudio
        hitPlayer, err := abd.audioContext.NewPlayer(hDcopy)
        if err != nil {
            log.Fatal(err)
        }

        hitPlayer.Rewind()
        hitPlayer.Play()
        */
        event.RemoveEntityEvent.Publish(
            abd.ecs.World, 
            event.RemoveEntity{Entity:abd.entity},
        )
    }

    donburi.SetValue(abd.entry, component.Actions, abd.actions)

    // Damage
    abd.damage = component.NewDamageData()
    *abd.damage.Value = AlienBulletDamage
    *abd.damage.DestroyOnDamage = true
    donburi.SetValue(abd.entry, component.Damage, abd.damage)

    //TODO alien fire noise?
    /*
    fDcopy := *asset.AudioAssets[AlienBulletFireSoundName].DecodedAudio
    firePlayer, err := abd.audioContext.NewPlayer(&fDcopy)
    */
    //firePlayer, err := abd.audioContext.NewPlayer(asset.AudioAssets[AlienBulletFireSoundName].DecodedAudio)
    asset.PlaySound(audioContext, AlienBulletFireSoundName)
    /*
    fDcopy := asset.AudioAssets[AlienBulletFireSoundName].DecodedAudio
    firePlayer, err := abd.audioContext.NewPlayer(fDcopy)
    if err != nil {
        log.Fatal(err)
    }

    firePlayer.Rewind()
    firePlayer.Play()
    */

    return abd.entity
}
