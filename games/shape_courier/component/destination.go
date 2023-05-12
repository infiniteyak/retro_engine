package shape_courier_component

import (
    "github.com/yohamta/donburi"
    "github.com/infiniteyak/retro_engine/engine/utility"
)

type DestinationData struct {
    Point *utility.Point
}

var Destination = donburi.NewComponentType[DestinationData]()

func NewDestinationData(x, y float64) DestinationData {
    return DestinationData{Point: &utility.Point{X: x, Y: y}}
}
