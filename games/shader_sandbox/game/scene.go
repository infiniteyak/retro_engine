package game

import (
	//"github.com/infiniteyak/retro_engine/engine/entity"
	//"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/infiniteyak/retro_engine/engine/scene"
)

const (
    Undefined_sceneEvent scene.SceneEventId = iota
    Init_sceneEvent
)

const (
    Undefined_sceneId scene.SceneId = iota
    Main_sceneId
)

func (this *Game) InitStates() {
    this.states = map[scene.SceneId]map[scene.SceneEventId]func() {
        Undefined_sceneId: { 
            Init_sceneEvent: this.LoadMainScene,
        },
    }
}

func (this *Game) Transition(event scene.SceneEventId) {
    if this.states[this.curScene.GetId()][event] != nil {
        this.curScene.Cleanup()
        sid := this.curScene.GetId()
        this.curScene = scene.NewScene(this.ecs)
        this.states[sid][event]()
    } else {
        println("states map miss")
    }
}

