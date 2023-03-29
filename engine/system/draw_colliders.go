package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/layer"
	"image/color"
    "github.com/infiniteyak/retro_engine/engine/shader"
)

type drawColliders struct {
	query *query.Query
}

var DrawColliders = &drawColliders{
	query: ecs.NewQuery(
		layer.Foreground,
		filter.Contains(
			component.Position,
			component.View,
            component.Collider,
		)),
}

func (this *drawColliders) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		position := component.Position.Get(entry)
		view := component.View.Get(entry).View
		hitboxes := component.Collider.Get(entry).Hitboxes

        for _, hb := range hitboxes {
            op := &ebiten.DrawImageOptions{}

            img := ebiten.NewImage(hb.Radius*2, hb.Radius*2)
            color := color.RGBA{120, 226, 160, 255} //TODO configurable
            drawCircle(img, hb.Radius, hb.Radius, hb.Radius, color)

            hbw, hbh := img.Size()

            newX := position.Point.X-float64(hbw)/2.0
            newX += hb.Offset.X
            newX += view.Offset.X

            newY := position.Point.Y-float64(hbh)/2.0
            newY += hb.Offset.Y
            newY += view.Offset.Y

            op.GeoM.Translate(newX, newY)
            //screen.DrawImage(img, op)
            shader.Image0.DrawImage(img, op) // TODO why is this so slow?
        }
	})
}

func drawCircle(img *ebiten.Image, x0, y0, r int, c color.Color) {
    x, y, dx, dy := r-1, 0, 1, 1
    err := dx - (r * 2)

    for x > y {
        img.Set(x0+x, y0+y, c)
        img.Set(x0+y, y0+x, c)
        img.Set(x0-y, y0+x, c)
        img.Set(x0-x, y0+y, c)
        img.Set(x0-x, y0-y, c)
        img.Set(x0-y, y0-x, c)
        img.Set(x0+y, y0-x, c)
        img.Set(x0+x, y0-y, c)

        if err <= 0 {
            y++
            err += dy
            dy += 2
        }
        if err > 0 {
            x--
            dx += 2
            err += dx - (r * 2)
        }
    }
}
