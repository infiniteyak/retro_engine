package shape_courier_entity

import (
	//gMath "math"
	//"math/rand"
	//"strconv"

	"github.com/infiniteyak/retro_engine/engine/component"
	//"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	//"github.com/yohamta/donburi/features/math"
    //"math"
)

const (
    powerPointValue = 50
    powerHealth = 1.0
)

type powerData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    position component.PositionData
    view component.ViewData
    collider component.ColliderData
    actions component.ActionsData
    graphicObject component.GraphicObjectData
}

func AddPower( ecs *ecs.ECS,
               x, y float64,
               view *utility.View) {
    this := &dotData{}
    this.ecs = ecs

    entity := this.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.GraphicObject,
        component.Actions,
        component.Collider,
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

    // Graphic Object
    this.graphicObject = component.NewGraphicObjectData()
    spriteData := component.SpriteData{}
    spriteData.Load("Items", nil)
    spriteData.Play("dot")
    this.graphicObject.Renderables = append(this.graphicObject.Renderables, &spriteData)
    donburi.SetValue(this.entry, component.GraphicObject, this.graphicObject)

    // Actions
    this.actions = component.NewActions()
    this.actions.AddNormalAction(component.Destroy_actionid, func() {
        this.graphicObject.HideAllRenderables(true)
        se := event.Score{Value:powerPointValue}
        event.ScoreEvent.Publish(this.ecs.World, se)

        runEvent := event.RunMode{}
        event.SetRunModeEvent.Publish(this.ecs.World, runEvent)

        ree := event.RemoveEntity{Entity:this.entity}
        event.RemoveEntityEvent.Publish(this.ecs.World, ree)
    })

    this.actions.AddUpkeepAction(func(){
		c := component.Collider.Get(this.entry)
        for _, target := range c.Collisions {
            if target.HasComponent(component.PlayerTag) {
                this.actions.TriggerMap[component.Destroy_actionid] = true
            }
        }
    })

    donburi.SetValue(this.entry, component.Actions, this.actions)

    // View
    donburi.SetValue(this.entry, component.View, component.ViewData{View:view})
}
