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
)

func AddPlayerShip( ecs *ecs.ECS, 
                    x, y float64, 
                    view *utility.View) *donburi.Entity {
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
    inputs := component.NewInput()
    inputs.AddContinuousInput(component.MoveLeft_actionid, ebiten.KeyLeft)
    inputs.AddContinuousInput(component.MoveRight_actionid, ebiten.KeyRight)
    inputs.AddContinuousInput(component.Shoot_actionid, ebiten.KeySpace)
    inputs.AddContinuousInput(component.ShootSecondary_actionid, ebiten.KeyControl)
    donburi.SetValue(entry, component.Inputs, inputs)

    // Actions
    ad := component.NewActions()

    // Shoot
    bulletVelocity := dmath.Vec2{X:0, Y:-2.0}
    readyToFire := true
    power := 1
    ad.AddCooldownAction(component.Shoot_actionid, 50, func(){
        ad.ResetCooldown(component.Shoot_actionid)

        if readyToFire {
            AddBoomerang(ecs, pd.Point.X, pd.Point.Y, bulletVelocity, view, power, &entity)
            readyToFire = false
            shipSd.Play("Idle")
        }
    })

    ad.AddNormalAction(component.Reload_actionid, func(){
        ad.TriggerMap[component.Reload_actionid] = false
        readyToFire = true
        shipSd.Play("Ready")
    })

    secondaryBulletVelocity := dmath.Vec2{X:0, Y:-2.0}
    ad.AddCooldownAction(component.ShootSecondary_actionid, 75, func(){
        ad.ResetCooldown(component.ShootSecondary_actionid)

        AddLaser(ecs, pd.Point.X-3, pd.Point.Y-4, secondaryBulletVelocity, view)
        AddLaser(ecs, pd.Point.X+3, pd.Point.Y-4, secondaryBulletVelocity, view)
    })
    ad.SetCooldown(component.ShootSecondary_actionid, 50)

    ad.AddNormalAction(component.IncreasePower_actionid, func(){
        ad.TriggerMap[component.IncreasePower_actionid] = false
        power++
    })

    ad.AddNormalAction(component.ResetPower_actionid, func(){
        ad.TriggerMap[component.ResetPower_actionid] = false
        if power > 1 {
            power -= 1
        }
    })

    // Shield - actually turns off shield
    ad.AddCooldownAction(component.Shield_actionid, 300, func(){
        ad.TriggerMap[component.Shield_actionid] = false
    })
    ad.TriggerMap[component.Shield_actionid] = true //start out invulnerable

    // Move Left
    moveSpeed := 1.0
    ad.AddNormalAction(component.MoveLeft_actionid, func(){
        vd.Velocity.X = -1.0 * moveSpeed
    })

    // Move Right
    ad.AddNormalAction(component.MoveRight_actionid, func(){
        vd.Velocity.X = moveSpeed
    })

    ad.AddNormalAction(component.Destroy_actionid, func(){
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
        asset.PlaySound("PlayerShipDestroyed")

        event.ShipDestroyedEvent.Publish(ecs.World, event.ShipDestroyed{})

        AddExplosion(ecs, pd.Point.X, pd.Point.Y, "AstralianShip", view)
    })

    blinkCounter := 0
    ad.AddUpkeepAction(func(){
        // do a blinking effect if we are shielded
        if ad.TriggerMap[component.Shield_actionid] {
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
        if ad.TriggerMap[component.MoveRight_actionid] == ad.TriggerMap[component.MoveLeft_actionid] {
            vd.Velocity.X = 0
        }
    })

    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}
