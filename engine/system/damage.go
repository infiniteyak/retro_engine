package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type damage struct {
	query *query.Query
}

var Damage = &damage{
	query: query.NewQuery(
		filter.Contains(
			component.Damage,
			component.Collider,
			component.Factions,
		)),
}

func (this *damage) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		damage := component.Damage.Get(entry)
		collider := component.Collider.Get(entry)
		factions := component.Factions.Get(entry)
        for _, target := range collider.Collisions {
            if target.HasComponent(component.Health) && target.HasComponent(component.Factions) {
                targetHealth := component.Health.Get(target)
                targetFactions := component.Factions.Get(target)

                // Ignore friendly fire
                hasFaction := false
                for _, f := range factions.Values {
                    if targetFactions.HasFaction(f) {
                        hasFaction = true
                    }
                }

                if !hasFaction {
                    if target.HasComponent(component.Actions) {
                        acts := component.Actions.Get(target)
                        if !acts.TriggerMap[component.Shield_actionid] {
                            targetHealth.Value -= damage.Value
                        }
                    } else {
                        targetHealth.Value -= damage.Value
                    }
                    if damage.DestroyOnDamage && entry.HasComponent(component.Actions) {
                        acts := component.Actions.Get(entry)
                        acts.TriggerMap[component.Destroy_actionid] = true
                    }
                }
            }
        }
	})
}
