package main

import (
	"github.com/gookit/slog"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/tilemap"
	"github.com/matjam/sword/internal/tilemap/text"

	_ "image/png"
)

type Game struct {
	tm         *tilemap.Tilemap
	tmRenderer tilemap.Renderer
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.tmRenderer.Draw(screen, 20, 40,
		tilemap.Rectangle{
			X:      0,
			Y:      0,
			Width:  40,
			Height: 20,
		})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 768
}

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})

	game := &Game{}

	slog.Info("loading assets ...")
	assets.StartAssetManager()

	slog.Info("creating tilemap ...")
	game.tm = tilemap.NewTilemap(200, 120)

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
		slog.Fatal(err)
	}
}
