package game

import (
	"fmt"
	"github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/infiniteyak/retro_engine/games/shape_courier/entity"
    "github.com/infiniteyak/retro_engine/engine/layer"
)

var elroyModeThreshold = []int{ // and half this for elroy 2 mode
    20,
    30,
    40,
    40,
    40,
    50,
    50,
    50,
    60,
    60,
    60,
    80,
    80,
    80,
    100,
    100,
    100,
    100,
    120,
    120,
    120,
}

func (this *Game) LoadPlayingScene() {
    println("LoadPlayingScene")
    this.curScene.SetId(Playing_sceneId)

    //asset.PlayMusic("Music")

    // HUD
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

    gameView := utility.NewView(
        0.0, 
        hudView.Area.Max.Y,
        this.screenView.Area.Max.X, 
        this.screenView.Area.Max.Y - hudView.Area.Max.Y,
    )

    spawnPlayer := func() {panic("Tried to spawn player too early")}
    livesObjects := make([]*donburi.Entity, this.curLives)
    adjustLives := func(w donburi.World, e event.AdjustLives) {
        this.curLives += e.Value
        if this.curLives < 0 {
            gameOverText := entity.AddTitleText(
                this.ecs, 
                float64(gameView.Area.Max.X / 2), 
                float64(gameView.Area.Max.Y / 2) - 13.0, 
                gameView,
                "GAME OVER", //TODO adjust font
            )
            gameOverText.XAlign = entity.Center_fontalignx
            gameOverText.YAlign = entity.Middle_fontaligny
            gameOverText.Blink = true
            entity.AddTimer(this.ecs, 500, func(){
                gameOverText.Blink = false
                ree := event.RemoveEntity{Entity:gameOverText.Entity}
                event.RemoveEntityEvent.Publish(this.ecs.World, ree)
                this.Transition(GameOver_sceneEvent)
            })
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
        spawnPlayer()
    }
    event.AdjustLivesEvent.Subscribe(this.ecs.World, adjustLives)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.AdjustLivesEvent.Unsubscribe(this.ecs.World, adjustLives)
            },
        },
    )
    
    dots := 0
    adjustDots := func(w donburi.World, e event.AdjustDots) {
        dots += e.Value
        if e.Value > 0 {
            return
        }
        if dots == 0 {
            this.curWave++
            this.Transition(ScreenClear_sceneEvent)
        } else if this.curWave >= len(elroyModeThreshold) && dots <= elroyModeThreshold[len(elroyModeThreshold)-1] {
            // TODO test this
            event.ElroyModeEvent.Publish(this.ecs.World, event.ElroyMode{})
        } else if this.curWave < len(elroyModeThreshold) && dots <= elroyModeThreshold[this.curWave-1] {
            event.ElroyModeEvent.Publish(this.ecs.World, event.ElroyMode{})
        }
    }
    event.AdjustDotsEvent.Subscribe(this.ecs.World, adjustDots)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.AdjustDotsEvent.Unsubscribe(this.ecs.World, adjustDots)
            },
        },
    )

    mazeData := shape_courier_entity.AddMaze(
        this.ecs, 
        float64(gameView.Area.Max.X / 2), 
        float64(gameView.Area.Max.Y / 2), 
        gameView)

    var mandyData *shape_courier_entity.MandyData = nil
    var redGhostData *shape_courier_entity.GhostData = nil
    spawnGhost := func(varient shape_courier_entity.GhostVarient) *shape_courier_entity.GhostData {
        if redGhostData == nil && varient == shape_courier_entity.ClassicBlue_ghostvarient {
            panic("Must initialize red ghost before blue ghost!")
        }
        gd := shape_courier_entity.AddGhost(
            this.ecs, 
            gameView,
            mandyData,
            mazeData,
            varient,
            redGhostData,
            this.curWave)
        if varient == shape_courier_entity.ClassicRed_ghostvarient {
            redGhostData = gd
        }
        return gd
    }
    shape_courier_entity.AddGhostController(this.ecs, this.curWave, spawnGhost)

    shape_courier_entity.AddTaco(
        this.ecs, 
        gameView,
        mazeData)

    spawnPlayer = func() {
        readyText := entity.AddTitleText(
            this.ecs, 
            float64(gameView.Area.Max.X / 2), 
            float64(gameView.Area.Max.Y / 2) - 13.0, 
            gameView,
            "READY", //TODO adjust font, add !
        )
        readyText.XAlign = entity.Center_fontalignx
        readyText.YAlign = entity.Middle_fontaligny
        readyText.Blink = true

        entity.AddTimer(this.ecs, 200, func(){
            readyText.Blink = false
            ree := event.RemoveEntity{Entity:readyText.Entity}
            event.RemoveEntityEvent.Publish(this.ecs.World, ree)
            mandyData = shape_courier_entity.AddSpaceMandy(
                this.ecs, 
                gameView,
                mazeData)
            respEnemies := event.RespawnEnemies{}
            event.RespawnEnemiesEvent.Publish(this.ecs.World, respEnemies)
        })
    }
    event.AdjustLivesEvent.Publish(
        this.ecs.World, 
        event.AdjustLives{
            Value: 0,
        },
    )

    //play the start noise
    asset.PlaySound("StartNoise")
}
