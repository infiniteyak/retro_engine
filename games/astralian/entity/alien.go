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
)

var alienInitMap = map[AlienType]func(ad *alienData) {
    Blue_alientype: initBlueAlien,
    Purple_alientype: initPurpleAlien,
}

var alienSpriteMap = map[AlienType]string {
    Blue_alientype: "AlienA",
    Purple_alientype: "AlienB",
}

//TODO maybe these should be lowercase?
type alienData struct {
    Ecs *ecs.ECS
    Entry *donburi.Entry
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
}

//TODO make this more like func (this *alienData)destroyShip() {}
func destroyShip(ad *alienData) {
    for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
        *ad.GraphicObject.Renderables[i].GetTransInfo().Hide = true
    }
    event.ScoreEvent.Publish(ad.Ecs.World, event.Score{Value:ad.PointValue})
    AddExplosion(
        ad.Ecs, ad.Position.Point.X, ad.Position.Point.Y, alienSpriteMap[ad.Type], ad.View.View)
    event.RemoveFromFormationEvent.Publish(
        ad.Ecs.World, event.RemoveFromFormation{Entry:ad.Entry})
    event.RemoveEntityEvent.Publish(
        ad.Ecs.World, 
        event.RemoveEntity{Entity:ad.Entity},
    )
}

// Handles returning the alien ship to it's original position
func returnShip(ad *alienData) {
    //strafeVal = 0.0
    //strafeSet = false //TODO if this code is common, this might  need to come out
    ad.Actions.TriggerMap[component.Shoot_actionid] = false
    increments := int(ad.ReturnY - ad.Position.Point.Y)
    if increments <= 24 {
        tag := strconv.Itoa(270 - (15*(increments/2)))
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play(tag)
        }
    }
    if ad.Position.Point.Y + AlienReturnSpeed > ad.ReturnY {
        ad.Position.Point.Y = ad.ReturnY
        for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
            ad.GraphicObject.Renderables[i].Play("Idle")
        }
        ad.Actions.TriggerMap[component.ReturnShip_actionid] = false
    } else {
        ad.Position.Point.Y += AlienReturnSpeed
    }
}

func shoot(ad *alienData) {
    bulletVelocity := math.Vec2{X:0, Y:AlienBulletSpeed}

    ad.Actions.CooldownMap[component.Shoot_actionid] = component.Cooldown{
        Cur:AlienShootCoolDown, 
        Max:AlienShootCoolDown,
    }

    // TODO Make the bullet spawn at the front of the ship, not the middle
    AddAlienBullet(
        ad.Ecs, 
        ad.Position.Point.X, 
        ad.Position.Point.Y, 
        bulletVelocity, 
        ad.View.View, 
        ad.AudioContext)
}

func AddAlien( ecs *ecs.ECS, 
               x, y float64, 
               view *utility.View, 
               audioContext *audio.Context, 
               playerPos *component.PositionData,
               alienType AlienType) *donburi.Entity {
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
    ad.Health = component.HealthData{Value:AlienHealth}
    donburi.SetValue(ad.Entry, component.Health, ad.Health)

    // Collider
    ad.Collider = component.NewColliderData()
    ad.Collider.Hitboxes = append(ad.Collider.Hitboxes, component.NewHitbox(AlienHitRadius, AlienHitOffsetX, AlienHitOffsetY))
    donburi.SetValue(ad.Entry, component.Collider, ad.Collider)

    // Position
    ad.Position = component.NewPositionData(x, y)
    donburi.SetValue(ad.Entry, component.Position, ad.Position)

    //TODO is this the best way?
    ad.ReturnY = ad.Position.Point.Y //used to return to the original position
    ad.PlayerPos = playerPos //so the AI can track the player ship

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
    ad.Actions.ActionMap[component.Destroy_actionid] = func() {destroyShip(ad)}
    ad.Actions.ActionMap[component.ReturnShip_actionid] = func() {returnShip(ad)}// TODO is this common between all aliens?
    ad.Actions.ActionMap[component.Shoot_actionid] = func() {shoot(ad)}
    donburi.SetValue(ad.Entry, component.Actions, ad.Actions)

    alienInitMap[alienType](ad)
    //initBlueAlien(ad) //TODO add functionality to choose type
    //initPurpleAlien(ad) //TODO add functionality to choose type

    return ad.Entity
}

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
        angleDeg := 0.0
        angleDeg = curAngle * (180.0/gMath.Pi)
        /*
        if moveVect.Y > 0 {
            angleDeg = angleRad * (180.0/gMath.Pi)
        } else {
            angleDeg = curAngle * (180.0/gMath.Pi)
        }
        */
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
        //if moveVect.Y > 0 {
        if ad.Position.Point.Y > ad.ReturnY {
            if !strafeSet {
                //strafeVal = (2 * gMath.Pi / view.Area.Max.X)*(pd.Point.X)-(gMath.Pi)
                strafeVal = (gMath.Pi / ad.View.View.Area.Max.X)*(ad.Position.Point.X)
                /*
                if moveVect.X > 0 {
                    strafeVal = 0
                } else {
                    strafeVal = gMath.Pi
                }
                */
                strafeSet = true
            }
            //strafeMod = gMath.Sin(strafeVal) * 2
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
            strafeSet = false
            curAngle = gMath.Pi
            ad.Position.Point.Y = ad.View.View.Area.Min.Y - 30
            for i := 0; i < len(ad.GraphicObject.Renderables); i++ {
                ad.GraphicObject.Renderables[i].Play("90") //tags are labeled based on 15 degree incrments
            }
        }
    }
}
