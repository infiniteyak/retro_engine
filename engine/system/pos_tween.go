package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type posTween struct {
	query *query.Query
}

var PosTween = &posTween{
	query: query.NewQuery(
		filter.Contains(
			component.PosTween,
			component.Position,
		)),
}

func (this *posTween) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		pt := component.PosTween.Get(entry)
		pos := component.Position.Get(entry).Point

        if pt.Delay > 0 {
            pt.Delay--
        } else {
            if pt.XTween != nil {
                current, isFinished := pt.XTween.Update(float32(1.0/120.0)) //TODO constant?
                pos.X = float64(current)
                if isFinished {
                    pt.XTween = nil
                }
            }
            if pt.YTween != nil {
                current, isFinished := pt.YTween.Update(float32(1.0/120.0)) //TODO constant?
                if !isFinished {
                    pos.Y = float64(current) 
                }
            }
        }
	})
}
