package astra_entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/yohamta/donburi/features/math"
    "math/rand"
    gMath "math"
    "strconv"
)

var frameStart int //used for syncing up animations

const (
    AlienConvoySpeed = 0.1
    AlienConvoySpacing = 14.0
    AlienConvoySendInitCd = 100
    AlienConvoySendCd = 600

    AlienChargeSpeed = 0.5
    AlienTurnSpeed = 0.01
    AlienReturnSpeed = 0.5
    AlienDamage = 1.0
    AlienShootCoolDown = 80
    AlienShootDelay = 300
    AlienBulletSpeed = 1.3
    AlienHealth = 1.0
    AlienHitRadius = 4
    AlienHitOffsetX = 0
    AlienHitOffsetY = 0
    BlueAlienPointValue = 30
    AlienAimOffsetY = -15.0
)

type AlienType int

const (
    Undefined_alientype AlienType = iota
    Blue_alientype 
    Purple_alientype
    Green_alientype
    Grey_alientype
)

// Alien type to the name of the sprite to use for that type
var alienSpriteMap = map[AlienType]string {
    Blue_alientype: "AlienA",
    Purple_alientype: "AlienB",
    Green_alientype: "AlienC",
    Grey_alientype: "AlienD",
}

// Alien type to base point value for destroying alien
var alienPointMap = map[AlienType]int {
    Blue_alientype: 30,
    Purple_alientype: 40,
    Green_alientype: 50,
    Grey_alientype: 60,
}

type alienData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity

    factions component.FactionsData
    damage component.DamageData
    health component.HealthData
    collider component.ColliderData
    position component.PositionData
    view component.ViewData
    velocity component.VelocityData
    graphicObject component.GraphicObjectData
    actions component.ActionsData

    pointValue int
    aType AlienType
    returnY float64
    playerPos *component.PositionData
    boss *donburi.Entry
    curAngle float64
    strafeVal float64
    strafeSet bool
}

// Move the ship, if strafe is set it will use input strafe mod value
func (this *alienData) move(strafeMod float64) {
    if this.strafeSet {
        this.position.Point.X += strafeMod
        this.position.Point.Y += AlienChargeSpeed/ 1.5
    } else {
        moveVect := math.Vec2{X:0, Y:AlienChargeSpeed}
        moveVect = moveVect.Rotate(this.curAngle)
        this.position.Point.X += moveVect.X
        this.position.Point.Y += moveVect.Y
    }
}

// Calculates sin trajectory stuff, outputs strafemod for use with move
func (this *alienData) strafe() float64 {
    strafeMod := 0.0
    if this.position.Point.Y > this.returnY {
        if !this.strafeSet {
            this.strafeVal = (gMath.Pi / this.view.View.Area.Max.X)*(this.position.Point.X)
            this.strafeSet = true
        }
        strafeMod = gMath.Sin(this.strafeVal) / 2.0
        this.strafeVal += 0.005
    }
    return strafeMod
}

// Turn alien towards player
func (this *alienData) angleTowardsPlayer() {
    // find angle of player ship
    angleRad := gMath.Atan2(this.position.Point.X - this.playerPos.Point.X, (this.playerPos.Point.Y + AlienAimOffsetY) - this.position.Point.Y)

    // Turn towards the point we're aiming at
    a := this.curAngle - angleRad
    a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
    if a <= 0 {
        this.curAngle += AlienTurnSpeed
    } else {
        this.curAngle -= AlienTurnSpeed
    }

    // clean up if we looped around
    if this.curAngle >= (2 * gMath.Pi) {
        this.curAngle -= (2 * gMath.Pi)
    }
}

// Sets animation frame to match the current rotation
func (this *alienData) rotateAnimation() {
    // Determine the closest animation frame for our angle and load that
    angle := 90 + int(this.curAngle * (180.0/gMath.Pi))
    floor := (angle/15) * 15
    ceiling := floor + 15
    if gMath.Abs(float64(angle - floor)) > gMath.Abs(float64(angle - ceiling)) {
        angle = ceiling
    } else {
        angle = floor
    }
    if angle < 0 {
        angle += 360
    }
    angle = angle % 360 // make sure we're in the range of our tags
    for i := 0; i < len(this.graphicObject.Renderables); i++ {
        this.graphicObject.Renderables[i].Play(strconv.Itoa(angle)) //tags are labeled based on 15 degree incrments
    }
}

// Do all the things that must happen when the ship is destroyed
func (this *alienData) destroy() {
    // Hide so it disappears right away
    this.graphicObject.HideAllRenderables(true)

    AddExplosion(this.ecs, 
                 this.position.Point.X, 
                 this.position.Point.Y, 
                 alienSpriteMap[this.aType], 
                 this.view.View)

    pointVal := this.pointValue
    if this.actions.TriggerMap[component.Charge_actionid] ||
       this.actions.TriggerMap[component.Follow_actionid] {
        pointVal += 10 // bonus points for hitting chargers
    }
    se := event.Score{Value:pointVal}
    event.ScoreEvent.Publish(this.ecs.World, se)

    rfe := event.RemoveFromFormation{Entry:this.entry}
    event.RemoveFromFormationEvent.Publish(this.ecs.World, rfe)

    ree := event.RemoveEntity{Entity:this.entity}
    event.RemoveEntityEvent.Publish(this.ecs.World, ree)
}

func (this *alienData) prepareReturnToConvoy() {
    this.actions.TriggerMap[component.Follow_actionid] = false
    this.actions.TriggerMap[component.Charge_actionid] = false
    this.actions.TriggerMap[component.ReturnShip_actionid] = true
    this.curAngle = gMath.Pi
    this.position.Point.Y = this.view.View.Area.Min.Y - 30
    this.graphicObject.PlayAllRenderables("90")
    this.strafeSet = false
}

//TODO next time we do something like this, clean up this stuff
// Return alien to it's convoy position
func (this *alienData) returnToConvoy() {
    this.actions.TriggerMap[component.Shoot_actionid] = false

    // Orient the alien for returning
    increments := int(this.returnY - this.position.Point.Y)
    if increments <= 24 {
        tag := strconv.Itoa(270 - (15*(increments/2)))
        this.graphicObject.PlayAllRenderables(tag)
    }

    // Move the alien to it's correct position
    if this.position.Point.Y + AlienReturnSpeed >= this.returnY {
        this.position.Point.Y = this.returnY
        this.graphicObject.PlayAllRenderables("Idle")
        this.actions.TriggerMap[component.ReturnShip_actionid] = false
    } else {
        this.position.Point.Y += AlienReturnSpeed
    }
}

func (this *alienData) shoot() {
    cd := component.Cooldown{Cur:AlienShootCoolDown, Max:AlienShootCoolDown}
    this.actions.CooldownMap[component.Shoot_actionid] = cd

    bulletVelocity := math.Vec2{X:0, Y:AlienBulletSpeed}
    AddAlienBullet(
        this.ecs, 
        this.position.Point.X, 
        this.position.Point.Y + 4, 
        bulletVelocity, 
        this.view.View)
}

func AddAlien( ecs *ecs.ECS,
               x, y float64,
               view *utility.View, 
               playerPos *component.PositionData,
               aType AlienType,
               boss *donburi.Entry,
               wave int) *donburi.Entity {
    ad := &alienData{}
    ad.ecs = ecs

    entity := ad.ecs.Create(
        layer.Foreground,
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Velocity,
        component.Collider,
        component.Health,
        component.Factions,
        component.Actions,
        component.Damage,
    )
    ad.entity = &entity

    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:ad.entity})

    ad.entry = ecs.World.Entry(*ad.entity)

    // Factions
    factions := []component.FactionId{component.Enemy_factionid}
    ad.factions = component.FactionsData{Values: factions}
    donburi.SetValue(ad.entry, component.Factions, ad.factions)

    // Damage
    ad.damage = component.NewDamageData(AlienDamage)
    donburi.SetValue(ad.entry, component.Damage, ad.damage)

    // Health
    healthAmount := AlienHealth
    ad.health = component.HealthData{Value: &healthAmount}
    donburi.SetValue(ad.entry, component.Health, ad.health)

    // Collider
    ad.collider = component.NewColliderData()
    hb := component.NewHitbox(AlienHitRadius, AlienHitOffsetX, AlienHitOffsetY)
    ad.collider.Hitboxes = append(ad.collider.Hitboxes, hb)
    donburi.SetValue(ad.entry, component.Collider, ad.collider)

    // Position
    ad.position = component.NewPositionData(x, y)
    donburi.SetValue(ad.entry, component.Position, ad.position)

    ad.returnY = ad.position.Point.Y //used to return to the original position
    ad.playerPos = playerPos //so the AI can track the player ship
    ad.boss = boss //for when the alien has a boss to follow

    // View
    ad.view = component.ViewData{View:view}
    donburi.SetValue(ad.entry, component.View, ad.view)

    // Velocity
    ad.velocity = component.VelocityData{Velocity: &math.Vec2{}}
    donburi.SetValue(ad.entry, component.Velocity, ad.velocity)

    // Action
    ad.actions = component.NewActions()
    ad.actions.AddNormalAction(component.Destroy_actionid, ad.destroy)
    ad.actions.AddNormalAction(component.ReturnShip_actionid, ad.returnToConvoy)
    ad.actions.AddNormalAction(component.Shoot_actionid, ad.shoot)
    donburi.SetValue(ad.entry, component.Actions, ad.actions)

    ad.aType = aType
    ad.pointValue = alienPointMap[ad.aType]// + (10 * (wave - 1))

    //TODO need to make a slice of pointers?
    ad.graphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(alienSpriteMap[ad.aType], nil)
    nsd.Play("Idle")
    nsd.SetFrame(frameStart)
    frameStart = (frameStart + rand.Intn(3)) % 10
    ad.graphicObject.Renderables = append(ad.graphicObject.Renderables, &nsd)
    donburi.SetValue(ad.entry, component.GraphicObject, ad.graphicObject)

    ad.curAngle = gMath.Pi

    ad.init()

    return ad.entity
}

func (this *alienData) init() {
    switch this.aType {
    case Blue_alientype:
        this.initBlue()
    case Purple_alientype:
        this.initPurple()
    case Green_alientype:
        this.initGreen()
    case Grey_alientype:
        this.initGrey()
    }
}

func (this *alienData) initBlue() {
    this.actions.AddNormalAction(component.Charge_actionid, func() {
        this.actions.TriggerMap[component.Shoot_actionid] = true

        if this.position.Point.Y < (this.playerPos.Point.Y + AlienAimOffsetY) {
            this.angleTowardsPlayer()
        } 

        this.rotateAnimation()

        this.move(0.0)

        if this.position.Point.Y > (this.view.View.Area.Max.Y + 30) {
            this.prepareReturnToConvoy()
        }
    })
}

func (this *alienData) initPurple() {
    this.actions.AddNormalAction(component.Charge_actionid, func() {
        this.actions.TriggerMap[component.Shoot_actionid] = true

        if this.position.Point.Y < (this.playerPos.Point.Y + AlienAimOffsetY) {
            this.angleTowardsPlayer()
        } 

        this.rotateAnimation()

        strafeMod := this.strafe()

        this.move(strafeMod)

        if this.position.Point.Y > (this.view.View.Area.Max.Y + 30) {
            this.prepareReturnToConvoy()
        }
    })
}

func (this *alienData) initGreen() {
    this.actions.AddNormalAction(component.Charge_actionid, func() {
        // Don't charge if our boss is still alive.
        if this.boss.Valid() {
            this.actions.TriggerMap[component.Charge_actionid] = false
            return
        } 
        
        this.actions.TriggerMap[component.Shoot_actionid] = true

        if this.position.Point.Y < (this.playerPos.Point.Y + AlienAimOffsetY) {
            this.angleTowardsPlayer()
        } 

        this.rotateAnimation()

        strafeMod := this.strafe()

        this.move(strafeMod)

        if this.position.Point.Y > (this.view.View.Area.Max.Y + 30) {
            this.prepareReturnToConvoy()
        }
    })

    bossOffsetX := 0.0
    bossOffsetY := 0.0
    if this.boss.Valid() {
		bossPos := component.Position.Get(this.boss).Point
        bossOffsetX = this.position.Point.X - bossPos.X
        bossOffsetY = this.position.Point.Y - bossPos.Y
    }
    this.actions.AddNormalAction(component.Follow_actionid, func() {
        if !this.boss.Valid() {
            println("boss died")
            this.actions.TriggerMap[component.Charge_actionid] = true
            this.actions.TriggerMap[component.Follow_actionid] = false
            return
        } 

        this.actions.TriggerMap[component.Shoot_actionid] = true

        bossPos := component.Position.Get(this.boss).Point
        this.position.Point.X = bossPos.X + bossOffsetX
        this.position.Point.Y = bossPos.Y + bossOffsetY


        if this.position.Point.Y < (this.playerPos.Point.Y + AlienAimOffsetY) {
            this.angleTowardsPlayer()
        } 

        this.rotateAnimation()
    })

    this.actions.AddUpkeepAction(func() {
        if this.boss.Valid() {
            //if the boss is charging bodyguards should follow
            acts := component.Actions.Get(this.boss)
            if acts.TriggerMap[component.Charge_actionid] &&
               !this.actions.TriggerMap[component.Follow_actionid] {
                this.actions.TriggerMap[component.Follow_actionid] = true
                this.actions.CooldownMap[component.Shoot_actionid] = component.Cooldown{
                    Cur:AlienShootDelay, 
                    Max:AlienShootDelay,
                }
            }

            //If the boss is returning so should the bodyguards 
            if acts.TriggerMap[component.ReturnShip_actionid] {
                this.prepareReturnToConvoy()
            }
        }
    })
}

func (this *alienData) initGrey() {
    this.actions.AddNormalAction(component.Charge_actionid, func() {
        this.actions.TriggerMap[component.Shoot_actionid] = true

        if this.position.Point.Y < (this.playerPos.Point.Y + AlienAimOffsetY) {
            this.angleTowardsPlayer()
        } 

        this.rotateAnimation()

        strafeMod := this.strafe()

        this.move(strafeMod)

        if this.position.Point.Y > (this.view.View.Area.Max.Y + 30) {
            this.prepareReturnToConvoy()
        }
    })
}
