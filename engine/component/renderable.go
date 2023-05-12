package component

import (
    "github.com/yohamta/donburi"
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/infiniteyak/retro_engine/engine/utility"
    "github.com/infiniteyak/retro_engine/engine/asset"
	"image"
	"github.com/solarlune/goaseprite"
	"math"
    "image/color"
    _ "image/png"
    "math/rand"
)

// TODO split this stuff up into multiple files!

type TransformInfo struct {
    Rotation *float64 //radians
    Scale *float64
    Offset *utility.Point
    Hide *bool
}

type Renderable interface {
    GetTransInfo() *TransformInfo
    SetTransInfo(ti *TransformInfo)
    GetDepth() *int
    SetDepth(d *int)
    Init()
    Draw(screen *ebiten.Image, ti *TransformInfo)
    Animate()
    Play(tag string)
}

type GraphicObjectData struct {
    Renderables []Renderable
    TransInfo *TransformInfo
}

func (this *GraphicObjectData) Init() {
    this.Renderables = []Renderable{}
    this.TransInfo = &TransformInfo{
        Rotation: new(float64),
        Scale: new(float64),
        Offset: new(utility.Point),
        Hide: new(bool),
    }
}

func (this *GraphicObjectData) HideAllRenderables(hide bool) {
    for i := 0; i < len(this.Renderables); i++ {
        *this.Renderables[i].GetTransInfo().Hide = hide
    }
}

//TODO what if the renderables don't all have play?
func (this *GraphicObjectData) PlayAllRenderables(tag string) {
    for i := 0; i < len(this.Renderables); i++ {
        this.Renderables[i].Play(tag)
    }
}

func NewGraphicObjectData() GraphicObjectData {
    this := GraphicObjectData{}
    this.Init()
    return this
}

var GraphicObject = donburi.NewComponentType[GraphicObjectData]()

type RenderableData struct {
    tinfo *TransformInfo
    depth *int
}

func (this *RenderableData) GetTransInfo() *TransformInfo {
    return this.tinfo
}

func (this *RenderableData) SetTransInfo(ti *TransformInfo) {
    this.tinfo = ti
}

func (this *RenderableData) GetDepth() *int {
    return this.depth
}

func (this *RenderableData) SetDepth(d *int) {
    this.depth = d
}

func (this *RenderableData) baseInit() {
    this.tinfo = &TransformInfo{
        Rotation: new(float64),
        Scale: new(float64),
        Offset: &utility.Point{X: 0.0, Y:0.0},
        Hide: new(bool),
    }
    this.depth = new(int)
}

type SpriteData struct {
    RenderableData
    file *goaseprite.File
    image *ebiten.Image
    mask *image.Rectangle
}

func NewSpriteData(name string, mask *image.Rectangle, tag string) SpriteData {
    sd := SpriteData{}
    sd.Load(name, mask)
    sd.Play(tag)
    return sd
}

func (this *SpriteData) Init() {
    this.baseInit()
    this.file = nil
    this.image = nil
}

func (this *SpriteData) Load(name string, mask *image.Rectangle) {
    this.baseInit()
    sa := asset.SpriteAssets[name]
    tempFile := *sa.File
    this.file = &tempFile
    this.image = sa.Image
    this.mask = mask
    //this.file.Play("") //TODO probably move this...?
}

func (this *SpriteData) Play(tag string) {
    this.file.Play(tag)
}

func (this *SpriteData) SetLoopCallback(foo func()) {
    this.file.OnLoop = foo
}

func (this *SpriteData) SetPlaySpeed(s float32) {
    this.file.PlaySpeed = s
}

func (this *SpriteData) SetFrame(i int) {
    this.file.SetFrame(i)
}

func (this *SpriteData) Animate() {
    this.file.Update(float32(1.0/120.0))
}

func (this *SpriteData) Draw(screen *ebiten.Image, ti *TransformInfo) {
    if *this.RenderableData.tinfo.Hide {
        return
    }

    op := &ebiten.DrawImageOptions{}

    var sub image.Image
    if this.mask == nil { //This is used for animation
        sub = this.image.SubImage(
            image.Rect(this.file.CurrentFrameCoords()))
    } else { // This is used for text
        sub = this.image.SubImage(*this.mask)
    }

    sw, sh := sub.(*ebiten.Image).Size()

    newX := ti.Offset.X + this.RenderableData.tinfo.Offset.X - float64(sw/2)
    newY := ti.Offset.Y + this.RenderableData.tinfo.Offset.Y - float64(sh/2)

    op.GeoM.Translate(newX, newY)
    screen.DrawImage(sub.(*ebiten.Image), op)
}


type PolygonData struct {
    RenderableData
    vertices []ebiten.Vertex
}

func (this *PolygonData) Init() {
    this.baseInit()
    this.vertices = []ebiten.Vertex{}
}

func (this *PolygonData) Load(verts []ebiten.Vertex) {
    this.baseInit()
    this.vertices = verts
}

func (this *PolygonData) Animate() {
}

func (this *PolygonData) Play(tag string) {
}

func (this *PolygonData) transform(ti *TransformInfo) []ebiten.Vertex {
    nvs := []ebiten.Vertex{}
    for i := 0; i < len(this.vertices); i++ {
        // apply offset
        nx := this.vertices[i].DstX + float32(this.RenderableData.tinfo.Offset.X)
        ny := this.vertices[i].DstY + float32(this.RenderableData.tinfo.Offset.Y)

        // apply rotation
        theta := *ti.Rotation + *this.RenderableData.tinfo.Rotation
        rx := nx * float32(math.Cos(theta)) - ny * float32(math.Sin(theta))
        ry := ny * float32(math.Cos(theta)) + nx * float32(math.Sin(theta))
        nx = rx
        ny = ry
        // apply translation
        nvs = append(nvs, ebiten.Vertex{
            DstX: nx + float32(ti.Offset.X),
            DstY: ny + float32(ti.Offset.Y),
            SrcX: this.vertices[i].SrcX,
            SrcY: this.vertices[i].SrcY,
            ColorR: this.vertices[i].ColorR,
            ColorG: this.vertices[i].ColorG,
            ColorB: this.vertices[i].ColorB,
            ColorA: this.vertices[i].ColorA,
        })
    }
    return nvs
}

func (this *PolygonData) Draw(screen *ebiten.Image, ti *TransformInfo) {
    if *this.RenderableData.tinfo.Hide {
        return
    }

    op := &ebiten.DrawTrianglesOptions{}
    op.Address = ebiten.AddressUnsafe
    indices := []uint16{}
    nverts := len(this.vertices) - 1
    for i := 0; i < nverts; i++ {
        indices = append(indices, uint16(i), uint16(i+1)%uint16(nverts), uint16(nverts))
    }

    screen.DrawTriangles(
        this.transform(ti), 
        indices, 
        asset.PolyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), 
        op,
    )
}

type StarFieldData struct {
    RenderableData
    image *ebiten.Image
}

func (this *StarFieldData) Init() {
    this.baseInit()
    this.image = nil
}

func (this *StarFieldData) Generate(w, h int) {
    this.image = ebiten.NewImage(w, h)
    color := color.RGBA{255, 255, 255, 255}

    for i := 0; i < 40; i++ {
        x := rand.Intn(w)
        y := rand.Intn(h)
        this.image.Set(x, y, color)
    }
}

func (this *StarFieldData) Animate() {
}

func (this *StarFieldData) Play(tag string) {
}

func (this *StarFieldData) Draw(screen *ebiten.Image, ti *TransformInfo) {
    if *this.RenderableData.tinfo.Hide {
        return
    }

    op := &ebiten.DrawImageOptions{}

    w, h := this.image.Size()

    newX := ti.Offset.X + this.RenderableData.tinfo.Offset.X - float64(w/2)
    newY := ti.Offset.Y + this.RenderableData.tinfo.Offset.Y - float64(h/2)

    op.GeoM.Translate(newX, newY)
    screen.DrawImage(this.image, op)
}

type BlackBarData struct {
    RenderableData
    image *ebiten.Image
}

func (this *BlackBarData) Init() {
    this.baseInit()
    this.image = nil
}

func (this *BlackBarData) Generate(w, h int) {
    this.image = ebiten.NewImage(w, h)
    this.image.Fill(color.Black)
}

func (this *BlackBarData) Animate() {
}

func (this *BlackBarData) Play(tag string) {
}

func (this *BlackBarData) Draw(screen *ebiten.Image, ti *TransformInfo) {
    if *this.RenderableData.tinfo.Hide {
        return
    }

    op := &ebiten.DrawImageOptions{}

    w, h := this.image.Size()

    newX := ti.Offset.X + this.RenderableData.tinfo.Offset.X - float64(w/2)
    newY := ti.Offset.Y + this.RenderableData.tinfo.Offset.Y - float64(h/2)

    op.GeoM.Translate(newX, newY)
    screen.DrawImage(this.image, op)
}
