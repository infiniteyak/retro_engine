package component

import (
    "github.com/yohamta/donburi"
)

type WrapData struct {
    Distance *float64
}

var Wrap = donburi.NewComponentType[WrapData]()
