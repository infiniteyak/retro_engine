package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/asset"
)

func (this *Game) LoadOptionsScene() {
    println("LoadOptionsScene")
    this.curScene.SetId(Options_sceneId)

    menuFormat := entity.OptionMenuFormat{
        YAlign: entity.Middle_fontaligny,
        XAlign: entity.Inner_optionalignx,
        ItemFont: "WhiteFont",
        SelectFont: "WhiteFont",
        Kerning: 0,
        Spacing: 10,
        SelectPad: 2,
        SelectSprite: "WhiteSelect",
        Gap: 10,
        Title: "OPTIONS",
        TitleFont: "TitleFont",
        TitleKerning: 2,
    }

    //TODO there must be a better design for this
    menuDisplay := []string{
        "SFX Vol",
        "Music Vol",
        "Lives",
        "Back",
    }

    setSFXVol := func(i int) {
        asset.SetSFXVolume(float64(i)/float64(entity.MaxIncrement))
    }
    setMusicVol := func(i int) {
        asset.SetMusicVolume(float64(i)/float64(entity.MaxIncrement))
    }

    menuOptions := map[string]entity.Option {
        "SFX Vol": entity.AddSliderOption(this.ecs, 
                                          menuFormat, 
                                          int(asset.GetSFXVolume()*10.0), 
                                          setSFXVol,
                                          this.screenView,
            ),
        "Music Vol": entity.AddSliderOption(this.ecs, 
                                           menuFormat, 
                                           int(asset.GetMusicVolume()*10.0), 
                                           setMusicVol,
                                           this.screenView,
            ),
        "Lives": entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(Options.StartingLives), 
                                        1.0, 
                                        10.0, 
                                        1.0, 
                                        1.0, 
                                        2,
                                        func(f float64){
                                            Options.StartingLives = int(f)
                                            this.ResetScore()
                                        },
                                        this.screenView,
            ),
        "Back": entity.AddButtonOption(this.ecs, 
                                       menuFormat, 
                                       func() {this.Transition(Advance_sceneEvent)},
                                       this.screenView,
            ),
    }

    entity.AddOptionMenu(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        menuOptions,
        menuDisplay,
        menuFormat,
        "MenuNoise", //TODO move to struct
        this.screenView,
    )
}

