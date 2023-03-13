package game

import (
	"fmt"
	"github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

func (this *Game) LoadInitialsScene() {
    println("LoadInitialsScene")
    this.curScene.SetId(EnterInitials_sceneId)

    this.GenerateStars(this.screenView)

    asset.StopMusic()

    asset.PlaySound("MenuNoise")

    entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2 - 36),
        this.screenView,
        "WhiteFont",
        "SCORE " + fmt.Sprintf("%06d", this.curScore.Value),
    )

    entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2 - 22),
        this.screenView,
        "WhiteFont",
        "ENTER YOUR INITIALS",
    )

    initialsText := entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        this.screenView,
        "WhiteFont",
        "___",
    )

    entity.AddTextInput(
        this.ecs, 
        &initialsText.String,
        3,
        func(){
            this.curScore.Ident = initialsText.String
        },
    )

    confirmText := entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2 + 22),
        this.screenView,
        "WhiteFont",
        "FIRE TO CONFIRM",
    )
    confirmText.Blink = true

    entity.AddInputTrigger(
        this.ecs, 
        ebiten.KeySpace,
        func() {
            // If we've fully filled in the initals advance
            done := true
            for _, c := range initialsText.String {
                if string(c) == "_" {
                    done = false
                }
            }
            if done {
                this.Transition(Advance_sceneEvent)
            }
        },
    )
}
