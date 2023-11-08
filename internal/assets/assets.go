package assets

import (
	"image"
	"image/color"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/matjam/sword/internal/config"
	"github.com/matjam/sword/internal/tileset"
	woff "github.com/tdewolff/canvas/font"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const dpi = 72

var globalAssetManager *AssetManager

type AssetManager struct {
	images    map[string]image.Image
	tiles     map[string][]*ebiten.Image
	fonts     map[string]font.Face
	fontSizes map[string]int
	tileSet   map[string]*tileset.Tileset
}

type fontConfig struct {
	Path string  `json:"path"`
	Size float64 `json:"size"`
}

type assetConfig struct {
	Images map[string]string     `json:"images"`
	Fonts  map[string]fontConfig `json:"fonts"`
}

func StartAssetManager(configPath string) {
	if globalAssetManager != nil {
		slog.Error("asset manager already started")
		return
	}

	m := AssetManager{
		images:    make(map[string]image.Image),
		tiles:     make(map[string][]*ebiten.Image),
		fonts:     make(map[string]font.Face),
		fontSizes: make(map[string]int),
		tileSet:   make(map[string]*tileset.Tileset),
	}

	assetConfig := config.Load().Assets

	// load images
	for name, path := range assetConfig.Images {
		m.images[name] = m.loadImage(path, name)
	}

	// load fonts
	for name, fontConfig := range assetConfig.Fonts {
		m.loadFont(fontConfig.Path, name, fontConfig.Size)
		m.images[name] = m.CreateTilesheet(name, int(fontConfig.Size))
	}

	// load tilesets
	for name, tilesetConfig := range assetConfig.Tilesets {
		atlas := m.loadImage(tilesetConfig.Path, name)

		m.tileSet[name] = tileset.Load(name,
			atlas,
			tilesetConfig.TileSize,
			tilesetConfig.Columns,
			tilesetConfig.Rows,
			tilesetConfig.Autotiles,
			tilesetConfig.Fixtures)
	}
	globalAssetManager = &m
}

func (am *AssetManager) loadImage(path string, name string) *ebiten.Image {
	reader, err := os.Open(path)
	if err != nil {
		slog.Error("error opening image", err)
		panic(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		slog.Error("error decoding image", err)
		panic(err)
	}

	img := ebiten.NewImageFromImage(m)

	slog.Info("image loaded", "name", name, "path", path)

	return img
}

func (am *AssetManager) loadFont(fontPath string, name string, size float64) {
	var err error
	var data []byte
	var fnt *sfnt.Font
	var fntData []byte

	data, err = os.ReadFile(fontPath)
	if err != nil {
		slog.Error("error reading font file", err)
		panic(err)
	}

	ext := path.Ext(fontPath)
	switch strings.ToLower(ext) {
	case ".ttf":
		fnt, err = opentype.Parse(data)
		if err != nil {
			slog.Error("error parsing ttf font", err)
			panic(err)
		}
	case ".woff":
		fntData, err = woff.ParseWOFF(data)
		if err != nil {
			slog.Error("error parsing woff font", err)
			panic(err)
		}
		fnt, err = sfnt.Parse(fntData)
	case ".woff2":
		fntData, err = woff.ParseWOFF2(data)
		if err != nil {
			slog.Error("error parsing woff2 font", err)
			panic(err)
		}
		fnt, err = sfnt.Parse(fntData)
	}

	if err != nil {
		slog.Error("error parsing font", err)
		panic(err)
	}

	f, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		slog.Error("error creating font face", err)
		panic(err)
	}

	am.fonts[name] = f
	am.fontSizes[name] = int(size)

	slog.Info("font loaded", "name", name, "fontPath", fontPath)
}

// CreateTilesheet creates a 16x16 tilesheet from the given font, with
// each character being pixelSize x pixelSize.
func (am *AssetManager) CreateTilesheet(fontName string, pixelSize int) image.Image {
	face := am.fonts[fontName]
	size := am.fontSizes[fontName]

	// create the tilesheet
	tilesheet := ebiten.NewImage(16*pixelSize, 16*pixelSize)

	offset := 0
	// draw each character to the tilesheet
	for i := 32; i < 128; i++ {
		x := (offset % 16) * pixelSize
		y := (offset / 16) * pixelSize

		char := string([]rune{rune(i)})
		text.Draw(tilesheet, char, face, x, y+size, color.White)
		offset++
	}

	for i := 129792; i < 129792+128; i++ {
		x := (offset % 16) * pixelSize
		y := (offset / 16) * pixelSize

		char := string([]rune{rune(i)})
		text.Draw(tilesheet, char, face, x, y+size, color.White)
		offset++
	}

	return tilesheet
}

func (am *AssetManager) GetImage(name string) image.Image {
	return am.images[name]
}

func (am *AssetManager) GetFont(name string) font.Face {
	return am.fonts[name]
}

func (am *AssetManager) GetFontSize(name string) int {
	return am.fontSizes[name]
}

func GetFont(name string) font.Face {
	return globalAssetManager.GetFont(name)
}

func GetFontSize(name string) int {
	return globalAssetManager.GetFontSize(name)
}

func GetImage(name string) image.Image {
	return globalAssetManager.GetImage(name)
}

func GetTileset(name string) *tileset.Tileset {
	return globalAssetManager.tileSet[name]
}
