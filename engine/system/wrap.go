package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type wrap struct {
	query *query.Query
}

var Wrap = &wrap{
	query: query.NewQuery(
		filter.Contains(
			component.Wrap,
			component.Position,
			component.View,
		)),
}

func (this *wrap) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		wDist := *component.Wrap.Get(entry).Distance
		pos := component.Position.Get(entry).Point
		view := component.View.Get(entry).View

        if pos.X - wDist >= view.Area.Max.X {
            pos.X = view.Area.Min.X - wDist + 1
        } else if pos.X + wDist <= view.Area.Min.X {
            pos.X = view.Area.Max.X + wDist - 1
        }
        if pos.Y - wDist > view.Area.Max.Y {
            pos.Y = view.Area.Min.Y - wDist + 1
        } else if pos.Y + wDist < view.Area.Min.Y {
            pos.Y = view.Area.Max.Y + wDist - 1
        }
	})
}
