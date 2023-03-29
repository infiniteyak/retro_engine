package game

import (
	"github.com/infiniteyak/retro_engine/engine/entity"
)

func (this *Game) LoadInfoScene() {
    println("LoadInfoScene")
    this.curScene.SetId(Info_sceneId)

    titleText := entity.AddTitleText(
        this.ecs, 
        float64(this.screenView.Area.Max.X / 2), 
        float64(this.screenView.Area.Min.Y + 16),
        this.screenView,
        "INFO SCENE",
    )
    titleText.YAlign = entity.Top_fontaligny
}
