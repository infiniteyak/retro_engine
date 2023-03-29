package shader

import (
	//"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
    "image/color"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

    //"fmt"
)

var (
	//go:embed passthrough.go
	passthrough_go []byte

	//go:embed palette_swap.go
	paletteswap_go []byte
)

var Shaders map[string]*ebiten.Shader

var Image0 *ebiten.Image
var Image1 *ebiten.Image
var ImageWidth int
var ImageHeight int

func InitShaders(w, h float64) {
    ImageWidth = int(w)
    ImageHeight = int(h)
    Shaders = map[string]*ebiten.Shader{}

    s, err := ebiten.NewShader(passthrough_go)
    if err != nil {
        panic(err)
    }
    Shaders["passthrough"] = s

    s, err = ebiten.NewShader(paletteswap_go)
    if err != nil {
        panic(err)
    }
    Shaders["paletteSwap"] = s

    // Init image 0 and 1
    Image0 = ebiten.NewImage(ImageWidth, ImageHeight)

    Image1 = ebiten.NewImage(ImageWidth, ImageHeight)
    ebitenutil.DrawRect(Image1, 0.0, 0.0, w/2.0, h, color.White)
}

func RunNoShader(screen *ebiten.Image) {
    op := &ebiten.DrawImageOptions{}
    screen.DrawImage(Image0, op)
    /*
    msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
    */
	Image0.Clear()
}

func RunPassthroughShader(screen *ebiten.Image) {
	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"ScreenSize": []float32{float32(ImageWidth), float32(ImageHeight)},
	}
	op.Images[0] = Image0
	op.Images[1] = Image1

	screen.DrawRectShader(ImageWidth, ImageHeight, Shaders["passthrough"], op)
    Image0 = ebiten.NewImage(ImageWidth, ImageHeight)
}

func RunPaletteSwapShader(screen *ebiten.Image) {
	op := &ebiten.DrawRectShaderOptions{}
    // TODO These need to be slices I think, so I can't hard code the size like
    // this. However I maybe can declare them as arrays, of my max size, then
    // convert those to slices and pass those in.
	op.Uniforms = map[string]any{
		"ScreenSize": []float32{float32(ImageWidth), float32(ImageHeight)},
		"TextureColors": []float32{
            float32(0), float32(0), float32(0),
        },
		"SourcePalette": []float32{
            float32(255), float32(85), float32(85),
        },
		"Palette0": []float32{
            float32(255), float32(255), float32(255),
        },
	}
	op.Images[0] = Image0
	op.Images[1] = Image1
	//screen.DrawRectShader(ImageWidth, ImageHeight, Shaders["passthrough"], op)
	screen.DrawRectShader(ImageWidth, ImageHeight, Shaders["paletteSwap"], op)
    Image0 = ebiten.NewImage(ImageWidth, ImageHeight)
}
