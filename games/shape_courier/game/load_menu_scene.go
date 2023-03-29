package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
)

func (this *Game) LoadMenuScene() {
    println("LoadMenuScene")
    this.curScene.SetId(Menu_sceneId)

    menuFormat := entity.GameMenuFormat{
        YAlign: entity.Middle_fontaligny,
        XAlign: entity.Center_fontalignx,
        Font: "WhiteFont",
        Kerning: 0,
        Spacing: 10,
        SelectPad: 2,
        SelectSprite: "WhiteSelect",
    }

    menuDisplay := []string{
        "Start Game",
        "Info",
        "High Scores",
        "Attract",
        "Options",
    }

    menuOptions := map[string]func() {
        "Start Game": func() {
            println("Start Game")
            this.Transition(GameStart_sceneEvent)
        },
        "Info": func() {
            println("Info")
            this.Transition(GoInfo_sceneEvent)
        },
        "Attract": func() {
            println("Attract")
            this.Transition(GoAttract_sceneEvent)
        },
        "High Scores": func() {
            println("Scores")
            this.Transition(GoScores_sceneEvent)
        },
        "Options": func() {
            println("Options")
            this.Transition(GoOptions_sceneEvent)
        },
    }

    entity.AddGameMenu(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        menuOptions,
        menuDisplay,
        menuFormat,
        "MenuNoise",
        this.screenView,
    )
}

