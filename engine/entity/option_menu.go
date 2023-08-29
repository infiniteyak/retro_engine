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

//TODO combine with game menu? or with optionMenuData?
type OptionMenuFormat struct {
    XAlign OptionAlignX
    YAlign FontAlignY
    ItemFont string //TODO name?
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
    rSelect map[string]*donburi.Entity //each menu item has hidden sprites that are shown when selected
    lSelect map[string]*donburi.Entity
    selectIndex int
    displayOrder []string
    displayText map[string]*StringData
    options map[string]Option
    format OptionMenuFormat
    menuNoise string
    view *utility.View
    x float64
    y float64
    title StringData 
    ecs *ecs.ECS
}

func (this *optionMenuData) init() {
    this.rSelect = map[string]*donburi.Entity{}
    this.lSelect = map[string]*donburi.Entity{}
    this.displayOrder = []string{}
    this.options = map[string]Option{} //TODO should this be a map of string to pointer?
    this.displayText = map[string]*StringData{}
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
    this.x = x
    this.y = y

    le := this.ecs.World.Entry(*this.lineEntity)
    lePos := component.Position.Get(le)
    lePos.Point.X = this.x + float64(this.lineWidth)/2.0
    lePos.Point.Y = this.y + float64(this.lineHeight)/2.0

    this.UpdateCursor()
}

func (this *SliderOptionData) Increment() bool {
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
}

func (this *ButtonOptionData) Increment() bool {
    return false
}

func (this *ButtonOptionData) Decrement() bool {
    return false
}

func (this *ButtonOptionData) Toggle() bool {
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
    if this.value + this.increment < this.maxValue {
        this.value += this.increment
    } else {
        this.value = this.minValue
    }
    this.numberStringData.String = fmt.Sprintf(this.numberFormatString, int(this.value * this.valueAdjust))
    this.setFunc(this.value)
    return true
}

func (this *NumberOptionData) Decrement() bool {
    if this.value - this.increment > this.minValue {
        this.value -= this.increment
    } else {
        this.value = this.maxValue
    }
    this.numberStringData.String = fmt.Sprintf(this.numberFormatString, int(this.value * this.valueAdjust))
    this.setFunc(this.value)
    return true
}

func (this *NumberOptionData) Toggle() bool {
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

func AddOptionMenu( ecs *ecs.ECS, 
                    x, y float64, 
                    format OptionMenuFormat,
                    menuNoise string,
                    view *utility.View ) *optionMenuData {
    this := &optionMenuData{}
    this.init()
    this.format = format
    this.menuNoise = menuNoise
    this.view = view
    this.x = x
    this.y = y
    this.ecs = ecs

    //create a title object
    if this.format.Title != "" {
        var curY float64
        switch this.format.YAlign {
        case Top_fontaligny:
            curY = this.y
        case Middle_fontaligny:
            curY = this.y - asset.FontHeight/2.0
        case Bottom_fontaligny:
            curY = this.y -asset.FontHeight 
        }
        this.title = StringData{
            String: this.format.Title,
            XAlign: Center_fontalignx,
            YAlign: this.format.YAlign,
            Kerning: this.format.TitleKerning,
            Font: this.format.TitleFont,
        }
        AddSpriteText(
            this.ecs, 
            this.x, 
            curY,
            this.view,
            layer.HudForeground,
            &this.title,
        )
    }

    moveSelect := func() {
        for index, opt := range this.displayOrder {
            rs := component.GraphicObject.Get(ecs.World.Entry(*this.rSelect[opt]))
            ls := component.GraphicObject.Get(ecs.World.Entry(*this.lSelect[opt]))
            *rs.TransInfo.Hide = !(index == this.selectIndex)
            *ls.TransInfo.Hide = !(index == this.selectIndex)
        }
    }
    //moveSelect() //TODO should this be here or no?

    AddHybridInputTrigger(
        this.ecs,
        ebiten.KeyUp,
        inputDelay,
        inputFreq,
        func() {
            oldIndex := this.selectIndex
            this.selectIndex--
            if this.selectIndex < 0 {
                this.selectIndex = len(this.displayOrder)-1
            }
            if oldIndex != this.selectIndex {
                asset.PlaySound(menuNoise)
                moveSelect()
            }
        },
    )
    AddHybridInputTrigger(
        this.ecs,
        ebiten.KeyDown,
        inputDelay,
        inputFreq,
        func() {
            oldIndex := this.selectIndex
            this.selectIndex = (this.selectIndex + 1) % len(this.displayOrder)
            if oldIndex != this.selectIndex {
                asset.PlaySound(menuNoise)
                moveSelect()
            }
        },
    )
    AddHybridInputTrigger(
        this.ecs,
        ebiten.KeySpace,
        inputDelay,
        inputFreq,
        func() {
            if this.options[this.displayOrder[this.selectIndex]].Toggle() {
                asset.PlaySound(menuNoise)
            }
        },
    )
    AddHybridInputTrigger(
        this.ecs,
        ebiten.KeyRight,
        inputDelay,
        inputFreq,
        func() {
            if this.options[this.displayOrder[this.selectIndex]].Increment() {
                asset.PlaySound(menuNoise)
            }
        },
    )
    AddHybridInputTrigger(
        this.ecs,
        ebiten.KeyLeft,
        inputDelay,
        inputFreq,
        func() {
            if this.options[this.displayOrder[this.selectIndex]].Decrement() {
                asset.PlaySound(menuNoise)
            }
        },
    )
    return this
}

//TODO add some kind of pagination

func (this *optionMenuData) AddOption( displayName string,
                                       option Option) {
    //append option
    this.displayOrder = append(this.displayOrder, displayName)
    this.options[displayName] = option

    //create option text / label

    //determine x alignment and value

    var axt FontAlignX
    var textX float64

    //buttons should be centered
    if option.GetType() == Button_optiontype {
        axt = Center_fontalignx
        textX = this.x
    } else {
        switch this.format.XAlign {
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
        textX = this.x - float64(this.format.Gap)
    }

    // create display name text
    text := StringData{
        String: displayName,
        XAlign: axt,
        YAlign: this.format.YAlign,
        Kerning: this.format.Kerning,
        Font: this.format.ItemFont,
    }
    AddSpriteText(
        this.ecs, 
        textX, 
        0, //We'll set this later
        this.view,
        layer.HudForeground,
        &text,
    )
    this.displayText[displayName] = &text

    //add selector sprites for new option
    //TODO should the positioning account for the selector size?
    displayLength := float64(len(displayName) * asset.FontWidth + (len(displayName) - 1) * this.format.Kerning)
    var leftX, rightX float64
    selWidth := float64(asset.SpriteAssets[this.format.SelectSprite].File.FrameWidth)
    switch text.XAlign {
    case Left_fontalignx:
        leftX = textX - selWidth/2.0 - float64(this.format.Kerning) - float64(this.format.SelectPad)
        rightX = textX + displayLength + selWidth/2.0 + float64(this.format.Kerning) + float64(this.format.SelectPad -1)
    case Center_fontalignx:
        leftX = textX - displayLength/2.0 - selWidth/2.0 - float64(this.format.Kerning) - float64(this.format.SelectPad)
        rightX = textX + displayLength/2.0 + selWidth/2.0 + float64(this.format.Kerning) + float64(this.format.SelectPad -1)
    case Right_fontalignx:
        leftX = textX - displayLength - selWidth - float64(this.format.Kerning) - float64(this.format.SelectPad)
        rightX = textX + float64(this.format.Kerning) + float64(this.format.SelectPad -1)
    }

    //TODO FIX
    this.lSelect[displayName] = AddSpriteObject(
        this.ecs,
        layer.HudForeground,
        leftX,
        0, //this will get updated later
        this.format.SelectSprite,
        "left",
        this.view,
    )
    this.rSelect[displayName] = AddSpriteObject(
        this.ecs,
        layer.HudForeground,
        rightX,
        0, //this will get updated later
        this.format.SelectSprite,
        "right",
        this.view,
    )

    //calculate display height for options
    textHeight := len(this.displayOrder) * asset.FontHeight
    spacingHeight := (len(this.displayOrder) - 1) * this.format.Spacing
    displayHeight := float64(textHeight + spacingHeight)
    if this.format.Title != "" {
        displayHeight += asset.FontHeight + float64(this.format.Spacing) 
    }

    //set curY based on alignment options
    var curY float64
    switch this.format.YAlign {
    case Top_fontaligny:
        curY = this.y
    case Middle_fontaligny:
        curY = this.y - displayHeight/2.0
    case Bottom_fontaligny:
        curY = this.y - displayHeight
    }
    
    //adjust title location
    if this.format.Title != "" {
        this.title.Y = curY
        curY += asset.FontHeight + float64(this.format.Spacing)
    }

    //now go through and set all the position values as needed
    for index, opt := range this.displayOrder {
        dt := this.displayText[opt]
        dt.Y = curY
        this.displayText[opt] = dt

        this.options[opt].SetPosition(this.x + float64(this.format.Gap), curY -1) //TODO why -1?

        rsPos := component.Position.Get(this.ecs.World.Entry(*this.rSelect[opt]))
        rsPos.Point.Y = curY
        rsGO := component.GraphicObject.Get(this.ecs.World.Entry(*this.rSelect[opt]))
        *rsGO.TransInfo.Hide = !(index == this.selectIndex)

        lsPos := component.Position.Get(this.ecs.World.Entry(*this.lSelect[opt]))
        lsPos.Point.Y = curY
        lsGO := component.GraphicObject.Get(this.ecs.World.Entry(*this.lSelect[opt]))
        *lsGO.TransInfo.Hide = !(index == this.selectIndex)

        curY += asset.FontHeight + float64(this.format.Spacing)
    }
}
