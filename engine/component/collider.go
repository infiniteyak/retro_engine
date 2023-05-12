package component

import (
    "github.com/yohamta/donburi"
    "github.com/infiniteyak/retro_engine/engine/utility"
)

type Hitbox struct {
    Radius int
    Offset utility.Point
}

type ColliderData struct {
    Hitboxes []*Hitbox
    Collisions []*donburi.Entry
    Enable *bool
}

var Collider = donburi.NewComponentType[ColliderData]()

func NewColliderData() ColliderData {
    hbs := []*Hitbox{}
    cols := []*donburi.Entry{}
    enable := true
    return ColliderData{
        Hitboxes: hbs,
        Collisions: cols,
        Enable: &enable,
    }
}

func NewSingleHBCollider(radius int, x, y float64) ColliderData {
    hbs := []*Hitbox{}
    hbs = append(hbs, NewHitbox(radius, x, y))
    cols := []*donburi.Entry{}
    enable := true
    return ColliderData{
        Hitboxes: hbs,
        Collisions: cols,
        Enable: &enable,
    }
}

func NewHitbox(radius int, x, y float64) *Hitbox {
    return &Hitbox{
        Radius: radius,
        Offset: utility.Point{
            X: x,
            Y: y,
        },
    }
}
