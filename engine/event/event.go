package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type GameOver struct {}
var GameOverEvent = events.NewEventType[GameOver]()

type ScreenClear struct {}
var ScreenClearEvent = events.NewEventType[ScreenClear]()

type RegisterCleanupFunc struct {
    Function func()
}
var RegisterCleanupFuncEvent = events.NewEventType[RegisterCleanupFunc]()

type RegisterEntity struct {
    Entity *donburi.Entity
}
var RegisterEntityEvent = events.NewEventType[RegisterEntity]()

type RemoveEntity struct {
    Entity *donburi.Entity
}
var RemoveEntityEvent = events.NewEventType[RemoveEntity]()

type Score struct {
    Value int
}
var ScoreEvent = events.NewEventType[Score]()

///// MOVE TO GAMES
type AsteroidsCountUpdate struct {
    Value int
}
var AsteroidsCountUpdateEvent = events.NewEventType[AsteroidsCountUpdate]()

type RemoveFromFormation struct {
    Entry *donburi.Entry
}
var RemoveFromFormationEvent = events.NewEventType[RemoveFromFormation]()

type ShipDestroyed struct { }
var ShipDestroyedEvent = events.NewEventType[ShipDestroyed]()
