package entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    //"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

//TODO should I just combine this with option_menu?

const (
    inputDelay = 50
    inputFreq = 30
)

type GameMenuFormat struct {
    XAlign FontAlignX
    YAlign FontAlignY
    Font string
    Kerning int
    Spacing int
    SelectPad int
    SelectSprite string
}

type gameMenuData struct {
    rSelect map[string]*donburi.Entity
    lSelect map[string]*donburi.Entity
    selectIndex int
}

func (this *gameMenuData) init() {
    this.rSelect = map[string]*donburi.Entity{}
    this.lSelect = map[string]*donburi.Entity{}
}

func AddGameMenu(   ecs *ecs.ECS, 
                    x, y float64, 
                    options map[string]func(),
                    display []string,
                    format GameMenuFormat,
                    menuNoise string,
                    view *utility.View ) {
    this := &gameMenuData{}
    this.init()

    var curY float64

    displayHeight := float64(len(display) * asset.FontHeight + (len(display) - 1) * format.Spacing)

    switch format.YAlign {
    case Top_fontaligny:
        curY = y
    case Middle_fontaligny:
        curY = y - (displayHeight - asset.FontHeight)/2.0
    case Bottom_fontaligny:
        curY = y - (displayHeight - asset.FontHeight)
    }

    for _, opt := range display {
        text := StringData{
            String: opt,
            XAlign: format.XAlign,
            YAlign: format.YAlign,
            Kerning: format.Kerning,
            Font: format.Font,
        }

        AddSpriteText(
            ecs, 
            x, 
            curY,
            view,
            layer.HudForeground,
            &text,
        )

        //TODO should the positioning account for the selector size?
        displayLength := float64(len(opt) * asset.FontWidth + (len(opt) - 1) * format.Kerning)
        var leftX, rightX float64
        selWidth := float64(asset.SpriteAssets[format.SelectSprite].File.FrameWidth)
        switch format.XAlign {
        case Left_fontalignx:
            leftX = x - selWidth/2.0 - float64(format.Kerning) - float64(format.SelectPad)
            rightX = x + displayLength + selWidth/2.0 + float64(format.Kerning) + float64(format.SelectPad -1)
        case Center_fontalignx:
            leftX = x - displayLength/2.0 - selWidth/2.0 - float64(format.Kerning) - float64(format.SelectPad)
            rightX = x + displayLength/2.0 + selWidth/2.0 + float64(format.Kerning) + float64(format.SelectPad -1)
        case Right_fontalignx:
            leftX = x - displayLength - selWidth - float64(format.Kerning) - float64(format.SelectPad)
            rightX = x + float64(format.Kerning) + float64(format.SelectPad -1)
        }
        this.lSelect[opt] = AddSpriteObject(
            ecs,
            layer.HudForeground,
            leftX,
            curY,
            format.SelectSprite,
            "left",
            view,
        )
        this.rSelect[opt] = AddSpriteObject(
            ecs,
            layer.HudForeground,
            rightX,
            curY,
            format.SelectSprite,
            "right",
            view,
        )

        curY += asset.FontHeight + float64(format.Spacing)
    }

    moveSelect := func() {
        for index, opt := range display {
            rs := component.GraphicObject.Get(ecs.World.Entry(*this.rSelect[opt]))
            ls := component.GraphicObject.Get(ecs.World.Entry(*this.lSelect[opt]))
            *rs.TransInfo.Hide = !(index == this.selectIndex)
            *ls.TransInfo.Hide = !(index == this.selectIndex)
        }
    }
    moveSelect()

    AddHybridInputTrigger(
        ecs,
        ebiten.KeyUp,
        inputDelay,
        inputFreq,
        func() {
            oldIndex := this.selectIndex
            this.selectIndex--
            if this.selectIndex < 0 {
                this.selectIndex = len(display)-1
            }
            if oldIndex != this.selectIndex {
                asset.PlaySound(menuNoise)
                moveSelect()
            }
        },
    )
    AddHybridInputTrigger(
        ecs,
        ebiten.KeyDown,
        inputDelay,
        inputFreq,
        func() {
            oldIndex := this.selectIndex
            this.selectIndex = (this.selectIndex + 1) % len(display)
            if oldIndex != this.selectIndex {
                asset.PlaySound(menuNoise)
                moveSelect()
            }
        },
    )
    AddHybridInputTrigger(
        ecs,
        ebiten.KeySpace,
        inputDelay,
        inputFreq,
        func() {
            asset.PlaySound(menuNoise)
            options[display[this.selectIndex]]()
        },
    )
}
