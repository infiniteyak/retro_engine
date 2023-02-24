package astra_entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "math/rand"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

func findExtremeShipX(ships []*donburi.Entry) (float64, float64) {
    left := component.Position.Get(ships[0]).Point.X
    right := left
    for i, _ := range ships {
        acts := component.Actions.Get(ships[i])
        if acts.TriggerMap[component.Charge_actionid] {
            continue
        }
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

func AddAlienFormation(ecs *ecs.ECS, x, y float64, view *utility.View, playerPos *component.PositionData, audioContext *audio.Context) *donburi.Entity {
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
    sw := 14.0 //TODO make this a const
    sh := 14.0 

    columns := 10
    rows := 6

    ships := []*donburi.Entry{}
    shipSlots := map[*donburi.Entry]component.PositionData{}
    var curOffsetY float64 = 0
    for r := 0; r < rows; r++ {
        var curOffsetX float64 = sw * float64(1 - columns) / 2
        for c := 0; c < columns; c++ {
            aType := Undefined_alientype
            /*
            if r % 2 == 0 {
                aType = Blue_alientype
            } else {
                aType = Purple_alientype
            }
            */
            aType = Purple_alientype
            ship_entry := ecs.World.Entry(*AddAlien(ecs, x + curOffsetX, y + curOffsetY, view, audioContext, playerPos, aType))
            ships = append(ships, ship_entry)
            shipSlots[ship_entry] = component.NewPositionData(x + curOffsetX, y + curOffsetY)
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

    speed := AlienConvoySpeed
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

            if tempSpeed < 0 {
                tempSpeed = 0
            }
            shipSlots[ships[i]].Point.X += tempSpeed
            
            if !acts.TriggerMap[component.Charge_actionid] {
                pos.Point.X = shipSlots[ships[i]].Point.X 
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

            if tempSpeed < 0 {
                tempSpeed = 0
            }
            shipSlots[ships[i]].Point.X -= tempSpeed

            if !acts.TriggerMap[component.Charge_actionid] {
                pos.Point.X = shipSlots[ships[i]].Point.X 
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
        acts.CooldownMap[component.Shoot_actionid] = component.Cooldown{
            Cur:AlienShootDelay, 
            Max:AlienShootDelay,
        }
 
        //tm[component.SendShip_actionid] = false
        cdm[component.SendShip_actionid] = component.Cooldown{Cur:600, Max:600} 
    }

    donburi.SetValue(entry, component.Actions, component.ActionsData{
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    })

    return &entity
}
