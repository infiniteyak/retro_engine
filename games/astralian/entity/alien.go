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
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var frameStart int //used for syncing up animations

const (
    AlienConvoySpeed = 0.1
    AlienChargeSpeed = 0.5
    AlienTurnSpeed = 0.01
    AlienReturnSpeed = 0.5
    AlienDamage = 1.0
    AlienShootCoolDown = 100
    AlienShootDelay = 300
    AlienBulletSpeed = 1.3
    AlienHealth = 1.0
    AlienHitRadius = 4
    AlienHitOffsetX = 0
    AlienHitOffsetY = 0
    BlueAlienPointValue = 30
)

type AlienType int

const (
    Undefined_alientype AlienType = iota
    Blue_alientype 
    Purple_alientype
    Green_alientype
    Grey_alientype
)

// Alien type to the function that initializes that type
var alienInitMap = map[AlienType]func(ad *alienData) {
    Blue_alientype: initBlueAlien,
    Purple_alientype: initPurpleAlien,
    Green_alientype: initGreenAlien,
    Grey_alientype: initGreyAlien,
}

// Alien type to the name of the sprite to use for that type
var alienSpriteMap = map[AlienType]string {
    Blue_alientype: "AlienA",
    Purple_alientype: "AlienB",
    Green_alientype: "AlienC",
    Grey_alientype: "AlienD",
}

//TODO maybe these should be lowercase? (and above)
type alienData struct {
    Ecs *ecs.ECS
    Entry *donburi.Entry // You can get this from ecs and entity, but easier this way
    Entity *donburi.Entity
    AudioContext *audio.Context

    Factions component.FactionsData
    Damage component.DamageData
    Health component.HealthData
    Collider component.ColliderData
    Position component.PositionData
    View component.ViewData
    Velocity component.VelocityData
    GraphicObject component.GraphicObjectData
    Actions component.ActionsData

    PointValue int
    Type AlienType
    ReturnY float64
    PlayerPos *component.PositionData
    Boss *donburi.Entry
}

// Do all the things that must happen when the ship is destroyed
func (this *alienData) destroy() {
    // Hide so it disappears right away
    this.GraphicObject.HideAllRenderables(true)

    AddExplosion(this.Ecs, 
                 this.Position.Point.X, 
                 this.Position.Point.Y, 
                 alienSpriteMap[this.Type], 
                 this.View.View)

    se := event.Score{Value:this.PointValue}
    event.ScoreEvent.Publish(this.Ecs.World, se)

    rfe := event.RemoveFromFormation{Entry:this.Entry}
    event.RemoveFromFormationEvent.Publish(this.Ecs.World, rfe)

    ree := event.RemoveEntity{Entity:this.Entity}
    event.RemoveEntityEvent.Publish( this.Ecs.World, ree)
}

// Return alien to it's convoy position
func (this *alienData) returnToConvoy() {
    this.Actions.TriggerMap[component.Shoot_actionid] = false

    //TODO common code and consts for this section
    // Orient the alien for returning
    increments := int(this.ReturnY - this.Position.Point.Y)
    if increments <= 24 {
        tag := strconv.Itoa(270 - (15*(increments/2)))
        this.GraphicObject.PlayAllRenderables(tag)
    }

    // Move the alien to it's correct position
    if this.Position.Point.Y + AlienReturnSpeed >= this.ReturnY {
        this.Position.Point.Y = this.ReturnY
        this.GraphicObject.PlayAllRenderables("Idle")
        this.Actions.TriggerMap[component.ReturnShip_actionid] = false
    } else {
        this.Position.Point.Y += AlienReturnSpeed
    }
}

func (this *alienData) shoot() {
    //TODO think about cooldowns
    cd := component.Cooldown{Cur:AlienShootCoolDown, Max:AlienShootCoolDown}
    this.Actions.CooldownMap[component.Shoot_actionid] = cd

    bulletVelocity := math.Vec2{X:0, Y:AlienBulletSpeed}
    AddAlienBullet( //TODO interface? bullet spawn offset?
        this.Ecs, 
        this.Position.Point.X, 
        this.Position.Point.Y + 4, 
        bulletVelocity, 
        this.View.View, 
        this.AudioContext)
}

func AddAlien( ecs *ecs.ECS, //TODO can I make this global? or pass in game object? interface?
               x, y float64, //TODO can I make this better? and standard?
               view *utility.View, 
               audioContext *audio.Context, //TODO could/should this be global?
               playerPos *component.PositionData, //TODO improve position data
               aType AlienType,
               boss *donburi.Entry ) *donburi.Entity { //TODO does anything actually use return?
    ad := &alienData{}
    ad.Ecs = ecs

    entity := ad.Ecs.Create(
        layer.Foreground, // TODO argument?
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
    ad.Entity = &entity

    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:ad.Entity})

    ad.Entry = ecs.World.Entry(*ad.Entity)

    ad.AudioContext = audioContext

    // Factions
    factions := []component.FactionId{component.Enemy_factionid}
    ad.Factions = component.FactionsData{Values: factions}
    donburi.SetValue(ad.Entry, component.Factions, ad.Factions)

    // Damage
    damageAmount := AlienDamage
    ad.Damage = component.DamageData{Value: &damageAmount}
    donburi.SetValue(ad.Entry, component.Damage, ad.Damage)

    // Health
    healthAmount := AlienHealth
    ad.Health = component.HealthData{Value: &healthAmount}
    donburi.SetValue(ad.Entry, component.Health, ad.Health)

    // Collider
    ad.Collider = component.NewColliderData()
    hb := component.NewHitbox(AlienHitRadius, AlienHitOffsetX, AlienHitOffsetY)
    ad.Collider.Hitboxes = append(ad.Collider.Hitboxes, hb)
    donburi.SetValue(ad.Entry, component.Collider, ad.Collider)

    // Position
    ad.Position = component.NewPositionData(x, y)
    donburi.SetValue(ad.Entry, component.Position, ad.Position)

    //TODO is this the best way?
    ad.ReturnY = ad.Position.Point.Y //used to return to the original position
    ad.PlayerPos = playerPos //so the AI can track the player ship
    ad.Boss = boss //for when the alien has a boss to follow

    // View
    ad.View = component.ViewData{View:view}
    donburi.SetValue(ad.Entry, component.View, ad.View)

    // Velocity
    ad.Velocity = component.VelocityData{Velocity: &math.Vec2{}}
    donburi.SetValue(ad.Entry, component.Velocity, ad.Velocity)

    // Graphic Object
    /*
    ad.GraphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(sprite, nil)
    nsd.Play("Idle")
    nsd.SetFrame(frameStart)
    frameStart = (frameStart + rand.Intn(3)) % 10
    ad.GraphicObject.Renderables = append(ad.GraphicObject.Renderables, &nsd)
    donburi.SetValue(ad.Entry, component.GraphicObject, ad.GraphicObject)
    */

    // Action
    ad.Actions = component.NewActions()
    ad.Actions.ActionMap[component.Destroy_actionid] = ad.destroy
    ad.Actions.ActionMap[component.ReturnShip_actionid] = ad.returnToConvoy
    ad.Actions.ActionMap[component.Shoot_actionid] = ad.shoot
    donburi.SetValue(ad.Entry, component.Actions, ad.Actions)

    /*
    ad.Type = aType
    ad.init()
    */

    alienInitMap[aType](ad)
    //initBlueAlien(ad) //TODO add functionality to choose type
    //initPurpleAlien(ad) //TODO add functionality to choose type

    return ad.Entity
}

/*
func (this *alienData) init() {
    switch this.Type {
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
*/

func initBlueAlien(ad *alienData) {
    //TODO okay this kind of works, but for some reason I can't modify the 
    // do the donburi.SetValue above and then modify the slice later. I think
    // I need to make renderables a pointer to a slice or something? Probably
    // all the component data types need to have their members as pointers so
    // we can manipulate them like this... for now I'm just doing it here...
    spriteName := "AlienA"
    ad.GraphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(spriteName, nil)
    nsd.Play("Idle")
    nsd.SetFrame(frameStart)
    frameStart = (frameStart + rand.Intn(3)) % 10
    ad.GraphicObject.Renderables = append(ad.GraphicObject.Renderables, &nsd)
    donburi.SetValue(ad.Entry, component.GraphicObject, ad.GraphicObject)

    ad.PointValue = BlueAlienPointValue
    ad.Type = Blue_alientype

    curAngle := gMath.Pi
    //strafeVal := 0.0
    strafeSet := false
    ad.Actions.ActionMap[component.Charge_actionid] = func() {
        ad.Actions.TriggerMap[component.Shoot_actionid] = true
        aimOffsetY := -15.0 // Aim slightly above the ship

        angleRad := 0.0
        if ad.Position.Point.Y < (ad.PlayerPos.Point.Y + aimOffsetY) {
            // angle towards player ship
            angleRad = gMath.Atan2(ad.Position.Point.X - ad.PlayerPos.Point.X, (ad.PlayerPos.Point.Y + aimOffsetY) - ad.Position.Point.Y)

            // Turn towards the point we're aiming at
            a := curAngle - angleRad
            a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
            if a <= 0 {
                curAngle += AlienTurnSpeed
            } else {
                curAngle -= AlienTurnSpeed
            }
        } 

        // Use move rotation and charge speed to create a vector for movement
        moveVect := math.Vec2{X:0, Y:AlienChargeSpeed}
        moveVect = moveVect.Rotate(curAngle)

        // Determine the closest animation frame for our angle and load that
        angleDeg := curAngle * (180.0/gMath.Pi)
        //angleDeg := angleRad * (180.0/gMath.Pi)
        angle := 90 + int(angleDeg) // TODO make this better
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
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play(strconv.Itoa(angle)) //tags are labeled based on 15 degree incrments
        }

        // Strafe mod
        strafeMod := 0.0
        /*
        if moveVect.Y > 0 {
            if !strafeSet {
                //strafeVal = (2 * gMath.Pi / view.Area.Max.X)*(pd.Point.X)-(gMath.Pi)
                strafeVal = (gMath.Pi / view.Area.Max.X)*(pd.Point.X)
                strafeSet = true
            }
            strafeMod = gMath.Sin(strafeVal) * 2
            strafeVal += 0.03
        }
        */

        // move the ship
        if strafeSet {
            ad.Position.Point.X += strafeMod
            ad.Position.Point.Y += AlienChargeSpeed/2
        } else {
            ad.Position.Point.X += moveVect.X
            ad.Position.Point.Y += moveVect.Y
        }

        //pd.Point.X += strafeMod + moveVect.X
        //pd.Point.Y += AlienChargeSpeed/4 + moveVect.Y


        // clean up the current angle
        if curAngle >= (2 * gMath.Pi) {
            curAngle -= (2 * gMath.Pi)
        }

        if ad.Position.Point.Y > (ad.View.View.Area.Max.Y + 30) {
            ad.Actions.TriggerMap[component.Charge_actionid] = false
            ad.Actions.TriggerMap[component.ReturnShip_actionid] = true

            //Set values for return
            curAngle = gMath.Pi
            ad.Position.Point.Y = ad.View.View.Area.Min.Y - 30
            for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
                ad.GraphicObject.Renderables[i].Play("90") //tags are labeled based on 15 degree incrments
            }
        }
    }
}

func initPurpleAlien(ad *alienData) {
    spriteName := "AlienB"
    ad.GraphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(spriteName, nil)
    nsd.Play("Idle")
    nsd.SetFrame(frameStart)
    frameStart = (frameStart + rand.Intn(3)) % 10
    ad.GraphicObject.Renderables = append(ad.GraphicObject.Renderables, &nsd)
    donburi.SetValue(ad.Entry, component.GraphicObject, ad.GraphicObject)

    ad.PointValue = BlueAlienPointValue
    ad.Type = Purple_alientype

    curAngle := gMath.Pi
    strafeVal := 0.0
    strafeSet := false
    targetX := ad.PlayerPos.Point.X
    targetY := ad.PlayerPos.Point.Y
    ad.Actions.ActionMap[component.Charge_actionid] = func() {
        ad.Actions.TriggerMap[component.Shoot_actionid] = true
        aimOffsetY := -15.0 // Aim slightly above the ship (for movement)

        angleRad := 0.0
        if ad.Position.Point.Y < (ad.PlayerPos.Point.Y + aimOffsetY) {
            // angle towards player ship
            angleRad = gMath.Atan2(ad.Position.Point.X - targetX, (targetY + aimOffsetY) - ad.Position.Point.Y)

            // Turn towards the point we're aiming at
            a := curAngle - angleRad
            a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
            if a <= 0 {
                curAngle += AlienTurnSpeed
            } else {
                curAngle -= AlienTurnSpeed
            }
        } 

        // Use move rotation and charge speed to create a vector for movement
        moveVect := math.Vec2{X:0, Y:AlienChargeSpeed}
        moveVect = moveVect.Rotate(curAngle)

        // Determine the closest animation frame for our angle and load that
        angle := 90 + int(curAngle * (180.0/gMath.Pi))
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
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play(strconv.Itoa(angle)) //tags are labeled based on 15 degree incrments
        }

        // Strafe mod
        strafeMod := 0.0
        if ad.Position.Point.Y > ad.ReturnY {
            if !strafeSet {
                strafeVal = (gMath.Pi / ad.View.View.Area.Max.X)*(ad.Position.Point.X)
                strafeSet = true
            }
            strafeMod = gMath.Sin(strafeVal) / 2.0
            strafeVal += 0.005
        }

        // move the ship
        if strafeSet {
            ad.Position.Point.X += strafeMod
            ad.Position.Point.Y += AlienChargeSpeed/ 1.5
        } else {
            ad.Position.Point.X += moveVect.X
            ad.Position.Point.Y += moveVect.Y
        }

        // clean up the current angle
        if curAngle >= (2 * gMath.Pi) {
            curAngle -= (2 * gMath.Pi)
        }

        if ad.Position.Point.Y > (ad.View.View.Area.Max.Y + 30) {
            ad.Actions.TriggerMap[component.Charge_actionid] = false
            ad.Actions.TriggerMap[component.ReturnShip_actionid] = true

            //Set values for return
            strafeSet = false
            curAngle = gMath.Pi
            ad.Position.Point.Y = ad.View.View.Area.Min.Y - 30
            for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
                ad.GraphicObject.Renderables[i].Play("90") //tags are labeled based on 15 degree incrments
            }
        }
    }
}

func initGreenAlien(ad *alienData) {
    spriteName := "AlienC"
    ad.GraphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(spriteName, nil)
    nsd.Play("Idle")
    nsd.SetFrame(frameStart)
    frameStart = (frameStart + rand.Intn(3)) % 10
    ad.GraphicObject.Renderables = append(ad.GraphicObject.Renderables, &nsd)
    donburi.SetValue(ad.Entry, component.GraphicObject, ad.GraphicObject)

    ad.PointValue = BlueAlienPointValue
    ad.Type = Purple_alientype

    curAngle := gMath.Pi
    strafeVal := 0.0
    strafeSet := false
    targetX := ad.PlayerPos.Point.X //TODO wait does this just happen the first time?!
    targetY := ad.PlayerPos.Point.Y
    ad.Actions.ActionMap[component.Charge_actionid] = func() {
        // Don't charge if our boss is still alive.
        if ad.Boss.Valid() {
            ad.Actions.TriggerMap[component.Charge_actionid] = false
            return
        } 
        
        ad.Actions.TriggerMap[component.Shoot_actionid] = true
        aimOffsetY := -15.0 // Aim slightly above the ship (for movement)

        angleRad := 0.0
        if ad.Position.Point.Y < (ad.PlayerPos.Point.Y + aimOffsetY) {
            // angle towards player ship
            angleRad = gMath.Atan2(ad.Position.Point.X - targetX, (targetY + aimOffsetY) - ad.Position.Point.Y)

            // Turn towards the point we're aiming at
            a := curAngle - angleRad
            a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
            if a <= 0 {
                curAngle += AlienTurnSpeed
            } else {
                curAngle -= AlienTurnSpeed
            }
        } 

        // Use move rotation and charge speed to create a vector for movement
        moveVect := math.Vec2{X:0, Y:AlienChargeSpeed}
        moveVect = moveVect.Rotate(curAngle)

        // Determine the closest animation frame for our angle and load that
        angle := 90 + int(curAngle * (180.0/gMath.Pi))
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
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play(strconv.Itoa(angle)) //tags are labeled based on 15 degree incrments
        }

        // Strafe mod
        strafeMod := 0.0
        if ad.Position.Point.Y > ad.ReturnY {
            if !strafeSet {
                strafeVal = (gMath.Pi / ad.View.View.Area.Max.X)*(ad.Position.Point.X)
                strafeSet = true
            }
            strafeMod = gMath.Sin(strafeVal) / 2.0
            strafeVal += 0.005
        }

        // move the ship
        if strafeSet {
            ad.Position.Point.X += strafeMod
            ad.Position.Point.Y += AlienChargeSpeed/ 1.5
        } else {
            ad.Position.Point.X += moveVect.X
            ad.Position.Point.Y += moveVect.Y
        }

        // clean up the current angle
        if curAngle >= (2 * gMath.Pi) {
            curAngle -= (2 * gMath.Pi)
        }

        if ad.Position.Point.Y > (ad.View.View.Area.Max.Y + 30) {
            ad.Actions.TriggerMap[component.Charge_actionid] = false
            ad.Actions.TriggerMap[component.ReturnShip_actionid] = true

            //Set values for return
            strafeSet = false
            curAngle = gMath.Pi
            ad.Position.Point.Y = ad.View.View.Area.Min.Y - 30
            for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
                ad.GraphicObject.Renderables[i].Play("90") //tags are labeled based on 15 degree incrments
            }
        }
    }

    //determine boss offset
    bossOffsetX := 0.0
    bossOffsetY := 0.0
    if ad.Boss.Valid() {
		bossPos := component.Position.Get(ad.Boss).Point
        bossOffsetX = ad.Position.Point.X - bossPos.X
        bossOffsetY = ad.Position.Point.Y - bossPos.Y
    }
    ad.Actions.ActionMap[component.Follow_actionid] = func() {
        if !ad.Boss.Valid() {
            // TODO switch into some other mode? and return?
        } 
        aimOffsetY := -15.0 // Aim slightly above the ship (for movement)
        bossPos := component.Position.Get(ad.Boss).Point
        ad.Position.Point.X = bossPos.X + bossOffsetX
        ad.Position.Point.Y = bossPos.Y + bossOffsetY

        ad.Actions.TriggerMap[component.Shoot_actionid] = true

        angleRad := 0.0
        if ad.Position.Point.Y < (ad.PlayerPos.Point.Y + aimOffsetY) {
            // angle towards player ship
            angleRad = gMath.Atan2(ad.Position.Point.X - targetX, (targetY + aimOffsetY) - ad.Position.Point.Y)

            // Turn towards the point we're aiming at
            a := curAngle - angleRad
            a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
            if a <= 0 {
                curAngle += AlienTurnSpeed
            } else {
                curAngle -= AlienTurnSpeed
            }
        } 
        // Determine the closest animation frame for our angle and load that
        angle := 90 + int(curAngle * (180.0/gMath.Pi))
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
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play(strconv.Itoa(angle)) //tags are labeled based on 15 degree incrments
        }
    }

    ad.Actions.ActionMap[component.Upkeep_actionid] = func() {
        if ad.Boss.Valid() {
            acts := component.Actions.Get(ad.Boss)
            if acts.TriggerMap[component.Charge_actionid] &&
               !ad.Actions.TriggerMap[component.Follow_actionid] {
                ad.Actions.TriggerMap[component.Follow_actionid] = true
                ad.Actions.CooldownMap[component.Shoot_actionid] = component.Cooldown{
                    Cur:AlienShootDelay, 
                    Max:AlienShootDelay,
                }
            }
            if acts.TriggerMap[component.ReturnShip_actionid] {
                ad.Actions.TriggerMap[component.Follow_actionid] = false
                ad.Actions.TriggerMap[component.ReturnShip_actionid] = true

                //Set values for return
                ad.Position.Point.Y = ad.View.View.Area.Min.Y - 30
                for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
                    ad.GraphicObject.Renderables[i].Play("90") //tags are labeled based on 15 degree incrments
                }
            }
        }
    }
}

func initGreyAlien(ad *alienData) {
    spriteName := "AlienD"
    ad.GraphicObject = component.NewGraphicObjectData()
    nsd := component.SpriteData{}
    nsd.Load(spriteName, nil)
    nsd.Play("Idle")
    nsd.SetFrame(frameStart)
    frameStart = (frameStart + rand.Intn(3)) % 10
    ad.GraphicObject.Renderables = append(ad.GraphicObject.Renderables, &nsd)
    donburi.SetValue(ad.Entry, component.GraphicObject, ad.GraphicObject)

    ad.PointValue = BlueAlienPointValue
    ad.Type = Purple_alientype

    curAngle := gMath.Pi
    strafeVal := 0.0
    strafeSet := false
    targetX := ad.PlayerPos.Point.X
    targetY := ad.PlayerPos.Point.Y
    ad.Actions.ActionMap[component.Charge_actionid] = func() {
        ad.Actions.TriggerMap[component.Shoot_actionid] = true
        aimOffsetY := -15.0 // Aim slightly above the ship (for movement)

        angleRad := 0.0
        if ad.Position.Point.Y < (ad.PlayerPos.Point.Y + aimOffsetY) {
            // angle towards player ship
            angleRad = gMath.Atan2(ad.Position.Point.X - targetX, (targetY + aimOffsetY) - ad.Position.Point.Y)

            // Turn towards the point we're aiming at
            a := curAngle - angleRad
            a = gMath.Mod(a + gMath.Pi, 2 * gMath.Pi) - gMath.Pi 
            if a <= 0 {
                curAngle += AlienTurnSpeed
            } else {
                curAngle -= AlienTurnSpeed
            }
        } 

        // Use move rotation and charge speed to create a vector for movement
        moveVect := math.Vec2{X:0, Y:AlienChargeSpeed}
        moveVect = moveVect.Rotate(curAngle)

        // Determine the closest animation frame for our angle and load that
        angle := 90 + int(curAngle * (180.0/gMath.Pi))
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
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play(strconv.Itoa(angle)) //tags are labeled based on 15 degree incrments
        }

        // Strafe mod
        strafeMod := 0.0
        if ad.Position.Point.Y > ad.ReturnY {
            if !strafeSet {
                strafeVal = (gMath.Pi / ad.View.View.Area.Max.X)*(ad.Position.Point.X)
                strafeSet = true
            }
            strafeMod = gMath.Sin(strafeVal) / 2.0
            strafeVal += 0.005
        }

        // move the ship
        if strafeSet {
            ad.Position.Point.X += strafeMod
            ad.Position.Point.Y += AlienChargeSpeed/ 1.5
        } else {
            ad.Position.Point.X += moveVect.X
            ad.Position.Point.Y += moveVect.Y
        }

        // clean up the current angle
        if curAngle >= (2 * gMath.Pi) {
            curAngle -= (2 * gMath.Pi)
        }

        if ad.Position.Point.Y > (ad.View.View.Area.Max.Y + 30) {
            ad.Actions.TriggerMap[component.Charge_actionid] = false
            ad.Actions.TriggerMap[component.ReturnShip_actionid] = true

            //Set values for return
            strafeSet = false
            curAngle = gMath.Pi
            ad.Position.Point.Y = ad.View.View.Area.Min.Y - 30
            for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
                ad.GraphicObject.Renderables[i].Play("90") //tags are labeled based on 15 degree incrments
            }
        }
    }
}
