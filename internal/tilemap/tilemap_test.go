package tilemap_test

import (
	"testing"

	"github.com/matjam/sword/internal/tilemap"
)

func TestNewTileMap(t *testing.T) {
	tm := tilemap.NewTileMap(10, 10)
	if tm.Width != 10 {
		t.Errorf("expected width to be 10, got %d", tm.Width)
	}
	if tm.Height != 10 {
		t.Errorf("expected height to be 10, got %d", tm.Height)
	}
	if len(tm.Tiles) != 100 {
		t.Errorf("expected length of tiles to be 100, got %d", len(tm.Tiles))
	}
}

func TestGetTile(t *testing.T) {
	tm := tilemap.NewTileMap(10, 10)
	tile := tm.GetTile(0, 0)
	if tile == nil {
		t.Errorf("expected tile to not be nil")
	}
	tile = tm.GetTile(10, 10)
	if tile != nil {
		t.Errorf("expected tile to be nil")
	}
}

func TestSetTile(t *testing.T) {
	tm := tilemap.NewTileMap(10, 10)
	tile := tilemap.Tile{
		Type: tilemap.TileTypeFloor,
	}
	tm.SetTile(0, 0, &tile)
	tile = *tm.GetTile(0, 0)
	if tile.Type != tilemap.TileTypeFloor {
		t.Errorf("expected tile type to be floor, got %s", tile.Type)
	}
}

func TestIsVisible(t *testing.T) {
	tm := tilemap.NewTileMap(10, 10)
	tile := tilemap.Tile{
		Type: tilemap.TileTypeFloor,
	}
	tm.SetTile(0, 0, &tile)
	tile = tilemap.Tile{
		Type: tilemap.TileTypeFloor,
	}
	tm.SetTile(1, 0, &tile)

	tm.Dump()

	if !tm.IsVisible(0, 0, 1, 0) {
		t.Errorf("expected tile to be visible")
	}
	if tm.IsVisible(0, 0, 2, 0) {
		t.Errorf("expected tile to not be visible")
	}
}