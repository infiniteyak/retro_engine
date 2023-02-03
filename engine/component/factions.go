package component

import (
    "github.com/yohamta/donburi"
)

type FactionId int

const (
    Undefined_factionid FactionId = iota
    Enemy_factionid
    Player_factionid
)

type FactionsData struct {
    Values []FactionId
}

var Factions = donburi.NewComponentType[FactionsData]()

func (this *FactionsData) HasFaction(faction FactionId) bool {
    for _, fac := range this.Values {
        if faction == fac {
            return true
        }
    }
    return false
}
