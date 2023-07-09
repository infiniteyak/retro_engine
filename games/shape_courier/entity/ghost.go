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
    "github.com/infiniteyak/retro_engine/engine/asset"
	// "github.com/yohamta/donburi/features/math"
	//"github.com/hajimehoshi/ebiten/v2"
    //"math"
    "image/color"
)

const (
    ghostDamage = 1.0

    ghostPointValue = 200 //TODO is this correct?

    ghostRespawnDelay = 500 //TODO is this correct?

    ghostColliderRadius = 4
    ghostColliderOffsetX = 0
    ghostColliderOffsetY = 0

    ghostSpriteName = "Ghost"
    ghostSpriteInitialTag = "appear"
    ghostSpriteMoveLeftTag = "left"
    ghostSpriteMoveRightTag = "right"
    ghostSpriteMoveUpTag = "up"
    ghostSpriteMoveDownTag = "down"
    ghostSpriteIdleLeftTag = "left"
    ghostSpriteIdleRightTag = "right"
    ghostSpriteIdleUpTag = "up"
    ghostSpriteIdleDownTag = "down"
    ghostSpriteDeathTag = "death"

    ghostMoveSpeed = 0.6
    ghostMoveSpeedFrighten = 0.375
    ghostMoveSpeedFast = 0.65

    ghostFrightenTime = 960
)

type AiMode int
const (
    Undefined_aimode AiMode = iota
    Scatter_aimode
    Chase_aimode
)

type GhostVarient int
const (
    Undefined_ghostvarient GhostVarient = iota
    ClassicRed_ghostvarient 
    ClassicPink_ghostvarient 
    ClassicBlue_ghostvarient 
    ClassicOrange_ghostvarient 
)

type GhostData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity

    factions component.FactionsData
    damage component.DamageData
    collider component.ColliderData
    position component.PositionData
    view component.ViewData
    graphicObject component.GraphicObjectData
    spriteData component.SpriteData
    actions component.ActionsData
    mazeData *MazeData

    dir Direction
    targetDir Direction
    
    tpDestination sc_comp.DestinationData
    allowTp bool
    curR int
    curC int

    varient GhostVarient
    delayCountDown int
    animating bool
    scatterPoint utility.Point
    aiMode AiMode
    runMode bool
    fastMode bool
    runTimer int
}

var ghostDirMoveMap = map[Direction]string {
    North_direction: ghostSpriteMoveUpTag,
    South_direction: ghostSpriteMoveDownTag,
    East_direction: ghostSpriteMoveRightTag,
    West_direction: ghostSpriteMoveLeftTag,
}

var ghostDirIdleMap = map[Direction]string {
    North_direction: ghostSpriteIdleUpTag,
    South_direction: ghostSpriteIdleDownTag,
    East_direction: ghostSpriteIdleRightTag,
    West_direction: ghostSpriteIdleLeftTag,
}

var ghostColorMap = map[GhostVarient]color.NRGBA {
    ClassicRed_ghostvarient: color.NRGBA{219, 65, 97, 255},
    ClassicPink_ghostvarient: color.NRGBA{243, 97, 255, 255},
    ClassicBlue_ghostvarient: color.NRGBA{146, 211, 255, 255},
    ClassicOrange_ghostvarient: color.NRGBA{255, 121, 48, 255},
}

var ghostColorNameMap = map[GhostVarient]string {
    ClassicRed_ghostvarient: "RedGhost",
    ClassicPink_ghostvarient: "PinkGhost",
    ClassicBlue_ghostvarient: "BlueGhost",
    ClassicOrange_ghostvarient: "OrangeGhost",
}

var ghostColorDelayMap = map[GhostVarient]int {
    ClassicRed_ghostvarient: 25,
    ClassicPink_ghostvarient: 100,//200,
    ClassicBlue_ghostvarient: 175,//350,
    ClassicOrange_ghostvarient: 250, //500,
}

func (this *GhostData) reverse() {
    switch this.dir {
    case North_direction:
        this.dir = South_direction
    case South_direction:
        this.dir = North_direction
    case East_direction:
        this.dir = West_direction
    case West_direction:
        this.dir = East_direction
    }
    this.targetDir = this.dir
}

func (this *GhostData) move(direction Direction) {
    speed := ghostMoveSpeed
    if this.runMode {
        speed = ghostMoveSpeedFrighten
    } else if this.fastMode {
        speed = ghostMoveSpeedFast
    }
    *this.position.Point, direction = this.mazeData.ResolveMove(*this.position.Point, direction, speed)
    if this.runMode {
        this.spriteData.Play("f_" + ghostDirMoveMap[direction])
    } else {
        this.spriteData.Play(ghostDirMoveMap[direction])
    }
    this.dir = direction
}

func AddGhost(ecs *ecs.ECS,
              //x, y float64,
              view *utility.View,
              mandyData *MandyData,
              mazeData *MazeData,
              ghostType GhostVarient,
              redGhostData *GhostData,
              wave int) *GhostData{
    this := &GhostData{}
    this.ecs = ecs

    entity := this.ecs.Create(
        layer.Foreground, 
        component.Position, 
        component.View,
        component.GraphicObject,
        //component.Inputs,
        component.Actions,
        component.Collider,
        // component.Health,
        component.Factions,
        component.Damage,
        )
    this.entity = &entity

    event.RegisterEntityEvent.Publish(this.ecs.World, event.RegisterEntity{Entity:this.entity})
    this.entry = this.ecs.World.Entry(*this.entity)
    
    this.mazeData = mazeData
    this.dir = South_direction

    // Factions
    this.factions = component.NewSingleFaction(component.Enemy_factionid)
    donburi.SetValue(this.entry, component.Factions, this.factions)

    // Damage
    this.damage = component.NewDamageData(ghostDamage)
    donburi.SetValue(this.entry, component.Damage, this.damage)

    // Position
    this.position = component.NewPositionData(mazeData.GetSpawnPosition())
    donburi.SetValue(this.entry, component.Position, this.position)

    //Collider
    this.collider = component.NewSingleHBCollider(ghostColliderRadius, 
                                                  ghostColliderOffsetX, 
                                                  ghostColliderOffsetY)
    donburi.SetValue(this.entry, component.Collider, this.collider)

    // View
    this.view = component.ViewData{View:view}
    donburi.SetValue(this.entry, component.View, this.view)

    // Graphic Object
    this.varient = ghostType
    this.delayCountDown = ghostColorDelayMap[ghostType]
    this.animating = true
    _, ok := asset.SpriteAssets[ghostColorNameMap[ghostType]]
    if !ok {
        asset.DuplicateSpriteAsset("Ghost", ghostColorNameMap[ghostType])
        baseColor := color.NRGBA{255, 255, 255, 255} //TODO const?
        //ghostColor := color.NRGBA{146, 211, 255, 255}
        ghostColor := ghostColorMap[ghostType]
        asset.SwapColor(ghostColorNameMap[ghostType], baseColor, ghostColor) //TODO
    }

    this.graphicObject = component.NewGraphicObjectData()
    this.spriteData = component.NewSpriteData(ghostColorNameMap[ghostType], nil, ghostSpriteIdleDownTag)
    this.graphicObject.Renderables = append(this.graphicObject.Renderables, &this.spriteData)
    donburi.SetValue(this.entry, component.GraphicObject, this.graphicObject)

    switch this.varient {
    case ClassicRed_ghostvarient:
        this.scatterPoint = this.mazeData.GetNECorner()
    case ClassicPink_ghostvarient:
        this.scatterPoint = this.mazeData.GetNWCorner()
    case ClassicBlue_ghostvarient:
        this.scatterPoint = this.mazeData.GetSECorner()
    case ClassicOrange_ghostvarient:
        this.scatterPoint = this.mazeData.GetSWCorner()
    }
    this.aiMode = Scatter_aimode

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
        // Ghost can't teleport!
        //this.allowTp = true
        this.actions.TriggerMap[component.ReadyTeleport_actionid] = false
    })

    // Teleport
    this.actions.AddNormalAction(component.Teleport_actionid, func(){
    })

    // Destroy (killed)
    this.actions.AddNormalAction(component.Destroy_actionid, func() {
        this.actions.TriggerMap[component.Destroy_actionid] = false
        println("D")
        se := event.Score{Value:ghostPointValue}
        event.ScoreEvent.Publish(this.ecs.World, se)

        this.animating = true
        this.runTimer = 0
        this.spriteData.Play(ghostSpriteDeathTag)
        this.spriteData.SetLoopCallback(func() {
            //this.animating = false
            //this.spriteData.Play(ghostDirMoveMap[South_direction])
            //TODO not this
            *this.spriteData.RenderableData.GetTransInfo().Hide = true
            /*
            ree := event.RemoveEntity{Entity:this.entity}
            event.RemoveEntityEvent.Publish(this.ecs.World, ree)
            */
            this.delayCountDown = ghostRespawnDelay
            this.position.Point.X, this.position.Point.Y = mazeData.GetSpawnPosition()
            this.spriteData.SetLoopCallback(nil)
        })
    })

    // Destroy (despawn)
    this.actions.AddNormalAction(component.DestroySilent_actionid, func() {
        println("DS")
        *this.spriteData.RenderableData.GetTransInfo().Hide = true

        ree := event.RemoveEntity{Entity:this.entity}
        event.RemoveEntityEvent.Publish(this.ecs.World, ree)
    })

    this.actions.AddUpkeepAction(func(){
        if this.delayCountDown == 0 {
            this.delayCountDown = -1
            this.spriteData.Play(ghostSpriteInitialTag)
            *this.spriteData.RenderableData.GetTransInfo().Hide = false
            this.spriteData.SetLoopCallback(func() {
                println("X")
                this.animating = false
                this.spriteData.Play(ghostDirMoveMap[South_direction])
                this.spriteData.SetLoopCallback(nil)
            })
        } else {
            this.delayCountDown--
        }
        if this.runTimer == 0 {
            this.runTimer = -1
            this.runMode = false
            this.reverse()
        } else {
            this.runTimer--
        }

        if !this.animating {
            // TODO AI stuff goes here, need to support more ghost colors
            // also scatter mode, and chase mode etc
            var targetPoint utility.Point
            switch this.aiMode {
            case Scatter_aimode:
                targetPoint = this.scatterPoint
            case Chase_aimode:
                switch this.varient {
                case ClassicRed_ghostvarient:
                    targetPoint = *mandyData.position.Point
                case ClassicPink_ghostvarient:
                    targetPoint = *mandyData.position.Point
                    if mandyData.dir == North_direction {
                        targetPoint.Y -= 4 * wallSpriteHeight
                    } else if mandyData.dir == South_direction {
                        targetPoint.Y += 4 * wallSpriteHeight
                    } else if mandyData.dir == East_direction {
                        targetPoint.X += 4 * wallSpriteWidth
                    } else if mandyData.dir == South_direction {
                        targetPoint.X -= 4 * wallSpriteWidth
                    }
                case ClassicBlue_ghostvarient: //TODO need some way to verify this is right
                    rotationPoint := *mandyData.position.Point
                    if mandyData.dir == North_direction {
                        rotationPoint.Y -= 2 * wallSpriteHeight
                    } else if mandyData.dir == South_direction {
                        rotationPoint.Y += 2 * wallSpriteHeight
                    } else if mandyData.dir == East_direction {
                        rotationPoint.X += 2 * wallSpriteWidth
                    } else if mandyData.dir == South_direction {
                        rotationPoint.X -= 2 * wallSpriteWidth
                    }
                    //rotationPoint.X + VAL = *redGhostData.position.Point.X
                    valX := redGhostData.position.Point.X - rotationPoint.X
                    valY := redGhostData.position.Point.Y - rotationPoint.Y
                    targetPoint.X = rotationPoint.X - valX
                    targetPoint.Y = rotationPoint.Y - valY
                case ClassicOrange_ghostvarient:
                    distToPlayer := distance(*mandyData.position.Point, *this.position.Point)
                    boundary := 8 * wallSpriteHeight
                    if distToPlayer > boundary {
                        targetPoint = *mandyData.position.Point
                    } else {
                        targetPoint = this.scatterPoint
                    }
                }
            }

            //We need to do this so that the ghost can't double back at intersections
            row, col := this.mazeData.FindCoordinates(*this.position.Point)
            if row != this.curR || col != this.curC {
                this.curR = row
                this.curC = col

                if this.runMode {
                    this.targetDir = this.mazeData.GetRandomDirection(*this.position.Point, this.dir)
                } else {
                    //TODO move AI stuff in here?
                    this.targetDir = this.mazeData.GetDirectionToTarget(*this.position.Point, targetPoint, this.dir)
                }
            }

            this.move(this.targetDir)

            c := component.Collider.Get(this.entry)
            for _, target := range c.Collisions {
                if !this.animating && target.HasComponent(component.PlayerTag) {
                    if this.runMode {
                        //ghost is killed
                        this.actions.TriggerMap[component.Destroy_actionid] = true
                    } else {
                        //player is killed
                        targetActions := *component.Actions.Get(target)
                        targetActions.TriggerMap[component.Destroy_actionid] = true
                    }
                }
            }
        }
    })
    donburi.SetValue(this.entry, component.Actions, this.actions)

    changeAiMode := func(w donburi.World, event event.AiMode) {
        if this.fastMode {
            return
        }
        this.aiMode = AiMode(event.Value)
        if !this.runMode {
            this.reverse()
        }
    }
    event.SetAiModeEvent.Subscribe(this.ecs.World, changeAiMode)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.SetAiModeEvent.Unsubscribe(this.ecs.World, changeAiMode)
            },
        },
    )

    elroy := func(w donburi.World, event event.ElroyMode) {
        if this.varient != ClassicRed_ghostvarient {
            return
        }
        println("ELROY")
        this.fastMode = true
        this.aiMode = Chase_aimode 
    }
    event.ElroyModeEvent.Subscribe(this.ecs.World, elroy)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.ElroyModeEvent.Unsubscribe(this.ecs.World, elroy)
            },
        },
    )

    frighten := func(w donburi.World, event event.RunMode) {
        ft := ghostFrightenTime
        if wave < 5 {
            ft -= int(float64(ft) * 0.1 * float64(wave-1))
        } else {
            ft = ft/2
        }
        println("fright val ", ft)
        this.runTimer = ft
        this.runMode = true
        this.reverse()
    }
    event.SetRunModeEvent.Subscribe(this.ecs.World, frighten)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.SetRunModeEvent.Unsubscribe(this.ecs.World, frighten)
            },
        },
    )

    despawn := func(w donburi.World, event event.DespawnAllEnemies) {
        this.actions.TriggerMap[component.DestroySilent_actionid] = true
    }
    event.DespawnAllEnemiesEvent.Subscribe(this.ecs.World, despawn)
    event.RegisterCleanupFuncEvent.Publish(
        this.ecs.World, 
        event.RegisterCleanupFunc{
            Function: func() {
                event.DespawnAllEnemiesEvent.Unsubscribe(this.ecs.World, despawn)
            },
        },
    )
 
    return this
}


