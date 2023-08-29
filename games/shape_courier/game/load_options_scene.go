package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
	//sc_entity "github.com/infiniteyak/retro_engine/games/shape_courier/entity"
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

    optionsMenu := entity.AddOptionMenu(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        menuFormat,
        "MenuNoise", //TODO move to struct
        this.screenView,
    )

    setSFXVol := func(i int) {
        asset.SetSFXVolume(float64(i)/float64(entity.MaxIncrement))
    }
    optionsMenu.AddOption("SFX Vol", entity.AddSliderOption(this.ecs, 
                                          menuFormat, 
                                          int(asset.GetSFXVolume()*10.0), 
                                          setSFXVol,
                                          this.screenView))

    setMusicVol := func(i int) {
        asset.SetMusicVolume(float64(i)/float64(entity.MaxIncrement))
    }
    optionsMenu.AddOption("Music Vol", entity.AddSliderOption(this.ecs, 
                                          menuFormat, 
                                          int(asset.GetMusicVolume()*10.0), 
                                          setMusicVol,
                                          this.screenView))

    /*
    optionsMenu.AddOption("Lives", entity.AddNumberOption(this.ecs, 
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
                                        this.screenView)) 
    */
    optionsMenu.AddOption("Wave", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(Options.StartingWave), 
                                        1.0, 
                                        999.0, 
                                        1.0, 
                                        1.0, 
                                        3,
                                        func(f float64){
                                            Options.StartingWave = int(f)
                                            this.ResetScore()
                                        },
                                        this.screenView)) 

    /*
    optionsMenu.AddOption("Plyr Spd", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.MandyOptionsData.MoveSpeed), 
                                        0.001, 
                                        0.999, 
                                        0.010, 
                                        1000.0, 
                                        3,
                                        func(f float64){
                                            sc_entity.MandyOptionsData.MoveSpeed = f
                                        },
                                        this.screenView)) 

    optionsMenu.AddOption("Ghst Spd", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.GhostOptionsData.MoveSpeed), 
                                        0.001, 
                                        0.999, 
                                        0.010, 
                                        1000.0, 
                                        3,
                                        func(f float64){
                                            sc_entity.GhostOptionsData.MoveSpeed = f
                                        },
                                        this.screenView)) 

    optionsMenu.AddOption("Ghst Spd Fr", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.GhostOptionsData.MoveSpeedFrighten), 
                                        0.001, 
                                        0.999, 
                                        0.010, 
                                        1000.0, 
                                        3,
                                        func(f float64){
                                            sc_entity.GhostOptionsData.MoveSpeedFrighten = f
                                        },
                                        this.screenView)) 

    optionsMenu.AddOption("Ghst Spd Fa", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.GhostOptionsData.MoveSpeedFast), 
                                        0.001, 
                                        0.999, 
                                        0.010, 
                                        1000.0, 
                                        3,
                                        func(f float64){
                                            sc_entity.GhostOptionsData.MoveSpeedFast = f
                                        },
                                        this.screenView)) 

    optionsMenu.AddOption("Taco Spd", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.TacoOptionsData.MoveSpeed), 
                                        0.001, 
                                        0.999, 
                                        0.010, 
                                        1000.0, 
                                        3,
                                        func(f float64){
                                            sc_entity.TacoOptionsData.MoveSpeed = f
                                        },
                                        this.screenView)) 

    optionsMenu.AddOption("Taco Delay", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.TacoOptionsData.Delay), 
                                        100, 
                                        9000, 
                                        100, 
                                        0.01, 
                                        2,
                                        func(f float64){
                                            sc_entity.TacoOptionsData.Delay = int(f)
                                        },
                                        this.screenView)) 

    optionsMenu.AddOption("Taco Dur", entity.AddNumberOption(this.ecs, 
                                        menuFormat, 
                                        float64(sc_entity.TacoOptionsData.Duration), 
                                        100, 
                                        9000, 
                                        100, 
                                        0.01, 
                                        2,
                                        func(f float64){
                                            sc_entity.TacoOptionsData.Duration = int(f)
                                        },
                                        this.screenView)) 
    */

    optionsMenu.AddOption("Back", entity.AddButtonOption(this.ecs, 
                                       menuFormat, 
                                       func() {this.Transition(Advance_sceneEvent)},
                                       this.screenView))
}

