package astra_entity

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

var frameStart int //used for syncing up animations

func AddAlienFormation(ecs *ecs.ECS, x, y float64, view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.Actions,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // SHIPS

    //height and width plus padding...
    sw := 18.0 //TODO make this a const
    sh := 18.0 

    columns := 10
    rows := 5

    ships := []*donburi.Entry{}
    var curOffsetY float64 = 0
    for r := 0; r < rows; r++ {
        var curOffsetX float64 = sw * float64(1 - columns) / 2
        for c := 0; c < columns; c++ {
            ships = append(ships, ecs.World.Entry(*AddAlien(ecs, x + curOffsetX, y + curOffsetY, view)))
            curOffsetX += sw
        }
        curOffsetY += sh
    }

    // Creates an event to remove ships from formation when destroyed
    removeFromFormationFunc := func(w donburi.World, event event.RemoveFromFormation) {
        new_ships := []*donburi.Entry{}
        for i, _ := range ships {
            if ships[i] == event.Entry {
                ecs.World.Remove(ships[i].Entity())
            } else {
                new_ships = append(new_ships, ships[i])
            }
        }
        ships = new_ships
    }
    event.RemoveFromFormationEvent.Subscribe(ecs.World, removeFromFormationFunc)
    // This will clean up the above event when the scene ends
    event.RegisterCleanupFuncEvent.Publish(
        ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.RemoveFromFormationEvent.Unsubscribe(ecs.World, removeFromFormationFunc)
            },
        },
    )

    // Actions
    tm := make(map[component.ActionId]bool)
    tm[component.MoveRight_actionid] = true //Start out moving right
    tm[component.SendShip_actionid] = true 
    cdm := make(map[component.ActionId]component.Cooldown)
    cdm[component.SendShip_actionid] = component.Cooldown{Cur:200, Max:200}
    am := make(map[component.ActionId]func())

    speed := 0.10 //TODO make this a constant
    am[component.MoveRight_actionid] = func() {
        if len(ships) == 0 {
            //TODO dispatch game won event?
            return
        }
        _, right := findExtremeShipX(ships)
        tempSpeed := speed
        if right + speed >= view.Area.Max.X {
            tempSpeed = view.Area.Max.X - right
            tm[component.MoveRight_actionid] = false
            tm[component.MoveLeft_actionid] = true
        }
        for i, _ := range ships {
            pos := component.Position.Get(ships[i])
            acts := component.Actions.Get(ships[i])

            if !acts.TriggerMap[component.Charge_actionid] {
                pos.Point.X += tempSpeed
            }
        }
    }
    am[component.MoveLeft_actionid] = func() {
        if len(ships) == 0 {
            return
        }
        left, _ := findExtremeShipX(ships)
        tempSpeed := speed
        if left - speed <= view.Area.Min.X {
            tempSpeed = view.Area.Min.X + left
            tm[component.MoveLeft_actionid] = false
            tm[component.MoveRight_actionid] = true
        }
        for i, _ := range ships {
            pos := component.Position.Get(ships[i]) //.Point.X
            acts := component.Actions.Get(ships[i])

            if !acts.TriggerMap[component.Charge_actionid] {
                pos.Point.X -= tempSpeed
            }
        }
    }
    am[component.SendShip_actionid] = func() {
        if len(ships) == 0 {
            //TODO dispatch game won event?
            return
        }
        r := rand.Intn(len(ships))

        acts := component.Actions.Get(ships[r])
        acts.TriggerMap[component.Charge_actionid] = true
        tm[component.SendShip_actionid] = false
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}

func findExtremeShipX(ships []*donburi.Entry) (float64, float64) {
    left := component.Position.Get(ships[0]).Point.X
    right := left
    for i, _ := range ships {
        pos := component.Position.Get(ships[i]).Point
        if pos.X < left {
            left = pos.X
        }
        if pos.X > right {
            right = pos.X
        }
    }
    left -= 8 //TODO need some good method to get sprite size/edge
    right += 8
    return left, right
}

func AddAlien(ecs *ecs.ECS, x, y float64, view *utility.View) *donburi.Entity {
    entity := ecs.Create(
        layer.Foreground, // TODO argument?
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Velocity,
        component.Collider,
        component.Health,
        component.Factions,
        component.Actions,
        component.Damage,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    // Factions
    factions := []component.FactionId{component.Enemy_factionid}
    donburi.SetValue(entry, component.Factions, factions)

    // Damage
    donburi.SetValue(entry, component.Damage, component.DamageData{
        Value: 1.0,
    })

    // Health
    donburi.SetValue(entry, component.Health, component.HealthData{Value:1.0})

    // Collider
    collider := component.NewColliderData()
    collider.Hitboxes = append(collider.Hitboxes, component.NewHitbox(4, 0, 0))
    donburi.SetValue(entry, component.Collider, collider)

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    spriteName := "Alien1"
    nsd.Load(spriteName, nil)
    nsd.Play("Idle")
    //nsd.SetFrame(rand.Intn(10)) //TODO should be constant
    nsd.SetFrame(frameStart)
    //frameStart = (frameStart + 1) % 10
    frameStart = (frameStart + rand.Intn(3)) % 10 //TODO which effect is best?
    gobj.Renderables = append(gobj.Renderables, &nsd)
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // Velocity
    vd := component.VelocityData{Velocity: &math.Vec2{X:0, Y:0}}
    donburi.SetValue(entry, component.Velocity, vd)

    // Action
    tm := make(map[component.ActionId]bool)
    cdm := make(map[component.ActionId]component.Cooldown)
    am := make(map[component.ActionId]func())
    am[component.Destroy_actionid] = func() {
        *nsd.RenderableData.GetTransInfo().Hide = true
        event.ScoreEvent.Publish(ecs.World, event.Score{Value:10})
        AddAlienExplosion(ecs, pd.Point.X, pd.Point.Y, spriteName, view)
        event.RemoveFromFormationEvent.Publish(
            ecs.World, event.RemoveFromFormation{Entry:entry})
        event.RemoveEntityEvent.Publish(
            ecs.World, 
            event.RemoveEntity{Entity:&entity},
        )
        //event.AsteroidsCountUpdateEvent.Publish(ecs.World, event.AsteroidsCountUpdate{Value:-1})
    }
    am[component.Charge_actionid] = func() {
        //println("charge action")
        vd.Velocity.X = 0
        vd.Velocity.Y = 0
        // turn towards enemy
        //set velocity?
        //TODO continue here
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}
