package component

import (
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/yohamta/donburi"
)

type InputData struct {
    Mapping map[ActionId]ebiten.Key
}
var Inputs = donburi.NewComponentType[InputData]()
