package component

import (
    "github.com/yohamta/donburi"
)

type DamageData struct {
    Value *float64
    DestroyOnDamage *bool
    OnDamage func()
}

func NewDamageData() DamageData {
    return DamageData{
        Value: new(float64),
        DestroyOnDamage: new(bool),
    }
}

var Damage = donburi.NewComponentType[DamageData]()
