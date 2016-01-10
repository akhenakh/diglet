package mbt

import (
	"fmt"
	ts "github.com/buckhx/diglet/mbt/tile_system"
)

// Split features up by their tile coordinates. This is intended to be done at the deepest desired zoom level
// If a feature has any point in a tile, it will bind to that tile. A feature can be in multiple tiles
func splitFeatures(features <-chan *Feature, zoom uint) (tiles map[ts.Tile][]*Feature) {
	tiles = make(map[ts.Tile][]*Feature)
	for feature := range features {
		c := feature.Center()
		tile, _ := ts.CoordinateToTile(c.Lat, c.Lon, zoom)
		tiles[tile] = append(tiles[tile], feature)
	}
	return
}
