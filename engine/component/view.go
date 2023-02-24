package component

import (
    "github.com/yohamta/donburi"
    "github.com/infiniteyak/retro_engine/engine/utility"
)

// TODO results in a lot of View.View which is ugly...
type ViewData struct {
    View *utility.View
}

var View = donburi.NewComponentType[ViewData]()
