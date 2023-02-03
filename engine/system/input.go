package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type input struct {
	query *query.Query
}

var Input = &input{
	query: query.NewQuery(
		filter.Contains(
			component.Inputs,
			component.Actions,
		)),
}

func (this *input) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		inputs := component.Inputs.Get(entry)
		acts := component.Actions.Get(entry)
        
        for a, k := range inputs.Mapping {
            acts.TriggerMap[a] = ebiten.IsKeyPressed(k) 
        }
	})
}
