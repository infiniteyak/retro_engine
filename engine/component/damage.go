package component

import (
    "github.com/yohamta/donburi"
)

type DamageData struct {
    Value *float64
    DestroyOnDamage *bool
    OnDamage func()
}

func NewDamageData(v float64) DamageData {
    return DamageData{
        Value: &v,
        DestroyOnDamage: new(bool),
    }
}

var Damage = donburi.NewComponentType[DamageData]()
