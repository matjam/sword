package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lmittmann/tint"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/ecs"
	"github.com/matjam/sword/internal/ecs/component"
	"github.com/matjam/sword/internal/ecs/entity"
	"github.com/matjam/sword/internal/ecs/system"
	"github.com/matjam/sword/internal/tilemap"
	"github.com/matjam/sword/internal/tilemap/text"
	"github.com/mattn/go-colorable"

	_ "image/png"
	_ "net/http/pprof"
)

type Game struct {
	tm         *tilemap.Grid
	tmRenderer tilemap.Renderer
	world      *ecs.World
}

func (g *Game) Update() error {
	g.world.Update(time.Second / 60)

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

	g.world.Draw(screen)
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

func ConfigureWorld() *ecs.World {
	world := ecs.NewWorld()

	inputSystem := &system.Input{}

	world.AddSystem(inputSystem)
	world.AddSystem(&system.Movement{})
	world.AddSystem(&system.Renderer{GridSize: assets.GetFontSize("square")})

	player := world.AddEntity(&entity.Player{})
	playerLocation := ecs.GetComponent[*component.Location](world, player)
	playerLocation.X = 7
	playerLocation.Y = 7

	inputSystem.Player = player

	return world
}

func main() {
	ConfigureLogger()

	// go func() {
	// 	err := http.ListenAndServe("localhost:6060", nil)
	// 	if err != nil {
	// 		slog.Error("error running pprof server", err)
	// 	}
	// }()

	game := &Game{}

	slog.Info("loading assets ...")
	assets.StartAssetManager()

	slog.Info("creating tilemap ...")
	game.tm = tilemap.NewGrid(600, 400)

	slog.Info("creating world ...")
	game.world = ConfigureWorld()

	// lets clear out a room

	for y := 5; y < 35; y++ {
		for x := 5; x < 60; x++ {
			game.tm.SetTile(x, y, &tilemap.Tile{
				Type: tilemap.TileTypeFloor,
			})
		}
	}

	game.tmRenderer = text.NewRenderer(game.tm, "square")

	ebiten.SetWindowSize(1280, 768)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Panic("failed to run game: ", err)
	}
}
