package shape_courier_entity

import (
	"github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/infiniteyak/retro_engine/engine/component"
	"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const (
    tacoPointValue = 1000 //TODO is this correct?
    tacoColliderRadius = 4
    tacoColliderOffsetX = 0
    tacoColliderOffsetY = 0
    tacoSpriteName = "Items"
    tacoSpriteTag = "taco"
    tacoMoveSpeed = 0.1
    tacoDelay = 1500
    tacoDuration = 2000
)

type TacoOptions struct {
    MoveSpeed float64
    Delay int
    Duration int
}

var TacoOptionsData TacoOptions = TacoOptions{
    MoveSpeed: tacoMoveSpeed,
    Delay: tacoDelay,
    Duration: tacoDuration,
}

type TacoData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity
    collider component.ColliderData
    position component.PositionData
    view component.ViewData
    graphicObject component.GraphicObjectData
    spriteData component.SpriteData
    actions component.ActionsData
    mazeData *MazeData
    dir Direction
    targetDir Direction
    allowTp bool
    curR int
    curC int
    delayCountDown int
    durationCountDown int
}

func (this *TacoData) move(direction Direction) {
    speed := TacoOptionsData.MoveSpeed
    *this.position.Point, direction = this.mazeData.ResolveMove(*this.position.Point, direction, speed)
    this.dir = direction
}

func AddTaco (ecs *ecs.ECS,
              view *utility.View,
              mazeData *MazeData) *TacoData{
    this := &TacoData{}
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
    
    this.mazeData = mazeData
    this.dir = South_direction

    // Position
    this.position = component.NewPositionData(mazeData.GetStartPosition())
    donburi.SetValue(this.entry, component.Position, this.position)

    //Collider
    this.collider = component.NewSingleHBCollider(tacoColliderRadius, 
                                                  tacoColliderOffsetX, 
                                                  tacoColliderOffsetY)
    donburi.SetValue(this.entry, component.Collider, this.collider)

    // View
    this.view = component.ViewData{View:view}
    donburi.SetValue(this.entry, component.View, this.view)

    // Graphic Object
    this.delayCountDown = TacoOptionsData.Delay
    this.durationCountDown = TacoOptionsData.Duration
    this.graphicObject = component.NewGraphicObjectData()
    this.spriteData = component.NewSpriteData(tacoSpriteName, nil, tacoSpriteTag)
    this.graphicObject.Renderables = append(this.graphicObject.Renderables, &this.spriteData)
    donburi.SetValue(this.entry, component.GraphicObject, this.graphicObject)

    //START OUT HIDDEN
    *this.spriteData.RenderableData.GetTransInfo().Hide = true

    // Actions
    this.actions = component.NewActions()

    // Move Left
    this.actions.AddNormalAction(component.MoveLeft_actionid, func(){})

    // Move Right
    this.actions.AddNormalAction(component.MoveRight_actionid, func(){})

    // Move Up
    this.actions.AddNormalAction(component.MoveUp_actionid, func(){})

    // Move Down
    this.actions.AddNormalAction(component.MoveDown_actionid, func(){})
    
    // ReadyTeleport
    this.actions.AddNormalAction(component.ReadyTeleport_actionid, func(){
        // taco can't teleport!
        //this.allowTp = true
        this.actions.TriggerMap[component.ReadyTeleport_actionid] = false
    })

    // Teleport
    this.actions.AddNormalAction(component.Teleport_actionid, func(){
    })

    // Destroy (killed)
    this.actions.AddNormalAction(component.Destroy_actionid, func() {
        this.actions.TriggerMap[component.Destroy_actionid] = false
        se := event.Score{Value:tacoPointValue}
        event.ScoreEvent.Publish(this.ecs.World, se)

        *this.spriteData.RenderableData.GetTransInfo().Hide = true
        ree := event.RemoveEntity{Entity:this.entity}
        event.RemoveEntityEvent.Publish(this.ecs.World, ree)

        asset.PlaySound("Crunch")
    })

    // Destroy (despawn)
    this.actions.AddNormalAction(component.DestroySilent_actionid, func() {
        *this.spriteData.RenderableData.GetTransInfo().Hide = true
        ree := event.RemoveEntity{Entity:this.entity}
        event.RemoveEntityEvent.Publish(this.ecs.World, ree)
    })

    this.actions.AddUpkeepAction(func(){
        if this.delayCountDown == 0 {
            this.delayCountDown = -1
            this.spriteData.Play(tacoSpriteTag)
            *this.spriteData.RenderableData.GetTransInfo().Hide = false
        } else if this.delayCountDown > 0{
            this.delayCountDown--
        }

        if this.durationCountDown == 0 {
            this.actions.TriggerMap[component.DestroySilent_actionid] = true
        } else if this.delayCountDown == -1 {
            this.durationCountDown--
        }

        row, col := this.mazeData.FindCoordinates(*this.position.Point)
        if row != this.curR || col != this.curC {
            this.curR = row
            this.curC = col
            this.targetDir = this.mazeData.GetRandomDirection(*this.position.Point, this.dir)
        }

        this.move(this.targetDir)

        c := component.Collider.Get(this.entry)
        for _, target := range c.Collisions {
            if target.HasComponent(component.PlayerTag) && this.delayCountDown == -1 {
                this.actions.TriggerMap[component.Destroy_actionid] = true
            }
        }
    })
    donburi.SetValue(this.entry, component.Actions, this.actions)

    return this
}


