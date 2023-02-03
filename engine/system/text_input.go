package system

import (
	"strings"
	"unicode"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
    "github.com/infiniteyak/retro_engine/engine/component"
)

type textInput struct {
	query *query.Query
}

var TextInput = &textInput{
	query: query.NewQuery(
		filter.Contains(
			component.TextInput,
		)),
}

func (this *textInput) Update(ecs *ecs.ECS) {
	this.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		ti := component.TextInput.Get(entry)

        var runes []rune
        runes = ebiten.AppendInputChars(runes[:0])
        *ti.String += string(runes)

        var str string
        for _, c := range *ti.String {
            if unicode.IsLetter(c) {
                str += string(c)
            }
        }
        str = strings.ToUpper(str)
        *ti.String = str

        if len(*ti.String) > ti.Length {
            *ti.String = (*ti.String)[0:3]
        }
        if repeatingKeyPressed(ebiten.KeyBackspace) && len(*ti.String) >= 1 {
            *ti.String = (*ti.String)[:len(*ti.String)-1]
        }
        if len(*ti.String) == ti.Length {
            ti.Function()
        }
        if len(*ti.String) < ti.Length {
            *ti.String += strings.Repeat("_", ti.Length - len(*ti.String))
        }
	})
}

func repeatingKeyPressed(key ebiten.Key) bool {
    const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
