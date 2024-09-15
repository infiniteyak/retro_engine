package component

import (
    "github.com/yohamta/donburi"
    "github.com/infiniteyak/retro_engine/engine/utility"
)

type ViewData struct {
    View *utility.View
}

var View = donburi.NewComponentType[ViewData]()
