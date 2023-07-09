package shape_courier_entity

import (
	//gMath "math"
	//"math/rand"
	//"strconv"
	sc_comp "github.com/infiniteyak/retro_engine/games/shape_courier/component"

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

const (
    spaceMandyDamage = 1.0

    spaceMandyColliderRadius = 4
    spaceMandyColliderOffsetX = 0
    spaceMandyColliderOffsetY = 0

    spaceMandySpriteName = "SpaceMandy"
    spaceMandySpriteInitialTag = "stand_down"
    spaceMandySpriteMoveLeftTag = "move_left"
    spaceMandySpriteMoveRightTag = "move_right"
    spaceMandySpriteMoveUpTag = "move_up"
    spaceMandySpriteMoveDownTag = "move_down"
    spaceMandySpriteIdleLeftTag = "stand_left"
    spaceMandySpriteIdleRightTag = "stand_right"
    spaceMandySpriteIdleUpTag = "stand_up"
    spaceMandySpriteIdleDownTag = "stand_down"
    spaceMandySpriteDeathTag = "death"
    spaceMandySpriteDeadTag = "dead"

    spaceMandyMoveSpeed = 0.625
    //spaceMandyTeleportCd = 400
)

type MandyData struct {
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
    spriteData component.SpriteData
    actions component.ActionsData
    inputs component.InputsData
    mazeData *MazeData

    dir Direction
    
    disableControls bool
    tpDestination sc_comp.DestinationData
    allowTp bool
}

var spaceMandyDirMoveMap = map[Direction]string {
    North_direction: spaceMandySpriteMoveUpTag,
    South_direction: spaceMandySpriteMoveDownTag,
    East_direction: spaceMandySpriteMoveRightTag,
    West_direction: spaceMandySpriteMoveLeftTag,
}

var spaceMandyDirIdleMap = map[Direction]string {
    North_direction: spaceMandySpriteIdleUpTag,
    South_direction: spaceMandySpriteIdleDownTag,
    East_direction: spaceMandySpriteIdleRightTag,
    West_direction: spaceMandySpriteIdleLeftTag,
}

func (this *MandyData) move(direction Direction) {
    *this.position.Point, direction = this.mazeData.ResolveMove(*this.position.Point, direction, spaceMandyMoveSpeed)
    this.spriteData.Play(spaceMandyDirMoveMap[direction])
    this.dir = direction
}

func AddSpaceMandy( ecs *ecs.ECS,
              view *utility.View,
              md *MazeData) *MandyData {
    println("ADDING PLAYER")

    this := &MandyData{}
    this.ecs = ecs

    entity := this.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.GraphicObject,
        component.Inputs,
        component.Actions,
        component.Collider,
        // component.Health,
        component.Factions,
        component.Damage,
        component.PlayerTag,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)
    
    this.mazeData = md
    this.dir = South_direction

    // Factions
    this.factions = component.NewSingleFaction(component.Player_factionid)
    donburi.SetValue(this.entry, component.Factions, this.factions)

    // Damage
    this.damage = component.NewDamageData(spaceMandyDamage)
    donburi.SetValue(this.entry, component.Damage, this.damage)

    // Position
    this.position = component.NewPositionData(md.GetStartPosition())
    donburi.SetValue(this.entry, component.Position, this.position)

    //Collider
    this.collider = component.NewSingleHBCollider(spaceMandyColliderRadius, 
                                                  spaceMandyColliderOffsetX, 
                                                  spaceMandyColliderOffsetY)
    donburi.SetValue(this.entry, component.Collider, this.collider)

    // View
    this.view = component.ViewData{View:view}
    donburi.SetValue(this.entry, component.View, this.view)

    // Graphic Object
    this.graphicObject = component.NewGraphicObjectData()
    this.spriteData = component.NewSpriteData(spaceMandySpriteName, nil, spaceMandySpriteInitialTag)
    this.graphicObject.Renderables = append(this.graphicObject.Renderables, &this.spriteData)
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
    this.actions.AddNormalAction(component.MoveLeft_actionid, func(){})

    // Move Right
    this.actions.AddNormalAction(component.MoveRight_actionid, func(){})

    // Move Up
    this.actions.AddNormalAction(component.MoveUp_actionid, func(){})

    // Move Down
    this.actions.AddNormalAction(component.MoveDown_actionid, func(){})
    
    // ReadyTeleport
    this.actions.AddNormalAction(component.ReadyTeleport_actionid, func(){
        this.allowTp = true
        this.actions.TriggerMap[component.ReadyTeleport_actionid] = false
    })

    // Destroy (killed)
    this.actions.AddNormalAction(component.Destroy_actionid, func(){
        this.actions.TriggerMap[component.Destroy_actionid] = false
        this.disableControls = true

        //adjust lives count (and note if it's game over)

        //despawn ghosts
        event.DespawnAllEnemiesEvent.Publish(this.ecs.World, event.DespawnAllEnemies{})

        //play death animation
        //on loop 
        //  if it was game over, do game over stuff
        //  otherwise reset for new life
        this.spriteData.Play(spaceMandySpriteDeathTag)
        this.spriteData.SetLoopCallback(func() {
            this.spriteData.Play(spaceMandySpriteDeadTag)
            this.spriteData.SetLoopCallback(nil)

            event.AdjustLivesEvent.Publish(
                this.ecs.World, 
                event.AdjustLives{
                    Value: -1,
                },
            )

            ree := event.RemoveEntity{Entity:this.entity}
            event.RemoveEntityEvent.Publish(this.ecs.World, ree)
        })
    })

    // Teleport
    this.actions.AddNormalAction(component.Teleport_actionid, func(){
        if this.allowTp {
            this.allowTp = false
            *this.collider.Enable = false
            this.disableControls = true
            this.spriteData.Play("teleport_out")
            this.spriteData.SetLoopCallback(func() {
                //*this.spriteData.RenderableData.GetTransInfo().Hide = true
                this.spriteData.Play("teleport_in")
                //this.position.Point.X -= 20
                this.position.Point.X = this.tpDestination.Point.X
                this.position.Point.Y = this.tpDestination.Point.Y
                this.spriteData.SetLoopCallback(func() {
                    this.spriteData.SetLoopCallback(nil)
                    this.disableControls = false
                    *this.collider.Enable = true
                    this.spriteData.Play(spaceMandySpriteInitialTag)
                })
            })
        }
        this.actions.TriggerMap[component.Teleport_actionid] = false
        //this.actions.ResetCooldown(component.Teleport_actionid)
    })

    this.actions.AddUpkeepAction(func(){
        if !this.disableControls {
            if this.actions.TriggerMap[component.MoveRight_actionid] {
                this.move(East_direction)
            } else if this.actions.TriggerMap[component.MoveLeft_actionid] {
                this.move(West_direction)
            } else if this.actions.TriggerMap[component.MoveUp_actionid] {
                this.move(North_direction)
            } else if this.actions.TriggerMap[component.MoveDown_actionid] {
                this.move(South_direction)
            }

            if !this.actions.TriggerMap[component.MoveRight_actionid] &&
               !this.actions.TriggerMap[component.MoveLeft_actionid] &&
               !this.actions.TriggerMap[component.MoveUp_actionid] &&
               !this.actions.TriggerMap[component.MoveDown_actionid] {
                this.spriteData.Play(spaceMandyDirIdleMap[this.dir])
            }
        } 

		c := component.Collider.Get(this.entry)
        for _, target := range c.Collisions {
            if target.HasComponent(sc_comp.Destination) {
                if this.actions.CooldownMap[component.Teleport_actionid].Cur == 0 {
                    this.actions.TriggerMap[component.Teleport_actionid] = true
                    this.tpDestination = *sc_comp.Destination.Get(target)
                }
            }
        }
    })

    donburi.SetValue(this.entry, component.Actions, this.actions)
 
    return this
}
