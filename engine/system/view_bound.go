package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type viewBound struct {
	query *query.Query
}

var ViewBound = &viewBound{
	query: query.NewQuery(
		filter.Contains(
			component.Velocity,
			component.Position,
            component.View,
            component.ViewBound,
		)),
}

func (this *viewBound) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		vel := component.Velocity.Get(entry).Velocity
		pos := component.Position.Get(entry).Point
		view := component.View.Get(entry).View
		vBound := component.ViewBound.Get(entry)
        
        if (pos.X - vBound.XDistance) < (view.Offset.X + view.Area.Min.X) {
            pos.X = view.Offset.X + view.Area.Min.X + vBound.XDistance
            vel.X = 0
        } else if (pos.X + vBound.XDistance) > (view.Offset.X + view.Area.Max.X) {
            pos.X = view.Offset.X + view.Area.Max.X - vBound.XDistance
            vel.X = 0
        }
        //TODO also handle this for Y, or maybe should be a different comp for y?
	})
}
