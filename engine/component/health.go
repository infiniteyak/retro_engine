package component

import (
    "github.com/yohamta/donburi"
)

type HealthData struct {
    Value float64
}

var Health = donburi.NewComponentType[HealthData]()
