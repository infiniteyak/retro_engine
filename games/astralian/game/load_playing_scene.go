package game

import (
	"fmt"
	"github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/infiniteyak/retro_engine/games/astralian/entity"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/infiniteyak/retro_engine/engine/component"
	"strings"
	"github.com/yohamta/donburi"
	"github.com/hajimehoshi/ebiten/v2"
)

func (this *Game) LoadPlayingScene() {
    println("LoadPlayingScene")
    this.curScene.SetId(Playing_sceneId)

    asset.PlayMusic("Music")
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
        strings.Repeat("^", this.curShips),
    )
    shipsText.XAlign = entity.Right_fontalignx
    shipsText.YAlign = entity.Top_fontaligny

    // GAME
    gameView := utility.NewView(
        0.0, 
        hudView.Area.Max.Y,
        this.screenView.Area.Max.X, 
        this.screenView.Area.Max.Y - hudView.Area.Max.Y,
    )

    var playerPos *component.PositionData
    shipDestFunc := func(w donburi.World, event event.ShipDestroyed) {
        this.curShips--
        if this.curShips < 0 {
            println("game over")
            entity.AddTitleText(
                this.ecs, 
                float64(gameView.Area.Max.X / 2), 
                float64(gameView.Area.Max.Y / 2), 
                gameView,
                "GAME OVER",
            )
            
            entity.AddInputTrigger(
                this.ecs, 
                ebiten.KeySpace,
                func() {
                    this.Transition(GameOver_sceneEvent)
                },
            )
        } else {
            shipsText.String = strings.Repeat("^", this.curShips) 
            psEntity := astra_entity.AddPlayerShip(
                this.ecs, 
                float64(gameView.Area.Max.X / 2), 
                float64(gameView.Area.Max.Y - 10), 
                gameView,
            )
            *playerPos = *component.Position.Get(this.ecs.World.Entry(*psEntity))
        }
    }
    event.ShipDestroyedEvent.Subscribe(this.ecs.World, shipDestFunc)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.ShipDestroyedEvent.Unsubscribe(this.ecs.World, shipDestFunc)
            },
        },
    )

    screenClearFunc := func(w donburi.World, event event.ScreenClear) {
        println("Screen Clear")
        this.curWave++
        entity.AddTitleText(
            this.ecs, 
            float64(gameView.Area.Max.X / 2), 
            float64(gameView.Area.Max.Y / 2), 
            gameView,
            fmt.Sprintf("PREPARE FOR WAVE %03d", this.curWave),
        )
        
        entity.AddInputTrigger(
            this.ecs, 
            ebiten.KeySpace,
            func() {
                this.Transition(ScreenClear_sceneEvent)
            },
        )
    }
    event.ScreenClearEvent.Subscribe(this.ecs.World, screenClearFunc)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.ScreenClearEvent.Unsubscribe(this.ecs.World, screenClearFunc)
            },
        },
    )

    this.GenerateStars(gameView)

    psEntity := astra_entity.AddPlayerShip(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Max.Y - 10), 
        gameView,
    )
    playerPos = component.Position.Get(this.ecs.World.Entry(*psEntity))

    astra_entity.AddAlienFormation(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Min.Y + 40), 
        gameView,
        playerPos,
        this.curWave,
    )

    asset.PlaySound("Wave")
}
