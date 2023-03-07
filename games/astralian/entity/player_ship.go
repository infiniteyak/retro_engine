package astra_entity

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/asset"
    dmath "github.com/yohamta/donburi/features/math"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

func AddPlayerShip( ecs *ecs.ECS, 
                    x, y float64, 
                    view *utility.View, 
                    audioContext *audio.Context) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.ViewBound,
        component.GraphicObject,
        component.Inputs,
        component.Actions,
        component.Velocity,
        component.Collider,
        component.Health,
        component.Factions,
        component.Damage,
        )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Damage
    dD := component.NewDamageData()
    *dD.Value = 1.0
    donburi.SetValue(entry, component.Damage, dD)

    // Factions
    factions := []component.FactionId{component.Player_factionid}
    donburi.SetValue(entry, component.Factions, factions)

    // Health
    healthAmount := 1.0
    donburi.SetValue(entry, component.Health, component.HealthData{Value:&healthAmount})

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(5, 0, -2))
    donburi.SetValue(entry, component.Collider, collider)

    // Velocity
    vd := component.VelocityData{Velocity: &dmath.Vec2{}}
    donburi.SetValue(entry, component.Velocity, vd)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // ViewBound
    donburi.SetValue(entry, component.ViewBound, component.ViewBoundData{
        XDistance: 7.0,
    })

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    shipSd := component.SpriteData{}
    shipSd.Load("AstralianShip", nil)
    shipSd.Play("Ready") //TODO add this to other areas?
    gobj.Renderables = append(gobj.Renderables, &shipSd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // Inputs
    im := make(map[component.ActionId]ebiten.Key)
    im[component.MoveLeft_actionid] = ebiten.KeyLeft
    im[component.MoveRight_actionid] = ebiten.KeyRight
    im[component.Shoot_actionid] = ebiten.KeySpace
    im[component.ShootSecondary_actionid] = ebiten.KeyControl
    donburi.SetValue(entry, component.Inputs, component.InputData{Mapping: im})

    // Actions
    tm := make(map[component.ActionId]bool)
    tm[component.Shield_actionid] = true //start out invulnerable
    cdm := make(map[component.ActionId]component.Cooldown)
    cdm[component.Shoot_actionid] = component.Cooldown{Cur:50, Max:50}
    cdm[component.Shield_actionid] = component.Cooldown{Cur:300, Max:300}
    am := make(map[component.ActionId]func())

    // Shoot
    bulletVelocity := dmath.Vec2{X:0, Y:-2.0}
    readyToFire := true
    power := 1
    am[component.Shoot_actionid] = func() {
        max := cdm[component.Shoot_actionid].Max
        cooldown := component.Cooldown{Cur:max, Max:max}
        cdm[component.Shoot_actionid] = cooldown

        if readyToFire {
            AddBoomerang(ecs, pd.Point.X, pd.Point.Y, bulletVelocity, view, audioContext, power, &entity)
            readyToFire = false
            shipSd.Play("Idle")
        }
    }

    am[component.Reload_actionid] = func() {
        tm[component.Reload_actionid] = false
        readyToFire = true
        shipSd.Play("Ready")
    }

    cdm[component.ShootSecondary_actionid] = component.Cooldown{Cur:50, Max:75}
    secondaryBulletVelocity := dmath.Vec2{X:0, Y:-2.0}
    am[component.ShootSecondary_actionid] = func() {
        max := cdm[component.ShootSecondary_actionid].Max
        cooldown := component.Cooldown{Cur:max, Max:max}
        cdm[component.ShootSecondary_actionid] = cooldown

        AddLaser(ecs, pd.Point.X-3, pd.Point.Y-4, secondaryBulletVelocity, view, audioContext)
        AddLaser(ecs, pd.Point.X+3, pd.Point.Y-4, secondaryBulletVelocity, view, audioContext)
    }

    am[component.IncreasePower_actionid] = func() {
        tm[component.IncreasePower_actionid] = false
        power++
    }

    am[component.ResetPower_actionid] = func() {
        tm[component.ResetPower_actionid] = false
        if power > 1 {
            power -= 1
        }
    }

    // Shield - actually turns off shield
    am[component.Shield_actionid] = func() {
        tm[component.Shield_actionid] = false
    }

    // Move Left
    moveSpeed := 1.0
    am[component.MoveLeft_actionid] = func() {
        vd.Velocity.X = -1.0 * moveSpeed
    }

    // Move Right
    am[component.MoveRight_actionid] = func() {
        vd.Velocity.X = moveSpeed
    }

    am[component.Destroy_actionid] = func() {
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
        //shipDestroyedDcopy := *asset.DestroyedD
        /*
        shipDestroyedDcopy := *asset.AudioAssets["PlayerShipDestroyed"].DecodedAudio
        destroyedPlayer, err := audioContext.NewPlayer(&shipDestroyedDcopy)
        */
        asset.PlaySound(audioContext, "PlayerShipDestroyed")
        /*
        destroyedPlayer, err := audioContext.NewPlayer(asset.AudioAssets["PlayerShipDestroyed"].DecodedAudio)
        if err != nil {
            log.Fatal(err)
        }

        destroyedPlayer.Rewind()
        destroyedPlayer.Play()
        */

        event.ShipDestroyedEvent.Publish(ecs.World, event.ShipDestroyed{})

        AddExplosion(ecs, pd.Point.X, pd.Point.Y, "AstralianShip", view)
    }

    blinkCounter := 0
    am[component.Upkeep_actionid] = func() {
        // do a blinking effect if we are shielded
        if tm[component.Shield_actionid] {
            blinkCounter++
            if (blinkCounter / 10) % 2 == 0 {
                *shipSd.RenderableData.GetTransInfo().Hide = true
            } else {
                *shipSd.RenderableData.GetTransInfo().Hide = false
            }
        } else {
            *shipSd.RenderableData.GetTransInfo().Hide = false
        }
        // if they're both on or both off...
        if tm[component.MoveRight_actionid] == tm[component.MoveLeft_actionid] {
            vd.Velocity.X = 0
        }
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}
