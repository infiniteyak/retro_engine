package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
	//"github.com/infiniteyak/retro_engine/engine/utility"
	//aEntity "github.com/infiniteyak/retro_engine/games/astralian/entity"
	"strings"
	//"github.com/hajimehoshi/ebiten/v2"
)

func (this *Game) LoadMainScene() {
    println("LoadMainScene")
    this.curScene.SetId(Main_sceneId)

    entity.AddNormalText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Max.Y / 2),
        this.screenView,
        "WhiteFont",
        strings.ToUpper(Title),
    )
}

