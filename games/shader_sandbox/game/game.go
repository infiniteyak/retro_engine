package game

import (
	"os"
	"github.com/infiniteyak/retro_engine/games/shader_sandbox/asset"
	"github.com/infiniteyak/retro_engine/engine/system"

	"github.com/infiniteyak/retro_engine/engine/scene"
    "github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/infiniteyak/retro_engine/engine/shader"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

const (
    Title string = "Shader Sandbox"
)

type Game struct {
    screenView *utility.View //view equiv of the full screen
    curScene *scene.Scene
    ecs *ecs.ECS
    states map[scene.SceneId]map[scene.SceneEventId]func() 
}

func NewGame(width, height float64) *Game {
	world := donburi.NewWorld()
	ecs := ecs.NewECS(world)
    this := &Game{
        screenView: utility.NewView(0.0, 0.0, width, height),
        ecs: ecs,
    }

    this.InitStates()
    shader_sb_assets.InitAssets()
    shader.InitShaders(width, height)

    this.curScene = scene.NewScene(this.ecs)
    this.Transition(Init_sceneEvent)
    
    //this.ecs.AddSystem(system.Velocity.Update)
    //this.ecs.AddSystem(system.ViewBound.Update)
    //this.ecs.AddSystem(system.PosTween.Update)
    //this.ecs.AddSystem(system.Wrap.Update)
    //this.ecs.AddSystem(system.Collisions.Update)

    //this.ecs.AddSystem(system.AnimateGraphicObjects.Update)
    this.ecs.AddSystem(system.Input.Update)
    //this.ecs.AddSystem(system.TextInput.Update)
    //this.ecs.AddSystem(system.Damage.Update)
    //this.ecs.AddSystem(system.Health.Update)

    this.ecs.AddSystem(system.Action.Update)

    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsBG.Draw)
    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsFG.Draw)
    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsHudBG.Draw)
    this.ecs.AddRenderer(layer.Foreground, system.DrawGraphicObjectsHudFG.Draw)
    //this.ecs.AddRenderer(layer.Foreground, system.DrawColliders.Draw)
    return this
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return int(this.screenView.Area.Max.X), int(this.screenView.Area.Max.Y)
}

func (this *Game) Update() error {
    if ebiten.IsWindowBeingClosed() {
        this.Exit()
        return nil
    }
    events.ProcessAllEvents(this.ecs.World)
	this.ecs.Update()
	return nil
}

func (this *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	this.ecs.DrawLayer(layer.Background, screen)
	this.ecs.DrawLayer(layer.Foreground, screen)
	this.ecs.DrawLayer(layer.HudBackground, screen)
	this.ecs.DrawLayer(layer.HudForeground, screen)
    shader.RunShaders(screen)
}

func (this *Game) Exit() {
    os.Exit(0)
}
