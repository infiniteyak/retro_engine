package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
    "github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/hajimehoshi/ebiten/v2"
)

func (this *Game) LoadInfoScene() {
    println("LoadInfoScene")
    this.curScene.SetId(Info_sceneId)

    titleText := entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 16),
        this.screenView,
        "INFORMATION",
    )
    titleText.YAlign = entity.Top_fontaligny

    briefText := []string{
        //"XXXXXXXXXXXXXXXXXXXXXXXXXX",
        "You are the shape courier,",
        "a mysterious yet generic",
        "figure who's mission is to",
        "gather missing shapes from",
        "across the universe. In this",
        "case the shapes are in a",
        "haunted space ship for some",
        "reason. Look I'm not a",
        "writer, okay?",
    }
    totalDelay := 0
    for i := 0; i < len(briefText); i++ {
        line := entity.AddNormalText(
            this.ecs, 
            float64(this.screenView.Area.Max.X / 2), 
            this.screenView.Area.Min.Y + float64(34 + i * 15),
            this.screenView,
            "WhiteFont",
            briefText[i],
        )
        line.YAlign = entity.Top_fontaligny
        line.TypeWriter = 7
        line.Delay = totalDelay
        totalDelay += (len(briefText[i]) + 4) * line.TypeWriter
    }

    pointsText := entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 175),
        this.screenView,
        "POINTS",
    )
    pointsText.YAlign = entity.Top_fontaligny

    spriteY := float64(this.screenView.Area.Min.Y + 198)
    textY := float64(this.screenView.Area.Min.Y + 195)
    gapSize := 20.0
    spriteX := 65.0
    textX := 80.0

    //GHOST
    entity.AddSpriteObject(
                this.ecs,
                layer.HudForeground,
                spriteX,
                spriteY, 
                "Ghost",
                "fb_down",
                this.screenView,
                )
    ghostText := entity.AddNormalText(
        this.ecs, 
        textX,
        textY,
        this.screenView,
        "WhiteFont",
        "200 Points",
    )
    ghostText.YAlign = entity.Top_fontaligny
    ghostText.XAlign = entity.Left_fontalignx
    spriteY += gapSize
    textY += gapSize

    //TACO
    entity.AddSpriteObject(
                this.ecs,
                layer.HudForeground,
                spriteX,
                spriteY, 
                "Items",
                "taco",
                this.screenView,
                )
    tacoText := entity.AddNormalText(
        this.ecs, 
        textX,
        textY,
        this.screenView,
        "WhiteFont",
        "1000 Points",
    )
    tacoText.YAlign = entity.Top_fontaligny
    tacoText.XAlign = entity.Left_fontalignx
    spriteY += gapSize
    textY += gapSize

    //Powerup
    entity.AddSpriteObject(
                this.ecs,
                layer.HudForeground,
                spriteX,
                spriteY, 
                "Items",
                "dot",
                this.screenView,
                )
    puText := entity.AddNormalText(
        this.ecs, 
        textX,
        textY,
        this.screenView,
        "WhiteFont",
        "50 Points",
    )
    puText.YAlign = entity.Top_fontaligny
    puText.XAlign = entity.Left_fontalignx
    spriteY += gapSize
    textY += gapSize

    //Dot
    entity.AddSpriteObject(
                this.ecs,
                layer.HudForeground,
                spriteX,
                spriteY, 
                "Items",
                "small_dot",
                this.screenView,
                )
    dotText := entity.AddNormalText(
        this.ecs, 
        textX,
        textY,
        this.screenView,
        "WhiteFont",
        "10 Points",
    )
    dotText.YAlign = entity.Top_fontaligny
    dotText.XAlign = entity.Left_fontalignx

    logoText := entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 278),
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
}
