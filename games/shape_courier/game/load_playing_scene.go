package game

import (
	"fmt"
	"github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"strings"
	"github.com/yohamta/donburi"
	"github.com/infiniteyak/retro_engine/games/shape_courier/entity"
)

func (this *Game) LoadPlayingScene() {
    println("LoadPlayingScene")
    this.curScene.SetId(Playing_sceneId)

    //asset.PlayMusic("Music")

    // HUD
    hudView := utility.NewView(0.0, 0.0, this.screenView.Area.Max.X, asset.FontHeight)

    entity.AddBlackBar(
        this.ecs, 
        float64(hudView.Area.Max.X / 2),
        float64(hudView.Area.Max.Y / 2),
        hudView,
    )

    // Create score text
    curScoreText := entity.AddNormalText(
        this.ecs, 
        0,
        0,
        hudView,
        "WhiteFont",
        fmt.Sprintf("%06d", this.curScore.Value),
    )
    curScoreText.XAlign = entity.Left_fontalignx
    curScoreText.YAlign = entity.Top_fontaligny

    // Score text update code
    scoreFunction := func(w donburi.World, event event.Score) {
        this.curScore.Value += event.Value
        maxScore := 999999
        if this.curScore.Value > maxScore {
            this.curScore.Value = maxScore 
        }
        curScoreText.String = fmt.Sprintf("%06d", this.curScore.Value)
    }
    event.ScoreEvent.Subscribe(this.ecs.World, scoreFunction)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.ScoreEvent.Unsubscribe(this.ecs.World, scoreFunction)
            },
        },
    )

    // Create wave text
    waveText := entity.AddNormalText(
        this.ecs, 
        float64(hudView.Area.Max.X / 2), 
        0,
        hudView,
        "WhiteFont",
        fmt.Sprintf("%03d", this.curWave),
    )
    waveText.YAlign = entity.Top_fontaligny

    // Create ships text (lives)
    shipsText := entity.AddNormalText(
        this.ecs, 
        float64(hudView.Area.Max.X), 
        0,
        hudView,
        "WhiteFont",
        strings.Repeat("^", this.curLives),
    )
    shipsText.XAlign = entity.Right_fontalignx
    shipsText.YAlign = entity.Top_fontaligny
    
    gameView := utility.NewView(
        0.0, 
        hudView.Area.Max.Y,
        this.screenView.Area.Max.X, 
        this.screenView.Area.Max.Y - hudView.Area.Max.Y,
    )

    md := shape_courier_entity.AddMaze(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Max.Y / 2), 
        gameView)

    shape_courier_entity.AddSpaceMandy(
        this.ecs, 
        //float64(gameView.Area.Max.X / 2), 
        //float64(gameView.Area.Max.Y / 2)-1, 
        gameView,
        md)
}
