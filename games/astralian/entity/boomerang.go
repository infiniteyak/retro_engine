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
    gMath "math"
)

func AddBoomerang( ecs *ecs.ECS, 
                   pX, pY float64, 
                   velocity math.Vec2, 
                   view *utility.View, 
                   power int, 
                   parent *donburi.Entity) *donburi.Entity {
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

    playerShipEntry := ecs.World.Entry(*parent)

    // Factions
    factions := []component.FactionId{component.Player_factionid} //TODO should be arg
    donburi.SetValue(entry, component.Factions, factions)

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(3, 0, 0))
    donburi.SetValue(entry, component.Collider, collider)

    // Position
    pd := component.NewPositionData(pX, pY)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load("Boomerang", nil)
    nsd.Play("")
    nsd.SetPlaySpeed(2.0) //TODO should be constant
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    vd := component.VelocityData{Velocity: &velocity}
    donburi.SetValue(entry, component.Velocity, vd)

    // Actions
    ad := component.NewActions()

    ad.AddCooldownAction(component.SelfDestruct_actionid, 400, func() {
        ad.TriggerMap[component.SelfDestruct_actionid] = false
        ad.TriggerMap[component.DestroySilent_actionid] = true
    })
    ad.TriggerMap[component.SelfDestruct_actionid] = true

    caught := false
    hits := 0
    ad.AddNormalAction(component.DestroySilent_actionid, func() {
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
        if playerShipEntry.Valid() {
            playerActions := component.Actions.Get(playerShipEntry)
            playerActions.TriggerMap[component.Reload_actionid] = true
            if caught {
                playerActions.TriggerMap[component.IncreasePower_actionid] = true
                event.ScoreEvent.Publish(ecs.World, event.Score{Value:hits*hits})
            } else {
                playerActions.TriggerMap[component.ResetPower_actionid] = true
            }
        }
        caught = false
    })

    ad.AddNormalAction(component.Destroy_actionid, func() {
        asset.PlaySound("GenericHit")
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
    })

    playerPos := component.Position.Get(playerShipEntry)
    //TODO might be cool to have a random offset in x for sourcePos
    sourcePos := component.NewPositionData(playerPos.Point.X, playerPos.Point.Y)
    curAngle := gMath.Pi
    ad.AddNormalAction(component.ReturnProjectile_actionid, func() {
        aimOffsetY := 0.0 

        angleRad := 0.0
        if pd.Point.Y < (sourcePos.Point.Y + aimOffsetY) {
            // angle towards player ship
            angleRad = gMath.Atan2(pd.Point.X - sourcePos.Point.X, (sourcePos.Point.Y + aimOffsetY) - pd.Point.Y)

            // Turn towards the point we're aiming at
            a := curAngle - angleRad
            a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
            if a <= 0 {
                curAngle += 0.05
            } else {
                curAngle -= 0.05
            }
        } 

        // Use move rotation and charge speed to create a vector for movement
        moveVect := math.Vec2{X:0, Y:1.5}
        moveVect = moveVect.Rotate(curAngle)

        pd.Point.X += moveVect.X
        pd.Point.Y += moveVect.Y

        // clean up the current angle
        if curAngle >= (2 * gMath.Pi) {
            curAngle -= (2 * gMath.Pi)
        }

        if playerShipEntry.Valid() {
            playerCollider := component.Collider.Get(playerShipEntry)
            for _, cEntry := range playerCollider.Collisions {
                if cEntry.Entity() == entity {
                    ad.TriggerMap[component.DestroySilent_actionid] = true
                    caught = true
                }
            }
        }
    })

    donburi.SetValue(entry, component.Actions, ad)

    // Damage
    dd := component.NewDamageData()
    *dd.Value = 1.0
    *dd.DestroyOnDamage = false
    dd.OnDamage = func() {
        ad.TriggerMap[component.ReturnProjectile_actionid] = true
        vd.Velocity.X = 0
        vd.Velocity.Y = 0

        asset.PlaySound("GenericHit")

        hits++
        power--
        if power <= 0 {
            *dd.Value = 0
        }
    }
    donburi.SetValue(entry, component.Damage, dd)

    asset.PlaySound("SciFiProjectile")

    return &entity
}
