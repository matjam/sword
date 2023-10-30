package main

import (
	"github.com/gookit/slog"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/tilemap"

	"image/color"
	_ "image/png"
)

type Game struct {
	assetManager *assets.AssetManager
	tileMap *tilemap.TileMap
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, "██ Hello, World! ██", *g.assetManager.GetFont("square"), 40, 40, color.White)

	g.tileMap.Draw(screen)
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
	game.assetManager = assets.NewAssetManager()

	slog.Info("creating tilemap ...")
	game.tileMap = tilemap.NewTileMap(200, 120).WithFont(game.assetManager.GetFont("square"))

	ebiten.SetWindowSize(1280, 768)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		slog.Fatal(err)
	}
}
