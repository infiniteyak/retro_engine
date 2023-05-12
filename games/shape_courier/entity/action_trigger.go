package shape_courier_entity

import (
	//gMath "math"
	//"math/rand"
	//"strconv"

	"github.com/infiniteyak/retro_engine/engine/component"
	//sc_comp "github.com/infiniteyak/retro_engine/games/shape_courier/component"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	//"github.com/yohamta/donburi/features/math"
    //"math"
)

type actionTriggerData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    position component.PositionData
    view component.ViewData
    collider component.ColliderData
    actions component.ActionsData
}

func AddActionTrigger( ecs *ecs.ECS,
                       x, y float64,
                       actionId component.ActionId,
                       view *utility.View) {
    this := &actionTriggerData{}
    this.ecs = ecs

    entity := this.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.Collider,
        component.Actions,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)

    // Position
    this.position = component.NewPositionData(x, y)
    donburi.SetValue(this.entry, component.Position, this.position)

    //Collider
    this.collider = component.NewColliderData()
    hb := component.NewHitbox(1, 0, 0)
    this.collider.Hitboxes = append(this.collider.Hitboxes, hb)
    donburi.SetValue(this.entry, component.Collider, this.collider)

    // View
    donburi.SetValue(this.entry, component.View, component.ViewData{View:view})

    // Actions
    this.actions = component.NewActions()
    this.actions.AddUpkeepAction(func(){
		c := component.Collider.Get(this.entry)
        for _, target := range c.Collisions {
            if target.HasComponent(component.Actions) {
                targetTriggers := component.Actions.Get(target).TriggerMap
                targetTriggers[actionId] = true
            }
        }
    })
    donburi.SetValue(this.entry, component.Actions, this.actions)

    return
}
