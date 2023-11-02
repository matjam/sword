package assets

import (
	"encoding/json"
	"image"
	"log"
	"os"
	"path"
	"strings"

	"github.com/gookit/slog"
	"github.com/hajimehoshi/ebiten/v2"
	woff "github.com/tdewolff/canvas/font"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const dpi = 72

type AssetManager struct {
	images map[string]*image.Image
	tiles  map[string][]*ebiten.Image
	fonts  map[string]*font.Face
}

type fontConfig struct {
	Path string  `json:"path"`
	Size float64 `json:"size"`
}

type assetConfig struct {
	Images map[string]string     `json:"images"`
	Fonts  map[string]fontConfig `json:"fonts"`
}

func NewAssetManager() *AssetManager {
	m := AssetManager{
		images: make(map[string]*image.Image),
		tiles:  make(map[string][]*ebiten.Image),
		fonts:  make(map[string]*font.Face),
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

func (am *AssetManager) loadFont(fontPath string, name string, size float64) {
	var err error
	var data []byte
	var fnt *sfnt.Font
	var fntData []byte

	data, err = os.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}

	ext := path.Ext(fontPath)
	switch strings.ToLower(ext) {
	case ".ttf":
		fnt, err = opentype.Parse(data)
	case ".woff":
		fntData, err = woff.ParseWOFF(data)
		if err != nil {
			log.Fatal(err)
		}
		fnt, err = sfnt.Parse(fntData)
	case ".woff2":
		fntData, err = woff.ParseWOFF2(data)
		if err != nil {
			log.Fatal(err)
		}
		fnt, err = sfnt.Parse(fntData)
	}

	if err != nil {
		log.Fatal(err)
	}

	f, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	am.fonts[name] = &f

	slog.Infof("font: loaded %v:%v", name, fontPath)
}

func (am *AssetManager) GetImage(name string) *image.Image {
	return am.images[name]
}

func (am *AssetManager) GetFont(name string) *font.Face {
	return am.fonts[name]
}
