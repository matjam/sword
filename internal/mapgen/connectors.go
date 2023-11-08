package mapgen

import "github.com/matjam/sword/internal/terrain"

////////////////////////////////////////////////////////////////////////////////
// Connectors

func (mg *MapGenerator) generateConnectors() {
	// The generateConnectors() method is where we generate the connectors. We do
	// this by finding all the tiles that are adjacent to a corridor, and then
	// checking if they are adjacent to a room. If they are, we add them to the
	// list of connectors. We then shuffle the list of connectors, and then we
	// iterate over the list of connectors and try to connect them to a room.

	for y := 1; y < mg.Height-1; y += 1 {
		for x := 1; x < mg.Width-1; x += 1 {
			ok, region1, region2 := mg.isConnectorTile(x, y)
			if ok {
				connector := &Connector{
					x:       x,
					y:       y,
					region1: region1,
					region2: region2,
				}
				mg.connectorGrid.Set(x, y, connector)

				// add this connector to the list of connectors
				mg.connectors = append(mg.connectors, connector)
			}
		}
	}

	mg.Phase = PhaseConnectingRegions
}

func (mg *MapGenerator) isConnectorTile(x, y int) (isConnector bool, region1, region2 *Region) {
	// Determine if the current tile connects two different regions. We only
	// conside tiles that are rooms or corridors.

	e := mg.terrainGrid.Get(x+1, y)
	w := mg.terrainGrid.Get(x-1, y)

	if (e == terrain.Room && w == terrain.Room) ||
		(e == terrain.Corridor && w == terrain.Corridor) ||
		(e == terrain.Room && w == terrain.Corridor) ||
		(e == terrain.Corridor && w == terrain.Room) {
		eRegion := mg.regionGrid.Get(x+1, y)
		wRegion := mg.regionGrid.Get(x-1, y)
		if eRegion.id != wRegion.id {
			return true, eRegion, wRegion
		}
	}

	n := mg.terrainGrid.Get(x, y-1)
	s := mg.terrainGrid.Get(x, y+1)

	if (n == terrain.Room && s == terrain.Room) ||
		(n == terrain.Corridor && s == terrain.Corridor) ||
		(n == terrain.Room && s == terrain.Corridor) ||
		(n == terrain.Corridor && s == terrain.Room) {
		nRegion := mg.regionGrid.Get(x, y-1)
		sRegion := mg.regionGrid.Get(x, y+1)
		if nRegion.id != sRegion.id {
			return true, nRegion, sRegion
		}
	}

	return false, nil, nil
}
