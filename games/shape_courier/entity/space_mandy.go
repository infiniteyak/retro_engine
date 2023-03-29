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
	// "github.com/yohamta/donburi/features/math"
	"github.com/hajimehoshi/ebiten/v2"
)

type mandyData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity

    factions component.FactionsData
    damage component.DamageData
    //health component.HealthData
    collider component.ColliderData
    position component.PositionData
    view component.ViewData
    //velocity component.VelocityData
    graphicObject component.GraphicObjectData
    actions component.ActionsData
    inputs component.InputsData

    dir Direction
}

func AddSpaceMandy( ecs *ecs.ECS,
              view *utility.View,
              md *MazeData) {
    this := &mandyData{}
    this.ecs = ecs

    entity := this.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.GraphicObject,
        component.Inputs,
        component.Actions,
        //component.Velocity,
        component.Collider,
        // component.Health,
        component.Factions,
        component.Damage,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)
    
    // Velocity
    /*
    this.velocity = component.VelocityData{Velocity: &math.Vec2{}} //TODO add init func
    donburi.SetValue(this.entry, component.Velocity, this.velocity)
    */

    // Factions
    factions := []component.FactionId{component.Player_factionid}
    this.factions = component.FactionsData{Values: factions}
    donburi.SetValue(this.entry, component.Factions, this.factions)

    // Damage
    this.damage = component.NewDamageData() //TODO fix
    *this.damage.Value = 1.0
    donburi.SetValue(this.entry, component.Damage, this.damage)

    // Position
    x, y := md.GetStartPosition()
    this.position = component.NewPositionData(x, y)
    donburi.SetValue(this.entry, component.Position, this.position)

    //Collider
    this.collider = component.NewColliderData()
    hb := component.NewHitbox(4, 0, 0)
    this.collider.Hitboxes = append(this.collider.Hitboxes, hb)
    donburi.SetValue(this.entry, component.Collider, this.collider)

    // View
    this.view = component.ViewData{View:view}
    donburi.SetValue(this.entry, component.View, this.view)

    // Graphic Object
    this.graphicObject = component.NewGraphicObjectData() //TODO needs init functions?
    playerSd := component.SpriteData{}
    playerSd.Load("SpaceMandy", nil)
    playerSd.Play("stand_down")
    this.graphicObject.Renderables = append(this.graphicObject.Renderables, &playerSd)
    donburi.SetValue(this.entry, component.GraphicObject, this.graphicObject)

    // Inputs
    this.inputs = component.NewInput()
    this.inputs.AddContinuousInput(component.MoveLeft_actionid, ebiten.KeyLeft)
    this.inputs.AddContinuousInput(component.MoveRight_actionid, ebiten.KeyRight)
    this.inputs.AddContinuousInput(component.MoveUp_actionid, ebiten.KeyUp)
    this.inputs.AddContinuousInput(component.MoveDown_actionid, ebiten.KeyDown)
    donburi.SetValue(this.entry, component.Inputs, this.inputs)

    // Actions
    this.actions = component.NewActions()

    // Move Left
    //TODO clean up
    moveSpeed := 0.3 //TODO make this a const
    this.actions.AddNormalAction(component.MoveLeft_actionid, func(){
        //this.velocity.Velocity.X = -1.0 * moveSpeed //TODO accessor functions?
        *this.position.Point = md.ResolveMove(*this.position.Point, West_direction, moveSpeed)
        playerSd.Play("move_left") //TODO is there was better way?
        this.dir = West_direction
    })

    // Move Right
    this.actions.AddNormalAction(component.MoveRight_actionid, func(){
        *this.position.Point = md.ResolveMove(*this.position.Point, East_direction, moveSpeed)
        playerSd.Play("move_right") //TODO is there was better way?
        this.dir = East_direction
    })

    // Move Up
    this.actions.AddNormalAction(component.MoveUp_actionid, func(){
        *this.position.Point = md.ResolveMove(*this.position.Point, North_direction, moveSpeed)
        playerSd.Play("move_up") //TODO is there was better way?
        this.dir = North_direction
    })
    // Move Up
    this.actions.AddNormalAction(component.MoveDown_actionid, func(){
        *this.position.Point = md.ResolveMove(*this.position.Point, South_direction, moveSpeed)
        playerSd.Play("move_down") //TODO is there was better way?
        this.dir = South_direction
    })

    this.actions.AddUpkeepAction(func(){
        // if they're both on or both off...
        if !this.actions.TriggerMap[component.MoveRight_actionid] &&
           !this.actions.TriggerMap[component.MoveLeft_actionid] &&
           !this.actions.TriggerMap[component.MoveUp_actionid] &&
           !this.actions.TriggerMap[component.MoveDown_actionid] {
            switch this.dir {
            case East_direction:
                playerSd.Play("stand_right")
            case West_direction:
                playerSd.Play("stand_left")
            case North_direction:
                playerSd.Play("stand_up")
            case South_direction:
                playerSd.Play("stand_down")
            }
        }
    })

    donburi.SetValue(this.entry, component.Actions, this.actions)
 
    return 
}
