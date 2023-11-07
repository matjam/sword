package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lmittmann/tint"
	"github.com/matjam/sword/internal/mapgen"
	"github.com/mattn/go-colorable"
)

type Game struct {
	mg *mapgen.MapGenerator
}

func ConfigureLogger() {
	w := os.Stderr
	slog.SetDefault(slog.New(
		tint.NewHandler(colorable.NewColorable(w), &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func main() {
	ConfigureLogger()

	game := &Game{
		mg: mapgen.NewMapGenerator(1920/16-1, 1080/16, time.Now().UnixNano(), 250),
	}

	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("display the map!")
	if err := ebiten.RunGame(game); err != nil {
		log.Panic("failed to run game: ", err)
	}
}

func (g *Game) Update() error {
	g.mg.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mg.DrawDebug(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}
