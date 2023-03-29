package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
        
        for a, k := range inputs.KeyMap {
            switch inputs.TypeMap[a] {
            default:
                fallthrough
            case component.Undefined_inputtypeid:
                println("Unsupported input type %v", inputs.TypeMap[a])
                fallthrough
            case component.Continuous_inputtypeid:
                acts.TriggerMap[a] = ebiten.IsKeyPressed(k) 
            case component.Limited_inputtypeid:
                if inpututil.KeyPressDuration(k) == 1 {
                    acts.TriggerMap[a] = true
                }
            case component.Hybrid_inputtypeid:
                d := inpututil.KeyPressDuration(k)
                if d == 1 {
                    acts.TriggerMap[a] = true
                } else if (d >= inputs.DelayMap[a] && ((d - inputs.DelayMap[a]) % inputs.FrequencyMap[a]) == 0) {
                    acts.TriggerMap[a] = true
                }
            }
        }
	})
}
