package component

import (
    "github.com/yohamta/donburi"
)

type DamageData struct {
    Value *float64
    DestroyOnDamage bool
    OnDamage func() //TODO better name?
}

var Damage = donburi.NewComponentType[DamageData]()
