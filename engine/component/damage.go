package component

import (
    "github.com/yohamta/donburi"
)

type DamageData struct {
    Value float64
    DestroyOnDamage bool
}

var Damage = donburi.NewComponentType[DamageData]()
