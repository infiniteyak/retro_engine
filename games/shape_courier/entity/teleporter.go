package shape_courier_entity

import (
	//gMath "math"
	//"math/rand"
	//"strconv"

	"github.com/infiniteyak/retro_engine/engine/component"
	sc_comp "github.com/infiniteyak/retro_engine/games/shape_courier/component"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	//"github.com/yohamta/donburi/features/math"
    //"math"
)

type tpData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    position component.PositionData
    destination sc_comp.DestinationData
    view component.ViewData
    collider component.ColliderData
    graphicObject component.GraphicObjectData
}

func AddTeleporter( ecs *ecs.ECS,
             x, y float64,
             dx, dy float64,
             hbOffsetX, hbOffsetY float64,
             view *utility.View) {
    this := &tpData{}
    this.ecs = ecs

    entity := this.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.GraphicObject,
        component.Collider,
        sc_comp.Destination,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)

    // Position
    this.position = component.NewPositionData(x, y)
    donburi.SetValue(this.entry, component.Position, this.position)

    // Destination
    this.destination = sc_comp.NewDestinationData(dx, dy)
    donburi.SetValue(this.entry, sc_comp.Destination, this.destination)

    //Collider
    this.collider = component.NewColliderData()
    hb := component.NewHitbox(1, hbOffsetX, hbOffsetY)
    this.collider.Hitboxes = append(this.collider.Hitboxes, hb)
    donburi.SetValue(this.entry, component.Collider, this.collider)

    // Graphic Object
    this.graphicObject = component.NewGraphicObjectData()
    spriteData := component.SpriteData{}
    spriteData.Load("Items", nil)
    spriteData.Play("teleporter")
    this.graphicObject.Renderables = append(this.graphicObject.Renderables, &spriteData)
    donburi.SetValue(this.entry, component.GraphicObject, this.graphicObject)

    // View
    donburi.SetValue(this.entry, component.View, component.ViewData{View:view})
}
