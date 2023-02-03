package game

import (
	"os"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    w int
    h int
}

func NewGame(width, height float64) *Game {
    this := &Game{
        w: int(width),
        h: int(height),
    }

    return this
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return this.w, this.h
}

func (this *Game) Update() error {
    if ebiten.IsWindowBeingClosed() {
        this.Exit()
        return nil
    }
	return nil
}

func (this *Game) Draw(screen *ebiten.Image) {
}
