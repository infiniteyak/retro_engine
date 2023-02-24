package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/tanema/gween"
    "github.com/tanema/gween/ease"
)

func (this *Game) LoadInfoScene() {
    println("LoadInfoScene")
    this.curScene.SetId(Info_sceneId)

    // Add star field background
    this.GenerateStars(this.screenView)

    titleText := entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 16),
        this.screenView,
        "MISSION BRIEFING",
    )
    titleText.YAlign = entity.Top_fontaligny

    briefText := []string{
        //"XXXXXXXXXXXXXXXXXXXXXXXXXX",
        "Terra is being threatened by",
        "invasive alien species. Your",
        "mission is to intercept and",
        "destroy them.",
        "",
        "Good luck.",
    }
    totalDelay := 0
    for i := 0; i < len(briefText); i++ {
        line := entity.AddNormalText(
            this.ecs, 
            float64(this.screenView.Area.Max.X / 2), 
            this.screenView.Area.Min.Y + float64(34 + i * 15),
            this.screenView,
            "RedFont",
            briefText[i],
        )
        line.YAlign = entity.Top_fontaligny
        line.TypeWriter = 10
        line.Delay = totalDelay
        totalDelay += (len(briefText[i]) + 2) * line.TypeWriter
    }

    scoreInfoText := entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 130),
        this.screenView,
        "SCORE INFO",
    )
    scoreInfoText.YAlign = entity.Top_fontaligny

    scoreHeaderText := entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 148),
        this.screenView,
        "LightBlueFont",
        "CONVOY  CHARGER",
    )
    scoreHeaderText.YAlign = entity.Top_fontaligny


    scoreDelay := 0
    enemyOffset := 75.0

    enemyXTween1 := gween.New(float32(this.screenView.Area.Min.X - 20), 
                        float32(this.screenView.Area.Max.X / 2 - enemyOffset), 
                        2,
                        ease.Linear)
    entity.AddTweenSprite(
        this.ecs,
        layer.Foreground,
        float64(this.screenView.Area.Min.X - 20), 
        float64(this.screenView.Area.Min.Y + 169),
        enemyXTween1,
        nil,
        scoreDelay,
        this.screenView,
    )

    ptsTextTween1 := gween.New(float32(this.screenView.Area.Max.X + 20), 
                               float32(this.screenView.Area.Min.X + 68), 
                               2,
                               ease.Linear)
    ptsText1 := entity.AddNormalTweenText(
        this.ecs, 
        float64(this.screenView.Area.Max.X + 20), 
        float64(this.screenView.Area.Min.Y + 166),
        ptsTextTween1,
        nil,
        scoreDelay,
        this.screenView,
        "LightBlueFont",
        "60      120 points",
    )
    ptsText1.YAlign = entity.Top_fontaligny
    ptsText1.XAlign = entity.Left_fontalignx
    scoreDelay += 240

    enemyXTween2 := gween.New(float32(this.screenView.Area.Min.X - 20), 
                        float32(this.screenView.Area.Max.X / 2 - enemyOffset), 
                        2,
                        ease.Linear)
    entity.AddTweenSprite(
        this.ecs,
        layer.Foreground,
        float64(this.screenView.Area.Min.X - 20), 
        float64(this.screenView.Area.Min.Y + 187),
        enemyXTween2,
        nil,
        scoreDelay,
        this.screenView,
    )

    ptsTextTween2 := gween.New(float32(this.screenView.Area.Max.X + 20), 
                               float32(this.screenView.Area.Min.X + 68), 
                               2,
                               ease.Linear)
    ptsText2 := entity.AddNormalTweenText(
        this.ecs, 
        float64(this.screenView.Area.Max.X + 20), 
        float64(this.screenView.Area.Min.Y + 184),
        ptsTextTween2,
        nil,
        scoreDelay,
        this.screenView,
        "LightBlueFont",
        "50      100 points",
    )
    ptsText2.YAlign = entity.Top_fontaligny
    ptsText2.XAlign = entity.Left_fontalignx
    scoreDelay += 240

    enemyXTween3 := gween.New(float32(this.screenView.Area.Min.X - 20), 
                        float32(this.screenView.Area.Max.X / 2 - enemyOffset), 
                        2,
                        ease.Linear)
    entity.AddTweenSprite(
        this.ecs,
        layer.Foreground,
        float64(this.screenView.Area.Min.X - 20), 
        float64(this.screenView.Area.Min.Y + 205),
        enemyXTween3,
        nil,
        scoreDelay,
        this.screenView,
    )

    ptsTextTween3 := gween.New(float32(this.screenView.Area.Max.X + 20), 
                               float32(this.screenView.Area.Min.X + 68), 
                               2,
                               ease.Linear)
    ptsText3 := entity.AddNormalTweenText(
        this.ecs, 
        float64(this.screenView.Area.Max.X + 20), 
        float64(this.screenView.Area.Min.Y + 202),
        ptsTextTween3,
        nil,
        scoreDelay,
        this.screenView,
        "LightBlueFont",
        "40      80  points",
    )
    ptsText3.YAlign = entity.Top_fontaligny
    ptsText3.XAlign = entity.Left_fontalignx
    scoreDelay += 240

    enemyXTween4 := gween.New(float32(this.screenView.Area.Min.X - 20), 
                        float32(this.screenView.Area.Max.X / 2 - enemyOffset), 
                        2,
                        ease.Linear)
    entity.AddTweenSprite(
        this.ecs,
        layer.Foreground,
        float64(this.screenView.Area.Min.X - 20), 
        float64(this.screenView.Area.Min.Y + 223),
        enemyXTween4,
        nil,
        scoreDelay,
        this.screenView,
    )

    ptsTextTween4 := gween.New(float32(this.screenView.Area.Max.X + 20), 
                               float32(this.screenView.Area.Min.X + 68), 
                               2,
                               ease.Linear)
    ptsText4 := entity.AddNormalTweenText(
        this.ecs, 
        float64(this.screenView.Area.Max.X + 20), 
        float64(this.screenView.Area.Min.Y + 220),
        ptsTextTween4,
        nil,
        scoreDelay,
        this.screenView,
        "LightBlueFont",
        "30      60  points",
    )
    ptsText4.YAlign = entity.Top_fontaligny
    ptsText4.XAlign = entity.Left_fontalignx
    //scoreDelay += 240

    logoText := entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 240),
        this.screenView,
        "PurpleFont",
        "www.infiniteyak.com",
    )
    logoText.YAlign = entity.Top_fontaligny

    // Advance to the next state when you hit space
    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeySpace,
        func() {
            this.Transition(Advance_sceneEvent)
        },
    )

    // Start game when you hit Enter
    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeyEnter,
        func() {
            this.Transition(GameStart_sceneEvent)
        },
    )
}
