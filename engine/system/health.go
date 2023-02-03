package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type health struct {
	query *query.Query
}

var Health = &health{
	query: query.NewQuery(
		filter.Contains(
			component.Health,
		)),
}

func (this *health) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		hval := component.Health.Get(entry).Value
        if hval <= 0 {
            if entry.HasComponent(component.Actions) {
                actions := component.Actions.Get(entry).TriggerMap
                actions[component.Destroy_actionid] = true
            } else {
                ecs.World.Remove(entry.Entity())
            }
        }
	})
}
