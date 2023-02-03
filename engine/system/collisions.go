package system

import (
    "math"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type collisions struct {
	query *query.Query
}

var Collisions = &collisions{
	query: query.NewQuery(
		filter.Contains(
			component.View,
			component.Position,
			component.Collider,
		)),
}

func (r *collisions) Update(ecs *ecs.ECS) {
    var entries []*donburi.Entry
	r.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
        hb := component.Collider.Get(entry)
        hb.Collisions = []*donburi.Entry{}
        entries = append(entries, entry)
	})

    // TODO do something smarter so you don't compare everything to everything else...
    // not really needed for this game, but eventually we'll want it
    for i, entryA := range entries {
		positionA := component.Position.Get(entryA)
        colliderA := component.Collider.Get(entries[i])
        viewA := component.View.Get(entryA)
        for j, entryB := range entries {
            if entryB != nil && entryA != entryB {
                positionB := component.Position.Get(entryB)
                colliderB := component.Collider.Get(entries[j])
                viewB := component.View.Get(entryB)
                if viewA.View == viewB.View {
                    if CheckHitboxCollision(colliderA, colliderB, positionA, positionB) {
                        colliderA.Collisions = append(colliderA.Collisions, entries[j])
                        colliderB.Collisions = append(colliderB.Collisions, entries[i])
                    }
                }
            }
        }
        entries[i] = nil
    }
}

func CheckHitboxCollision(colliderA, colliderB *component.ColliderData, positionA, positionB *component.PositionData) bool {
    for _, hitboxA := range colliderA.Hitboxes {
        for _, hitboxB := range colliderB.Hitboxes {
            pax := positionA.Point.X + hitboxA.Offset.X
            pay := positionA.Point.Y + hitboxA.Offset.Y
            pbx := positionB.Point.X + hitboxB.Offset.X
            pby := positionB.Point.Y + hitboxB.Offset.Y
            if circlehit(pax, pay, hitboxA.Radius, pbx, pby, hitboxB.Radius) {
                return true
            }
        }
    }
    return false
}

// TODO probably just use a library for this
func circlehit(x1, y1 float64, r1 int, x2, y2 float64, r2 int) bool {
    dist := math.Sqrt(math.Pow(x1-x2, 2.0) + math.Pow(y1-y2, 2.0))
    return dist < float64(r1+r2)
}
