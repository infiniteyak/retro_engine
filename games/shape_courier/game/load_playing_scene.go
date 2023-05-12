package game

import (
	"fmt"
	"github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/utility"
	//"strings"
	"github.com/yohamta/donburi"
	"github.com/infiniteyak/retro_engine/games/shape_courier/entity"
    "github.com/infiniteyak/retro_engine/engine/layer"
)

func (this *Game) LoadPlayingScene() {
    println("LoadPlayingScene")
    this.curScene.SetId(Playing_sceneId)

    //asset.PlayMusic("Music")

    // HUD
    //hudView := utility.NewView(0.0, 0.0, this.screenView.Area.Max.X, asset.FontHeight)
    hudView := utility.NewView(0.0, 1.0, this.screenView.Area.Max.X, asset.FontHeight + 1)

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

    /*
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
    */
    livesObjects := make([]*donburi.Entity, this.curLives)
    adjustLives := func(w donburi.World, e event.AdjustLives) {
        this.curLives += e.Value
        if this.curLives < 0 {
            //TODO handle this
            println("game over")
            return
        }
        for i := 0; i < len(livesObjects); i++ {
            ree := event.RemoveEntity{Entity:livesObjects[i]}
            event.RemoveEntityEvent.Publish(this.ecs.World, ree)
        }
        livesXVal := float64(hudView.Area.Max.X) - 4.0
        for i := 0; i < this.curLives; i++ {
            livesObjects = append(livesObjects, entity.AddSpriteObject(
                this.ecs,
                layer.HudForeground,
                livesXVal,
                float64(hudView.Area.Max.Y / 2), 
                "Life",
                "",
                hudView,
                ))
            livesXVal -= 10
        }
    }
    event.AdjustLivesEvent.Subscribe(this.ecs.World, adjustLives)
    event.AdjustLivesEvent.Publish(
        this.ecs.World, 
        event.AdjustLives{
            Value: 0,
        },
    )
    
    gameView := utility.NewView(
        0.0, 
        hudView.Area.Max.Y,
        this.screenView.Area.Max.X, 
        this.screenView.Area.Max.Y - hudView.Area.Max.Y,
    )

    mazeData := shape_courier_entity.AddMaze(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Max.Y / 2), 
        gameView)

    mandyData := shape_courier_entity.AddSpaceMandy(
        this.ecs, 
        gameView,
        mazeData)

    shape_courier_entity.AddGhost(
        this.ecs, 
        gameView,
        mandyData,
        mazeData,
        shape_courier_entity.ClassicRed_ghostvarient,
        nil)
    redGhost := shape_courier_entity.AddGhost(
        this.ecs, 
        gameView,
        mandyData,
        mazeData,
        shape_courier_entity.ClassicPink_ghostvarient,
        nil)
    shape_courier_entity.AddGhost(
        this.ecs, 
        gameView,
        mandyData,
        mazeData,
        shape_courier_entity.ClassicBlue_ghostvarient,
        redGhost)
    shape_courier_entity.AddGhost(
        this.ecs, 
        gameView,
        mandyData,
        mazeData,
        shape_courier_entity.ClassicOrange_ghostvarient,
        nil)
    shape_courier_entity.AddGhostController(this.ecs)
}
