package component

import (
    "github.com/yohamta/donburi"
    "github.com/infiniteyak/retro_engine/engine/utility"
)

type PositionData struct {
    Point *utility.Point
}

var Position = donburi.NewComponentType[PositionData]()

func NewPositionData(x, y float64) PositionData {
    return PositionData{Point: &utility.Point{X: x, Y: y}}
}
