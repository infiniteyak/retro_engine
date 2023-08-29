package shape_courier_entity

import (
	"github.com/infiniteyak/retro_engine/engine/component"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const (
    scatterTime = 2000
    chaseTime = 2000
)

var scatterTimes = []([]int){
    {840, 840, 600},
    {840, 840, 600},
    {840, 840, 600},
    {840, 840, 600},
    {600, 600, 600},
}
var chaseTimes = []([]int){
    {2400, 2400, 2400},
    {2400, 2400, 120000},
    {2400, 2400, 120000},
    {2400, 2400, 120000},
    {2400, 2400, 120000},
}

type ghostControllerData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    actions component.ActionsData
    curTime int
    curMode AiMode
    ghostSpawnTimerMap map[GhostVarient]int
}

func AddGhostController(ecs *ecs.ECS, wave int, spawnGhost func(GhostVarient) *GhostData) {
    this := &ghostControllerData{}
    this.ecs = ecs

    entity := this.ecs.World.Create(
        component.Actions,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)

    this.curMode = Scatter_aimode //Start in scatter


    despawn := func(w donburi.World, event event.DespawnAllEnemies) {
        this.ghostSpawnTimerMap = make(map[GhostVarient]int)
    }
    event.DespawnAllEnemiesEvent.Subscribe(this.ecs.World, despawn)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.DespawnAllEnemiesEvent.Unsubscribe(this.ecs.World, despawn)
            },
        },
    )
    round := 0

    respawn := func(w donburi.World, event event.RespawnEnemies) {
        this.curTime = 0
        this.curMode = Scatter_aimode //Start in scatter
        round = 0
        this.ghostSpawnTimerMap = make(map[GhostVarient]int)
        this.ghostSpawnTimerMap[ClassicRed_ghostvarient] = ghostColorDelayMap[ClassicRed_ghostvarient]
        this.ghostSpawnTimerMap[ClassicPink_ghostvarient] = ghostColorDelayMap[ClassicPink_ghostvarient]
        this.ghostSpawnTimerMap[ClassicBlue_ghostvarient] = ghostColorDelayMap[ClassicBlue_ghostvarient]
        this.ghostSpawnTimerMap[ClassicOrange_ghostvarient] = ghostColorDelayMap[ClassicOrange_ghostvarient]
    }
    event.RespawnEnemiesEvent.Subscribe(this.ecs.World, respawn)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.RespawnEnemiesEvent.Unsubscribe(this.ecs.World, respawn)
            },
        },
    )

    if wave > len(scatterTimes) {
        wave = len(scatterTimes) - 1
    } else {
        wave--
    }

    // Actions
    this.actions = component.NewActions()
    this.actions.AddUpkeepAction(func(){
        this.curTime++

        for k, v := range this.ghostSpawnTimerMap {
            if v >= 0 {
                this.ghostSpawnTimerMap[k]-- //TODO does that work?
            }
            if v == 0 {
                spawnGhost(k)
            }
        }

        if this.curMode == Scatter_aimode && this.curTime >= scatterTimes[wave][round] {
            this.curTime = 0
            this.curMode = Chase_aimode
            event.SetAiModeEvent.Publish(this.ecs.World, event.AiMode{Value:int(Chase_aimode)})
        } else if this.curMode == Chase_aimode && this.curTime >= chaseTimes[wave][round] {
            this.curTime = 0
            this.curMode = Scatter_aimode
            event.SetAiModeEvent.Publish(this.ecs.World, event.AiMode{Value:int(Scatter_aimode)})
            round = (round + 1) % len(scatterTimes[wave])
        }
    })

    donburi.SetValue(this.entry, component.Actions, this.actions)
}
