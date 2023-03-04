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

type alienFormationData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    audioContext *audio.Context

    actions component.ActionsData
    position component.PositionData
    view component.ViewData

    playerPos *component.PositionData
    ships []*donburi.Entry
    shipSlots map[*donburi.Entry]component.PositionData
    sendCd int
    wave int
}

func findExtremeShipX(ships []*donburi.Entry) (float64, float64) {
    left := component.Position.Get(ships[0]).Point.X
    right := left
    for i, _ := range ships {
        acts := component.Actions.Get(ships[i])
        if  acts.TriggerMap[component.Charge_actionid] ||
            acts.TriggerMap[component.ReturnShip_actionid] ||
            acts.TriggerMap[component.Follow_actionid] {
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

func (this *alienFormationData) initConvoy(x, y float64) {
    sw := AlienConvoySpacing
    sh := AlienConvoySpacing

    columns := 10
    rows := 6
    pattern := [][]int {
        {0, 0, 0, 4, 0, 0, 4, 0, 0, 0,},
        {0, 0, 3, 3, 3, 3, 3, 3, 0, 0,},
        {0, 2, 2, 2, 2, 2, 2, 2, 2, 0,},
        {1, 1, 1, 1, 1, 1, 1, 1, 1, 1,},
        {0, 1, 1, 1, 1, 1, 1, 1, 1, 0,},
        {0, 0, 1, 1, 1, 1, 1, 1, 0, 0,},
    }

    this.ships = []*donburi.Entry{}
    this.shipSlots = map[*donburi.Entry]component.PositionData{}
    var curOffsetY float64 = 0
    bosses := []*donburi.Entry{}
    defenderCount := 0
    bossIndex := 0
    for r := 0; r < rows; r++ {
        var curOffsetX float64 = sw * float64(1 - columns) / 2
        for c := 0; c < columns; c++ {
            aType := Undefined_alientype
            if pattern[r][c] == 1 {
                aType = Blue_alientype
            } else if pattern[r][c] == 2 {
                aType = Purple_alientype
            } else if pattern[r][c] == 3 {
                aType = Green_alientype
            } else if pattern[r][c] == 4 {
                aType = Grey_alientype
            }
            if aType != Undefined_alientype { //TODO this is kind of hacky
                var curBoss *donburi.Entry
                if aType == Green_alientype {
                    if defenderCount >= 3 {
                        bossIndex++
                        defenderCount = 0
                    } else {
                        defenderCount++
                    }
                    if len(bosses) > bossIndex {
                        curBoss = bosses[bossIndex]
                    }
                }
                ship_entry := this.ecs.World.Entry(*AddAlien(this.ecs, x + curOffsetX, y + curOffsetY, this.view.View, this.audioContext, this.playerPos, aType, curBoss, this.wave))
                this.ships = append(this.ships, ship_entry)
                if aType == Grey_alientype {
                    bosses = append(bosses, ship_entry)
                }
                this.shipSlots[ship_entry] = component.NewPositionData(x + curOffsetX, y + curOffsetY)
            }
            curOffsetX += sw
        }
        curOffsetY += sh
    }

    // Creates an event to remove ships from formation when destroyed
    removeFromFormationFunc := func(w donburi.World, event event.RemoveFromFormation) {
        new_ships := []*donburi.Entry{}
        for i, _ := range this.ships {
            if this.ships[i] == event.Entry {
                this.ecs.World.Remove(this.ships[i].Entity())
            } else {
                new_ships = append(new_ships, this.ships[i])
            }
        }
        this.ships = new_ships
    }
    event.RemoveFromFormationEvent.Subscribe(this.ecs.World, removeFromFormationFunc)
    // This will clean up the above event when the scene ends
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.RemoveFromFormationEvent.Unsubscribe(this.ecs.World, removeFromFormationFunc)
            },
        },
    )
}

func AddAlienFormation(ecs *ecs.ECS, 
                       x, y float64, 
                       view *utility.View, 
                       playerPos *component.PositionData, 
                       audioContext *audio.Context,
                       wave int) *donburi.Entity {
    afd := &alienFormationData{}
    afd.ecs = ecs
    afd.wave = wave

    afd.playerPos = playerPos //so the AI can track the player ship

    entity := afd.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.Actions,
    )
    afd.entity = &entity

    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:afd.entity})

    afd.entry = ecs.World.Entry(*afd.entity)

    afd.audioContext = audioContext

    // Position
    afd.position = component.NewPositionData(x, y)
    donburi.SetValue(afd.entry, component.Position, afd.position)

    // View
    afd.view = component.ViewData{View:view}
    donburi.SetValue(afd.entry, component.View, afd.view)

    // SHIPS
    afd.initConvoy(x,y)

    // Actions
    afd.actions = component.NewActions()

    afd.sendCd = AlienConvoySendCd - (30 * (afd.wave-1))
    print(afd.sendCd)
    if afd.sendCd < AlienConvoySendInitCd {
        afd.sendCd = AlienConvoySendInitCd 
    }

    afd.actions.TriggerMap[component.MoveRight_actionid] = true //Start out moving right
    afd.actions.TriggerMap[component.SendShip_actionid] = true //Dispatch ships
    afd.actions.CooldownMap[component.SendShip_actionid] = component.Cooldown{
        Cur:AlienConvoySendInitCd, 
        Max:afd.sendCd,
    }

    cleared := false
    afd.actions.ActionMap[component.Upkeep_actionid] = func() {
        if len(afd.ships) == 0 && !cleared {
            cleared = true
            se := event.Score{Value:400 * afd.wave}
            event.ScoreEvent.Publish(ecs.World, se)
            event.ScreenClearEvent.Publish(ecs.World, event.ScreenClear{})
            return
        }
    }

    afd.actions.ActionMap[component.MoveRight_actionid] = func() {
        if len(afd.ships) == 0 {
            return
        }
        _, right := findExtremeShipX(afd.ships)
        tempSpeed := AlienConvoySpeed 
        if right + AlienConvoySpeed >= afd.view.View.Area.Max.X {
            tempSpeed = afd.view.View.Area.Max.X - right
            afd.actions.TriggerMap[component.MoveRight_actionid] = false
            afd.actions.TriggerMap[component.MoveLeft_actionid] = true
        }
        for i, _ := range afd.ships {
            pos := component.Position.Get(afd.ships[i])
            acts := component.Actions.Get(afd.ships[i])

            if tempSpeed < 0 {
                tempSpeed = 0
            }
            afd.shipSlots[afd.ships[i]].Point.X += tempSpeed
            
            if  !acts.TriggerMap[component.Charge_actionid] &&
                !acts.TriggerMap[component.Follow_actionid] {
                pos.Point.X = afd.shipSlots[afd.ships[i]].Point.X 
            }
        }
    }
    afd.actions.ActionMap[component.MoveLeft_actionid] = func() {
        if len(afd.ships) == 0 {
            return
        }
        left, _ := findExtremeShipX(afd.ships)
        tempSpeed := AlienConvoySpeed 
        if left - AlienConvoySpeed <= afd.view.View.Area.Min.X {
            tempSpeed = afd.view.View.Area.Min.X + left
            afd.actions.TriggerMap[component.MoveLeft_actionid] = false
            afd.actions.TriggerMap[component.MoveRight_actionid] = true
        }
        for i, _ := range afd.ships {
            pos := component.Position.Get(afd.ships[i])
            acts := component.Actions.Get(afd.ships[i])

            if tempSpeed < 0 {
                tempSpeed = 0
            }
            afd.shipSlots[afd.ships[i]].Point.X -= tempSpeed

            if  !acts.TriggerMap[component.Charge_actionid] &&
                !acts.TriggerMap[component.Follow_actionid] {
                pos.Point.X = afd.shipSlots[afd.ships[i]].Point.X 
            }
        }
    }
    afd.actions.ActionMap[component.SendShip_actionid] = func() {
        if len(afd.ships) == 0 {
            return
        }
        r := rand.Intn(len(afd.ships))

        acts := component.Actions.Get(afd.ships[r])

        if acts.TriggerMap[component.Charge_actionid] ||
           acts.TriggerMap[component.Follow_actionid] ||
           acts.TriggerMap[component.ReturnShip_actionid] {
            afd.actions.CooldownMap[component.SendShip_actionid] = component.Cooldown{
                Cur: AlienConvoySendInitCd, 
                Max: AlienConvoySendCd,
            } 
            return
        }

        acts.TriggerMap[component.Charge_actionid] = true
        acts.CooldownMap[component.Shoot_actionid] = component.Cooldown{
            Cur:AlienShootDelay, 
            Max:AlienShootDelay,
        }
 
        afd.actions.CooldownMap[component.SendShip_actionid] = component.Cooldown{
            Cur: afd.sendCd, 
            Max: afd.sendCd,
        } 
    }

    donburi.SetValue(afd.entry, component.Actions, afd.actions)

    return &entity
}
