package component

import (
    "github.com/yohamta/donburi"
)

type ViewBoundData struct {
    XDistance float64
    YDistance float64
}

var ViewBound = donburi.NewComponentType[ViewBoundData]()
