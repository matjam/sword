package assets

import (
	"encoding/json"
	"image"
	"log"
	"os"

	"github.com/gookit/slog"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const dpi = 72

type AssetManager struct {
	images map[string]*image.Image
	tiles  map[string][]*ebiten.Image
	fonts map[string]*font.Face
}

type fontConfig struct {
	Path string `json:"path"`
	Size float64 `json:"size"`
}

type assetConfig struct {
	Images map[string]string `json:"images"`
	Fonts map[string]fontConfig `json:"fonts"`
}

func NewAssetManager() *AssetManager {
	m := AssetManager{
		images: make(map[string]*image.Image),
		tiles:  make(map[string][]*ebiten.Image),
		fonts: make(map[string]*font.Face),
	}

	// load config
	data, err := os.ReadFile("assets.json")
	if err != nil {
		log.Fatal(err)
	}
	config := assetConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	
	// load images
	for name, path := range config.Images {
		m.loadImage(path, name)
	}
	
	// load fonts
	for name, fontConfig := range config.Fonts {
		m.loadFont(fontConfig.Path, name, fontConfig.Size)
	}

	return &m
}

func (am *AssetManager) loadImage(path string, name string) {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	am.images[name] = &m

	slog.Infof("image: loaded %v:%v", name, path)
}

func (am *AssetManager) loadFont(path string, name string, size float64) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	f, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	am.fonts[name] = &f

	slog.Infof("font: loaded %v:%v", name, path)
}

func (am *AssetManager) GetImage(name string) *image.Image {
	return am.images[name]
}

func (am *AssetManager) GetFont(name string) *font.Face {
	return am.fonts[name]
}