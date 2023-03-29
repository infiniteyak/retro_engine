package game

import (
	"github.com/infiniteyak/retro_engine/engine/scene"
)

const (
    Undefined_sceneEvent scene.SceneEventId = iota
    Init_sceneEvent
    Advance_sceneEvent
    GameStart_sceneEvent
    GameOver_sceneEvent
    ScreenClear_sceneEvent
    GoOptions_sceneEvent
    GoScores_sceneEvent
    GoAttract_sceneEvent
    GoInfo_sceneEvent
)

const (
    Undefined_sceneId scene.SceneId = iota
    Attract_sceneId 
    Menu_sceneId
    Info_sceneId
    ScoreBoard_sceneId
    Playing_sceneId
    EnterInitials_sceneId
    Options_sceneId
)

func (this *Game) InitStates() {
    this.states = map[scene.SceneId]map[scene.SceneEventId]func() {
        Undefined_sceneId: { 
            Init_sceneEvent: this.LoadAttractModeScene,
        },
        Attract_sceneId: { 
            Advance_sceneEvent: this.LoadMenuScene,
        },
        Menu_sceneId: { 
            GameStart_sceneEvent: this.LoadPlayingScene,
            GoOptions_sceneEvent: this.LoadOptionsScene,
            GoScores_sceneEvent: this.LoadScoreBoardScene,
            GoAttract_sceneEvent: this.LoadAttractModeScene,
            GoInfo_sceneEvent: this.LoadInfoScene,
        },
        Info_sceneId: { 
            Advance_sceneEvent: this.LoadMenuScene,
        },
        Options_sceneId: { 
            Advance_sceneEvent: this.LoadMenuScene,
        },
        ScoreBoard_sceneId: { 
            Advance_sceneEvent: this.LoadMenuScene,
        },
        Playing_sceneId: {
            GameOver_sceneEvent: this.LoadInitialsScene,
            ScreenClear_sceneEvent: this.LoadPlayingScene, //TODO add support for intra-level cutscenes
        },
        EnterInitials_sceneId: {
            Advance_sceneEvent: this.LoadScoreBoardScene,
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

