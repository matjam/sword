package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lmittmann/tint"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/tilemap"
	"github.com/matjam/sword/internal/tilemap/text"
	"github.com/mattn/go-colorable"

	_ "image/png"
)

type Game struct {
	tm         *tilemap.Grid
	tmRenderer tilemap.Renderer
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.tmRenderer.Draw(screen, 28, 26,
		tilemap.Rectangle{
			X:      0,
			Y:      0,
			Width:  77,
			Height: 49,
		})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 768
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

	game := &Game{}

	slog.Info("loading assets ...")
	assets.StartAssetManager()

	slog.Info("creating tilemap ...")
	game.tm = tilemap.NewGrid(600, 400)

	// lets clear out a room

	for y := 5; y < 15; y++ {
		for x := 5; x < 15; x++ {
			game.tm.SetTile(x, y, &tilemap.Tile{
				Type: tilemap.TileTypeFloor,
			})
		}
	}

	// and a door
	game.tm.SetTile(10, 5, &tilemap.Tile{
		Type: tilemap.TileTypeClosedDoor,
	})

	game.tm.SetTile(0, 0, &tilemap.Tile{
		Type: tilemap.TileTypeFloor,
	})

	game.tmRenderer = text.NewRenderer(game.tm, "square")

	ebiten.SetWindowSize(1280, 768)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Panic("failed to run game: ", err)
	}
}
