package engine

import (
    "github.com/infiniteyak/retro_engine/engine/constants"
	"github.com/hajimehoshi/ebiten/v2"
    "log"
	"os"
	"os/signal"
	"syscall"
    "math/rand"
	"time"
)

func RunGame(g ebiten.Game, title string) {
    rand.Seed(time.Now().UnixNano())

    ebiten.SetWindowTitle(title)
    ebiten.SetWindowResizable(true)
    ebiten.SetMaxTPS(constants.MaxTPS)
    ebiten.SetWindowClosingHandled(true)

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigc
        os.Exit(0)
    }()

    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}
