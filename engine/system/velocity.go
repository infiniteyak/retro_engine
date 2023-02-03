package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type velocity struct {
	query *query.Query
}

var Velocity = &velocity{
	query: query.NewQuery(
		filter.Contains(
			component.Velocity,
			component.Position,
		)),
}

func (this *velocity) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		vel := component.Velocity.Get(entry).Velocity
		pos := component.Position.Get(entry).Point
        pos.X += vel.X
        pos.Y += vel.Y
	})
}
