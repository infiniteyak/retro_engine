package component

import (
    "github.com/yohamta/donburi"
    "github.com/tanema/gween"
)

type PosTweenData struct {
    XTween *gween.Tween
    YTween *gween.Tween
    Delay int
}

var PosTween = donburi.NewComponentType[PosTweenData]()
