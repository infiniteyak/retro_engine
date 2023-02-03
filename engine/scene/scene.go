package scene

import (
    "github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi"
)

type SceneEventId int

type SceneId int

type Scene struct {
    sceneId SceneId
    entities []*donburi.Entity
    cleanupFuncs []func() //These will be called before entities are removed
                          //Use them for unsubscribing from callbacks etc
    ecs *ecs.ECS
}

func NewScene(ecs *ecs.ECS) *Scene {
    scene := &Scene{}
    scene.entities = make([]*donburi.Entity, 0)
    scene.cleanupFuncs = make([]func(), 0)
    scene.ecs = ecs

    // Event to handle adding entities
    registerEntityFunc := func(w donburi.World, event event.RegisterEntity) {
        scene.entities = append(scene.entities, event.Entity)
    }
    event.RegisterEntityEvent.Subscribe(ecs.World, registerEntityFunc)
    scene.cleanupFuncs = append(scene.cleanupFuncs, func() {
        event.RegisterEntityEvent.Unsubscribe(ecs.World, registerEntityFunc)
    })

    // Event to handle removing entities
    removeEntityFunc := func(w donburi.World, event event.RemoveEntity) {
        for i, e := range scene.entities {
            if e == event.Entity {
                scene.entities[i] = scene.entities[len(scene.entities)-1]
                scene.entities = scene.entities[:len(scene.entities)-1]
                w.Remove(*e)
                break
            }
        }
    }
    event.RemoveEntityEvent.Subscribe(ecs.World, removeEntityFunc)
    scene.cleanupFuncs = append(scene.cleanupFuncs, func() {
        event.RemoveEntityEvent.Unsubscribe(ecs.World, removeEntityFunc)
    })

    // Event to handle adding cleanup functions which are called when the scene ends
    registerCleanupFunc := func(w donburi.World, event event.RegisterCleanupFunc) {
        scene.cleanupFuncs = append(scene.cleanupFuncs, event.Function)
    }
    event.RegisterCleanupFuncEvent.Subscribe(ecs.World, registerCleanupFunc)
    scene.cleanupFuncs = append(scene.cleanupFuncs, func() {
        event.RegisterCleanupFuncEvent.Unsubscribe(ecs.World, registerCleanupFunc)
    })

    return scene
}

func (this *Scene) GetId() SceneId {
    return this.sceneId
}

func (this *Scene) SetId(id SceneId) {
    this.sceneId = id
}

func (this *Scene) Cleanup() {
    for _, foo := range this.cleanupFuncs {
        foo()
    }
    for _, e := range this.entities {
        this.ecs.World.Remove(*e)
    }
}
