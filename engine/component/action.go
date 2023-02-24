package component

import (
    "github.com/yohamta/donburi"
)

type ActionId int

// TODO fix this by having the LAST value by the start of the per game ID space
// For now this needs to have all actions from all games
// Otherwise I need a more complicated solution for these IDs
const (
    Undefined_actionid ActionId = iota
    Upkeep_actionid //TODO this is kind of hacky, no? should just make this a separate component
    MoveLeft_actionid
    MoveRight_actionid
    MoveDown_actionid
    MoveUp_actionid
    RotateCW_actionid
    RotateCCW_actionid
    Accelerate_actionid
    Shoot_actionid
    ShootSecondary_actionid
    TriggerFunction_actionid
    SelfDestruct_actionid
    DestroySilent_actionid
    Destroy_actionid
    Shield_actionid
    Blink_actionid
    Charge_actionid
    SendShip_actionid
    ReturnShip_actionid
    ReturnProjectile_actionid
    Reload_actionid
    ReloadSecondary_actionid
    IncreasePower_actionid
    ResetPower_actionid
    Follow_actionid
)

type ActionsData struct {
    TriggerMap map[ActionId]bool
    CooldownMap map[ActionId]Cooldown
    ActionMap map[ActionId]func()
}

var Actions = donburi.NewComponentType[ActionsData]()

type Cooldown struct {
    Cur int
    Max int
}

func (this *Cooldown) Reset() {
    this.Cur = this.Max
}

func NewActions() ActionsData {
    tm := make(map[ActionId]bool)
    cdm := make(map[ActionId]Cooldown)
    am := make(map[ActionId]func())
    return ActionsData {
        TriggerMap: tm,
        CooldownMap: cdm,
        ActionMap: am,
    }
}
