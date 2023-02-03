package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

    "github.com/infiniteyak/retro_engine/engine/component"
)

type animateGraphicObjects struct {
	query *query.Query
}

var AnimateGraphicObjects = &animateGraphicObjects{
	query: query.NewQuery(
		filter.Contains(
			component.GraphicObject,
		)),
}

func (a *animateGraphicObjects) Update(ecs *ecs.ECS) {
	a.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
        renderables := component.GraphicObject.Get(entry).Renderables
        for _, r := range renderables {
            r.Animate()
        }
	})
}
