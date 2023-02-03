package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

    "github.com/infiniteyak/retro_engine/engine/component"
)

type action struct {
	query *query.Query
}

var Action = &action{
	query: query.NewQuery(
		filter.Contains(
			component.Actions,
		)),
}

func (this *action) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
        if !entry.Valid() { //TODO shoudl this be in every query?
            return
        }
		acts := component.Actions.Get(entry)

        // Tick down cooldowns
        for k, v := range acts.CooldownMap {
            if v.Cur > 0 {
                v.Cur -= 1
                acts.CooldownMap[k] = v
            }
        }

        // Run any triggered & cooled actions
        for k, v := range acts.TriggerMap {
            if v && acts.CooldownMap[k].Cur == 0 {
                if acts.ActionMap[k] != nil {
                    acts.ActionMap[k]()
                }
            }
        }

        // Do any upkeep actions (always runs)
        if acts.ActionMap[component.Upkeep_actionid] != nil {
            acts.ActionMap[component.Upkeep_actionid]()
        }
	})
}
