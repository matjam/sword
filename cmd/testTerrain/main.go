package main

import (
	"image"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lmittmann/tint"
	"github.com/matjam/sword/internal/assets"
	"github.com/matjam/sword/internal/mapgen"
	"github.com/matjam/sword/internal/terrain"
	"github.com/matjam/sword/internal/tileset"
	"github.com/mattn/go-colorable"

	_ "image/png"
)

type Game struct {
	mg          *mapgen.MapGenerator
	pressedKeys []ebiten.Key

	mapgenDone  bool
	renderDebug bool

	Terrain *terrain.Terrain
	Tileset *tileset.Tileset

	mouseX int
	mouseY int

	viewportX int
	viewportY int
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

	assets.StartAssetManager("assets.json")

	game := &Game{
		mg: mapgen.NewMapGenerator(1920/16-1, 1080/16, time.Now().UnixNano(), 1000),
	}

	game.Tileset = assets.GetTileset("rogue_environment")

	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("display the map!")
	if err := ebiten.RunGame(game); err != nil {
		log.Panic("failed to run game: ", err)
	}
}

func (g *Game) Update() error {
	if !g.mapgenDone {
		g.mg.Update()
		g.mapgenDone = g.mg.Phase == mapgen.PhaseDone
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.mouseX, g.mouseY = ebiten.CursorPosition()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// get the distance between the mouse and the last mouse position
		dx := g.mouseX - x
		dy := g.mouseY - y

		// scroll the viewport by the distance
		g.viewportX += dx
		g.viewportY += dy
	}

	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	if len(g.pressedKeys) == 0 {
		return nil
	}

	key := g.pressedKeys[0]
	g.pressedKeys = g.pressedKeys[1:]

	switch key {
	case ebiten.KeyEscape:
		return ebiten.Termination
	case ebiten.KeyF1:
		if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
			g.renderDebug = !g.renderDebug
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.renderDebug {
		g.mg.DrawDebug(screen)
	} else {
		g.Tileset.Render(g.mg.Terrain(), screen, g.viewportX, g.viewportY, image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 640, Y: 360}}, 3)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}
