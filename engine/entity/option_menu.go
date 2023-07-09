package entity

import (
    "fmt"
    "strconv"
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/layer"
    "github.com/infiniteyak/retro_engine/engine/component"
    //"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/yohamta/donburi"
    "github.com/yohamta/donburi/ecs"
    "github.com/infiniteyak/retro_engine/engine/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

type OptionAlignX int
const (
    Left_optionalignx OptionAlignX = iota
    Center_optionalignx
    Right_optionalignx
    Inner_optionalignx
    Outer_optionalignx
)

//TODO combine with game menu?
type OptionMenuFormat struct {
    XAlign OptionAlignX
    YAlign FontAlignY
    ItemFont string
    SelectFont string
    Kerning int
    Spacing int
    SelectPad int
    SelectSprite string
    Gap int
    Title string
    TitleFont string
    TitleKerning int
}

const (
    MaxIncrement = 10
    MinIncrement = 0
)

type OptionType int
const (
    Undefined_optiontype OptionType = iota
    Slider_optiontype
    Button_optiontype
    Number_optiontype
)

type SliderOptionData struct {
    ecs *ecs.ECS
    view *utility.View

    x float64
    y float64

    lineEntity *donburi.Entity
    lineHeight int
    lineWidth int

    cursorEntity *donburi.Entity
    cursorHeight int
    cursorWidth int

    increments int //0 to 10, 0 is off
    setIncFunc func(int)
}

type NumberOptionData struct {
    ecs *ecs.ECS
    view *utility.View

    x float64
    y float64

    value float64
    increment float64 //how much to change by
    valueAdjust float64 //multiply value for display purposes
    setFunc func(float64)
    numberStringData *StringData
    numberFormatString string
    displayDigits int
    minValue float64 
    maxValue float64 
    font string
}

type ButtonOptionData struct {
    ecs *ecs.ECS
    view *utility.View

    x float64
    y float64

    buttonFunction func()
}

type Option interface {
    SetPosition(x, y float64)
    Increment() bool
    Decrement() bool
    Toggle() bool
    GetType() OptionType
}

type optionMenuData struct {
    rSelect map[string]*donburi.Entity
    lSelect map[string]*donburi.Entity
    selectIndex int
}

func (this *optionMenuData) init() {
    this.rSelect = map[string]*donburi.Entity{}
    this.lSelect = map[string]*donburi.Entity{}
}

func AddOptionMenu( ecs *ecs.ECS, 
                    x, y float64, 
                    options map[string]Option,
                    display []string,
                    format OptionMenuFormat,
                    menuNoise string,
                    view *utility.View ) {
    this := &optionMenuData{}
    this.init()

    var curY float64

    displayHeight := float64(len(display) * asset.FontHeight + (len(display) - 1) * format.Spacing)
    switch format.YAlign {
    case Top_fontaligny:
        curY = y
    case Middle_fontaligny:
        curY = y - displayHeight/2.0
    case Bottom_fontaligny:
        curY = y - displayHeight
    }
    if format.Title != "" {
        displayHeight += asset.FontHeight + float64(format.Spacing) 
        text := StringData{
            String: format.Title,
            XAlign: Center_fontalignx,
            YAlign: format.YAlign,
            Kerning: format.TitleKerning, //TODO separate for title?
            Font: format.TitleFont,
        }

        AddSpriteText(
            ecs, 
            x, 
            curY,
            view,
            layer.HudForeground,
            &text,
        )
        curY += asset.FontHeight + float64(format.Spacing)
    }


    var axt FontAlignX
    switch format.XAlign {
    case Left_optionalignx:
        axt = Left_fontalignx
    case Right_optionalignx:
        axt = Right_fontalignx
    case Center_optionalignx:
        axt = Center_fontalignx
    case Inner_optionalignx:
        axt = Right_fontalignx
    case Outer_optionalignx:
        axt = Left_fontalignx
    }

    for _, opt := range display {
        text := StringData{
            String: opt,
            XAlign: axt,
            YAlign: format.YAlign,
            Kerning: format.Kerning,
            Font: format.ItemFont,
        }

        textX := x - float64(format.Gap)

        // In this case we want to center the button regardless
        if options[opt].GetType() == Button_optiontype {
            text.XAlign = Center_fontalignx   
            textX = x
        } 

        AddSpriteText(
            ecs, 
            textX, 
            curY,
            view,
            layer.HudForeground,
            &text,
        )

        //TODO should the positioning account for the selector size?
        displayLength := float64(len(opt) * asset.FontWidth + (len(opt) - 1) * format.Kerning)
        var leftX, rightX float64
        selWidth := float64(asset.SpriteAssets[format.SelectSprite].File.FrameWidth)
        switch text.XAlign {
        case Left_fontalignx:
            leftX = textX - selWidth/2.0 - float64(format.Kerning) - float64(format.SelectPad)
            rightX = textX + displayLength + selWidth/2.0 + float64(format.Kerning) + float64(format.SelectPad -1)
        case Center_fontalignx:
            leftX = textX - displayLength/2.0 - selWidth/2.0 - float64(format.Kerning) - float64(format.SelectPad)
            rightX = textX + displayLength/2.0 + selWidth/2.0 + float64(format.Kerning) + float64(format.SelectPad -1)
        case Right_fontalignx:
            leftX = textX - displayLength - selWidth - float64(format.Kerning) - float64(format.SelectPad)
            rightX = textX + float64(format.Kerning) + float64(format.SelectPad -1)
        }

        //TODO FIX
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

        optionX := x + float64(format.Gap)
        options[opt].SetPosition(optionX, curY - 1)

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
            println("up")
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
            println("down")
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
            if options[display[this.selectIndex]].Toggle() {
                asset.PlaySound(menuNoise)
            }
        },
    )
    AddHybridInputTrigger(
        ecs,
        ebiten.KeyRight,
        inputDelay,
        inputFreq,
        func() {
            if options[display[this.selectIndex]].Increment() {
                asset.PlaySound(menuNoise)
            }
        },
    )
    AddHybridInputTrigger(
        ecs,
        ebiten.KeyLeft,
        inputDelay,
        inputFreq,
        func() {
            if options[display[this.selectIndex]].Decrement() {
                asset.PlaySound(menuNoise)
            }
        },
    )
}

//TODO move to assets or something?
func generateRectVertices(width, height int) []ebiten.Vertex {
    h := float32(height)
    w := float32(width)

    shape := []ebiten.Vertex{}
    var cr float32 = 1.0
    var cg float32 = 1.0
    var cb float32 = 1.0
    shape = append(shape, ebiten.Vertex{
        DstX: -1.0 * w/2.0,
        DstY: h/2.0,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: w/2.0,
        DstY: h/2.0,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: w/2.0,
        DstY: -1.0 * h/2.0,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{
        DstX: -1.0 * w/2.0,
        DstY: -1.0 * h/2.0,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    shape = append(shape, ebiten.Vertex{ //center
        DstX: 0.0,
        DstY: 0.0,
        SrcX: 0.0,
        SrcY: 0.0,
        ColorR: cr,
        ColorG: cg,
        ColorB: cb,
        ColorA: 1,
    })
    return shape
}

// Slider
func AddSliderOption ( ecs *ecs.ECS, 
                       format OptionMenuFormat, 
                       initialIncr int, 
                       setIncFunc func(int), 
                       view *utility.View ) *SliderOptionData {
    this := &SliderOptionData{}
    this.ecs = ecs
    this.view = view

    this.increments = initialIncr 
    this.setIncFunc = setIncFunc

    this.lineHeight = 1
    this.lineWidth = 80
    this.lineEntity = AddPolygonObject( this.ecs,
                                        layer.HudForeground,
                                        0, 0,
                                        generateRectVertices(this.lineWidth, this.lineHeight),
                                        view)

    this.cursorHeight = 7
    this.cursorWidth = 2
    this.cursorEntity = AddPolygonObject( this.ecs,
                                        layer.HudForeground,
                                        0, 0,
                                        generateRectVertices(this.cursorWidth, this.cursorHeight),
                                        view)

    return this
}

func (this *SliderOptionData) UpdateCursor() {
    incrementOffset := float64((this.lineWidth / MaxIncrement)  * this.increments)
    ce := this.ecs.World.Entry(*this.cursorEntity)
    cePos := component.Position.Get(ce)
    cePos.Point.X = this.x + incrementOffset // + float64(this.lineWidth)/2.0
    cePos.Point.Y = this.y + float64(this.lineHeight)/2.0
}

func (this *SliderOptionData) SetPosition(x, y float64) {
    println("SetPosition")

    this.x = x
    this.y = y

    le := this.ecs.World.Entry(*this.lineEntity)
    lePos := component.Position.Get(le)
    lePos.Point.X = this.x + float64(this.lineWidth)/2.0
    lePos.Point.Y = this.y + float64(this.lineHeight)/2.0

    this.UpdateCursor()
}

func (this *SliderOptionData) Increment() bool {
    println("Increment")
    if this.increments < MaxIncrement {
        this.increments++
        this.setIncFunc(this.increments)
        this.UpdateCursor()
    } else {
        return false
    }
    return true
}

func (this *SliderOptionData) Decrement() bool {
    println("Decrement")
    if this.increments > MinIncrement {
        this.increments--
        this.setIncFunc(this.increments)
        this.UpdateCursor()
    } else {
        return false
    }
    return true
}

func (this *SliderOptionData) Toggle() bool {
    println("Toggle")
    return false
}

func (this *SliderOptionData) GetType() OptionType {
    return Slider_optiontype
}

// Button
func AddButtonOption ( ecs *ecs.ECS, 
                       format OptionMenuFormat, 
                       foo func(), 
                       view *utility.View ) *ButtonOptionData {
    this := &ButtonOptionData{}
    this.ecs = ecs
    this.view = view

    this.buttonFunction = foo

    return this
}

func (this *ButtonOptionData) SetPosition(x, y float64) {
    println("SetPosition")
}

func (this *ButtonOptionData) Increment() bool {
    return false
}

func (this *ButtonOptionData) Decrement() bool {
    return false
}

func (this *ButtonOptionData) Toggle() bool {
    println("Toggle")
    this.buttonFunction()
    return true
}

func (this *ButtonOptionData) GetType() OptionType {
    return Button_optiontype
}

// Number
//TODO standardize
func AddNumberOption ( ecs *ecs.ECS, 
                       format OptionMenuFormat, //TODO why do we want this? 
                       initialValue float64, 
                       minValue float64, 
                       maxValue float64, 
                       increment float64, 
                       valueAdjust float64, 
                       displayDigits int, 
                       setFunc func(float64), 
                       view *utility.View ) *NumberOptionData {
    this := &NumberOptionData{}
    this.ecs = ecs
    this.view = view

    this.value = initialValue 
    this.increment = increment 
    this.valueAdjust = valueAdjust 
    this.setFunc = setFunc
    //this.x = x
    //this.y = y
    this.minValue = minValue
    this.maxValue = maxValue
    this.displayDigits = displayDigits

    this.numberFormatString = "%0" + strconv.Itoa(displayDigits) + "d"

    this.numberStringData = AddNormalText(
        this.ecs, 
        this.x,
        this.y,
        view,
        "WhiteFont",
        fmt.Sprintf(this.numberFormatString, int(this.value * this.valueAdjust)),
    )
    this.numberStringData.XAlign = Left_fontalignx

    return this
}

// TODO okay now that I think about it I should really redo this whole thing so that
// you create a menu and then that menu has a function to add options, which then
// inherit formatting instead of creating the options first.
func NewNumberOptionData(format OptionMenuFormat) *NumberOptionData {
    return &NumberOptionData {
        font: format.SelectFont,
    }
}

func (this *NumberOptionData) Init() *NumberOptionData {
    this.numberFormatString = "%0" + strconv.Itoa(this.displayDigits) + "d"
    this.numberStringData = AddNormalText(
        this.ecs, 
        this.x,
        this.y,
        this.view,
        "WhiteFont",
        fmt.Sprintf(this.numberFormatString, int(this.value * this.valueAdjust)),
    )
    this.numberStringData.XAlign = Left_fontalignx
    return this
}

func (this *NumberOptionData) Increment() bool {
    println("Increment")
    if this.value < this.maxValue {
        this.value += this.increment
        this.numberStringData.String = fmt.Sprintf(this.numberFormatString, int(this.value * this.valueAdjust))
        this.setFunc(this.value)
    } else {
        return false
    }
    return true
}

func (this *NumberOptionData) Decrement() bool {
    println("Increment")
    if this.value > this.minValue {
        this.value -= this.increment
        this.numberStringData.String = fmt.Sprintf(this.numberFormatString, int(this.value * this.valueAdjust))
        this.setFunc(this.value)
    } else {
        return false
    }
    return true
}

func (this *NumberOptionData) Toggle() bool {
    println("Toggle")
    return false
}

func (this *NumberOptionData) GetType() OptionType {
    return Number_optiontype
}

func (this *NumberOptionData) SetPosition(x, y float64) {
    //this.x = x + 24
    this.x = x
    this.y = y
    this.numberStringData.X = this.x
    this.numberStringData.Y = this.y
}
