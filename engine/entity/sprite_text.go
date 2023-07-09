package entity

import (
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    "github.com/infiniteyak/retro_engine/engine/event"
    "github.com/infiniteyak/retro_engine/engine/asset"

	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/tanema/gween"
)

type FontAlignY int
const (
    Top_fontaligny FontAlignY = iota
    Middle_fontaligny
    Bottom_fontaligny
)
type FontAlignX int
const (
    Left_fontalignx FontAlignX = iota
    Center_fontalignx
    Right_fontalignx
)

type StringData struct {
    String string
    XAlign FontAlignX
    YAlign FontAlignY
    Font string
    Kerning int
    Blink bool
    TypeWriter int
    Delay int
    TweenDelay int
    XTween *gween.Tween
    YTween *gween.Tween
    Entity *donburi.Entity
    X float64
    Y float64
}

func writeText(curString *StringData, gobj *component.GraphicObjectData, hide bool) {
    var textY float64
    switch curString.YAlign {
    case Top_fontaligny:
        textY = float64(asset.FontHeight/2)
    case Middle_fontaligny:
        textY = 0.0
    case Bottom_fontaligny:
        textY = -1.0 * float64(asset.FontHeight/2)
    }

    var textX float64
    switch curString.XAlign {
    case Left_fontalignx:
        textX = float64(asset.FontWidth/2) //Because each sprite pos is char center
    case Center_fontalignx:
        textX = float64(
(asset.FontWidth * (1 - len(curString.String)) + -1 * curString.Kerning * (len(curString.String) - 1)) / 2)
    case Right_fontalignx:
        textX = float64(asset.FontWidth * (0 - (len(curString.String) + curString.Kerning * (len(curString.String) - 1))))
    }

    gobj.Renderables = []component.Renderable{}
    for _, c := range curString.String {
        nsd := component.SpriteData{}
        m := asset.FontMasks[string(c)]
        nsd.Load(curString.Font, &m)
        tinfo := nsd.RenderableData.GetTransInfo()
        tinfo.Offset.X = textX
        tinfo.Offset.Y = textY
        *tinfo.Hide = hide
        nsd.RenderableData.SetTransInfo(tinfo)
        gobj.Renderables = append(gobj.Renderables, &nsd)
        textX += float64(asset.FontWidth) + float64(curString.Kerning)
    }
}

func AddSpriteText(ecs *ecs.ECS, x, y float64, view *utility.View, layer ecs.LayerID, str *StringData) *donburi.Entity {
    entity := ecs.Create(
        layer,
        component.Position, 
        component.GraphicObject,
        component.View,
        component.Actions,
        component.PosTween,
    )
    event.RegisterEntityEvent.Publish(ecs.World, event.RegisterEntity{Entity:&entity})

    entry := ecs.World.Entry(entity)

    str.X = x   
    str.Y = y   

    curString := *str

    // Position
    pd := component.NewPositionData(x, y)
    donburi.SetValue(entry, component.Position, pd)

    // Graphic Object
    gobj := component.NewGraphicObjectData()
    *gobj.TransInfo.Hide = true // everything start out hidden, then is shown by upkeep
    writeText(&curString, &gobj, (curString.TypeWriter != 0))
    donburi.SetValue(entry, component.GraphicObject, gobj)

    // View
    donburi.SetValue(entry, component.View, component.ViewData{View:view})

    // PosTween
    comp := component.PosTweenData{
        XTween: str.XTween,
        YTween: str.YTween,
        Delay: str.TweenDelay,
    }
    donburi.SetValue(entry, component.PosTween, comp)

    // Actions
    blinkCounter := 0
    typeWriterCounter := curString.TypeWriter
    delayCounter := curString.Delay
    ad := component.NewActions()
    ad.AddUpkeepAction(func(){
        g := component.GraphicObject.Get(entry)
        if delayCounter <= 0 {
            *g.TransInfo.Hide = false
            if curString != *str {
                curString = *str
                if curString.Delay != 0 {
                    *g.TransInfo.Hide = true
                    delayCounter = curString.Delay
                }
                writeText(&curString, g, (curString.TypeWriter != 0))
                pd.Point.X = curString.X
                pd.Point.Y = curString.Y
                return
            }
            if curString.Blink {
                blinkCounter++
                if (blinkCounter / 20) % 2 == 0 {
                    *g.TransInfo.Hide = true
                } else {
                    *g.TransInfo.Hide = false
                }
            } else {
                *g.TransInfo.Hide = false
            } 
            if curString.TypeWriter != 0 {
                if typeWriterCounter <= 0 {
                    firstHidden := -1
                    for i := 0; i < len(g.Renderables); i++ {
                        if *g.Renderables[i].GetTransInfo().Hide == true {
                            firstHidden = i
                            break
                        }
                    }
                    if firstHidden == -1 {
                        //curString.TypeWriter = 0 //TODO
                    } else {
                        *g.Renderables[firstHidden].GetTransInfo().Hide = false
                        typeWriterCounter = curString.TypeWriter
                    }
                } else {
                    typeWriterCounter--
                }
            }
        } else {
            delayCounter--
            /*
            if delayCounter == 0 {
                *g.TransInfo.Hide = false
            }
            */
        }
    })

    donburi.SetValue(entry, component.Actions, ad)

    return &entity
}

//convenient preset for doing title text
func AddTitleText(ecs *ecs.ECS, x, y float64, view *utility.View, str string) *StringData {
    text := StringData{
        String: str,
        XAlign: Center_fontalignx,
        YAlign: Middle_fontaligny,
        Kerning: 2,
        Font: "TitleFont",
        TypeWriter: 0,
        Delay: 0,
    }

    entity := AddSpriteText(
        ecs, 
        x, 
        y,
        view,
        layer.HudForeground,
        &text,
    )

    text.Entity = entity

    return &text
}

//convenient preset for doing white text
func AddNormalText(ecs *ecs.ECS, 
                   x, y float64, 
                   view *utility.View, 
                   font, str string) *StringData {
    text := StringData{
        String: str,
        XAlign: Center_fontalignx,
        YAlign: Middle_fontaligny,
        Kerning: 0,
        Font: font,
        TypeWriter: 0,
        Delay: 0,
    }

    AddSpriteText(
        ecs, 
        x, 
        y,
        view,
        layer.HudForeground,
        &text,
    )

    return &text
}

//convenient preset for doing white text
func AddNormalTweenText(ecs *ecs.ECS, 
                        x, y float64, 
                        xTween, yTween *gween.Tween, 
                        tweenDelay int, 
                        view *utility.View, 
                        font, str string) *StringData {
    text := StringData{
        String: str,
        XAlign: Center_fontalignx,
        YAlign: Middle_fontaligny,
        Kerning: 0,
        Font: font,
        TypeWriter: 0,
        Delay: 0,
        TweenDelay: tweenDelay,
        XTween: xTween,
        YTween: yTween,
    }

    AddSpriteText(
        ecs, 
        x, 
        y,
        view,
        layer.HudForeground,
        &text,
    )

    return &text
}
