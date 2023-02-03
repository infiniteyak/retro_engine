package game

import (
	"math/rand"
	//"github.com/infiniteyak/retro_engine/games/astralian/entity"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/infiniteyak/retro_engine/engine/scene"
)

const (
    Undefined_sceneEvent scene.SceneEventId = iota
    Init_sceneEvent
    Advance_sceneEvent
    GameStart_sceneEvent
    GameOver_sceneEvent
    ScreenClear_sceneEvent
)

const (
    Undefined_sceneId scene.SceneId = iota
    Title_sceneId
    Info_sceneId
    ScoreBoard_sceneId
    Playing_sceneId
    EnterInitials_sceneId
)

func (this *Game) InitStates() {
    this.states = map[scene.SceneId]map[scene.SceneEventId]func() {
        Undefined_sceneId: { 
            Init_sceneEvent: this.LoadTitleScene,
        },
        Title_sceneId: { 
            Advance_sceneEvent: this.LoadInfoScene,
            GameStart_sceneEvent: this.LoadPlayingScene,
        },
        Info_sceneId: { 
            Advance_sceneEvent: this.LoadScoreBoardScene,
            GameStart_sceneEvent: this.LoadPlayingScene,
        },
        ScoreBoard_sceneId: {
            Advance_sceneEvent: this.LoadTitleScene,
            GameStart_sceneEvent: this.LoadPlayingScene,
        },
        Playing_sceneId: {
            GameOver_sceneEvent: this.LoadInitialsScene,
            ScreenClear_sceneEvent: this.LoadPlayingScene,
        },
        EnterInitials_sceneId: {
            Advance_sceneEvent: this.LoadScoreBoardScene,
        },
    }
}

//TODO move elsewhere
func (this *Game) GenerateStars(view *utility.View) {
    println("generating stars")
    w := int(view.Area.Max.X)
    h := int(view.Area.Max.Y)

    for i := 0; i < 50; i++ {
        x := rand.Intn(w)
        y := rand.Intn(h)
        entity.AddStar(this.ecs, float64(x), float64(y), view)
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

