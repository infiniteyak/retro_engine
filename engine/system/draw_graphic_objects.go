package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
)

type drawGraphicObjects struct {
	query *query.Query
}

var DrawGraphicObjectsFG = &drawGraphicObjects{
	query: ecs.NewQuery(
		layer.Foreground,
		filter.Contains(
			component.Position,
			component.GraphicObject,
			component.View,
		)),
}

var DrawGraphicObjectsBG = &drawGraphicObjects{
	query: ecs.NewQuery(
		layer.Background,
		filter.Contains(
			component.Position,
			component.GraphicObject,
			component.View,
		)),
}

var DrawGraphicObjectsHudBG = &drawGraphicObjects{
	query: ecs.NewQuery(
		layer.HudBackground,
		filter.Contains(
			component.Position,
			component.GraphicObject,
			component.View,
		)),
}

var DrawGraphicObjectsHudFG = &drawGraphicObjects{
	query: ecs.NewQuery(
		layer.HudForeground,
		filter.Contains(
			component.Position,
			component.GraphicObject,
			component.View,
		)),
}

//TODO need to add some kind of masking to prevent drawing outside of
// the view
func (this *drawGraphicObjects) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		position := component.Position.Get(entry)
		view := component.View.Get(entry).View
		gobj := component.GraphicObject.Get(entry)

        if *gobj.TransInfo.Hide {
            return
        }

        //TODO depth?
        for _, renderable := range gobj.Renderables {
            tinfo := &component.TransformInfo{
                Rotation: gobj.TransInfo.Rotation,
                Scale: gobj.TransInfo.Scale,
                Offset: &utility.Point{
                    X: position.Point.X + view.Offset.X + gobj.TransInfo.Offset.X, 
                    Y: position.Point.Y + view.Offset.Y + gobj.TransInfo.Offset.Y,
                },
            }
            renderable.Draw(screen, tinfo)
        }
	})
}
