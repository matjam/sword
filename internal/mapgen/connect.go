package mapgen

import (
	"fmt"
	"image/color"
	"log/slog"

	"github.com/matjam/sword/internal/terrain"
)

////////////////////////////////////////////////////////////////////////////////
// Connecting regions

func (mg *MapGenerator) connectRegions() {
	// The connectRegions() method is where we connect all the regions together.

	// if there's only one region, we're done.
	if len(mg.regions) == 1 {
		mg.phase = PhaseRemoveDeadEnds
		return
	}

	if mg.rootRegion == nil {
		mg.selectRootRegion()
	}

	if len(mg.rootConnectors) == 0 {
		mg.findRootConnectors()

		if len(mg.rootConnectors) == 0 {
			mg.phase = PhaseRemoveDeadEnds
			return
		}

		// shuffle the list of root connectors
		shuffleArray(mg.rng, mg.rootConnectors)
	}

	// The algorithm here is simple, we work through the list of root connectors,
	// and for each one we check if it connects the root region to a region that
	// is not yet connected to the root region. If it does, we connect them and
	// remove the connector from the list of root connectors. We keep doing this
	// until we run out of regions to connect.
	success := false

	// because this function is called every update tick, we don't want to
	// try to connect all the regions at once, because that would make the
	// map generation happen in one frame. Instead, we only try to connect
	// one region per update tick.
	for !success {
		if len(mg.rootConnectors) == 0 {
			return
		}
		// grab the first root connector from the list
		c := mg.rootConnectors[0]

		// remove the root connector from the list of root connectors
		mg.rootConnectors = mg.rootConnectors[1:]

		// check if the connector connects the root region to a region that
		// is not yet connected to the root region.
		if mg.connectorIsBesideDoor(c) {
			continue
		}

		if mg.connectsRootToUnconnectedRegion(c) {
			// set the location to a door, and set the region to the root region
			mg.terrainGrid.Set(c.x, c.y, terrain.Door)
			mg.regionGrid.Set(c.x, c.y, mg.rootRegion)

			// find the region that is not the root region
			var otherRegion *Region
			if c.region1.id == mg.rootRegion.id {
				otherRegion = c.region2
			} else {
				otherRegion = c.region1
			}

			// replace all instances of the region with the root region
			mg.replaceRegion(otherRegion, mg.rootRegion)

			// remove the region from the list of unconnected regions
			delete(mg.regions, otherRegion.id)

			// success!
			success = true
		}
	}
}

func (mg *MapGenerator) connectorIsBesideDoor(c *Connector) bool {
	// check if the connector is beside a door
	e := mg.terrainGrid.Get(c.x+1, c.y)
	w := mg.terrainGrid.Get(c.x-1, c.y)
	n := mg.terrainGrid.Get(c.x, c.y-1)
	s := mg.terrainGrid.Get(c.x, c.y+1)

	if e == terrain.Door || w == terrain.Door || n == terrain.Door || s == terrain.Door {
		return true
	}

	return false
}

func (mg *MapGenerator) connectsRootToUnconnectedRegion(connector *Connector) bool {
	// check if the connector connects the root region to an unconnected region
	if connector.region1.id == mg.rootRegion.id && connector.region2.id != mg.rootRegion.id {
		return true
	}

	if connector.region2.id == mg.rootRegion.id && connector.region1.id != mg.rootRegion.id {
		return true
	}

	return false
}

func (mg *MapGenerator) selectRootRegion() {
	slog.Info(fmt.Sprintf("there are %d regions", len(mg.regions)))
	slog.Info(fmt.Sprintf("there are %v rooms", len(mg.roomList)))

	// all rooms start out as unconnected
	for _, room := range mg.roomList {
		mg.unconnectedRooms = append(mg.unconnectedRooms, room)
	}

	// shuffle the unconnected regions
	shuffleArray(mg.rng, mg.unconnectedRooms)

	// grab the last room from the list
	rootRoom := mg.unconnectedRooms[len(mg.unconnectedRooms)-1]
	mg.unconnectedRooms = mg.unconnectedRooms[:len(mg.unconnectedRooms)-1]
	mg.rootRegion = rootRoom.Region

	// set the color of the root region to black
	mg.rootRegion.clr = color.RGBA{0x00, 0x00, 0x00, 0xff}

	slog.Info(fmt.Sprintf("room at %v,%v selected as root region", rootRoom.X, rootRoom.Y))
}

func (mg *MapGenerator) findRootConnectors() {
	// shuffle the list of connectors
	shuffleArray(mg.rng, mg.connectors)

	otherConnectors := make([]*Connector, 0)
	mg.rootConnectors = make([]*Connector, 0)

	// find all the connectors that connect the root region to another region
	for _, c := range mg.connectors {
		if (c.region1.id == mg.rootRegion.id && c.region2.id != mg.rootRegion.id) ||
			(c.region1.id != mg.rootRegion.id && c.region2.id == mg.rootRegion.id) {
			mg.rootConnectors = append(mg.rootConnectors, c)
		} else {
			otherConnectors = append(otherConnectors, c)
		}
	}

	shuffleArray(mg.rng, mg.rootConnectors)

	mg.connectors = otherConnectors
}

func (mg *MapGenerator) replaceRegion(oldRegion *Region, newRegion *Region) {
	// The replaceRegion() method is where we replace all instances of one region
	// with another region. We do this by iterating over the Grid, and replacing
	// all instances of the old region with the new region.

	for y := 0; y < mg.Height; y++ {
		for x := 0; x < mg.Width; x++ {
			r := mg.regionGrid.Get(x, y)
			if r != nil && r.id == oldRegion.id {
				mg.regionGrid.Set(x, y, newRegion)
			}

			c := mg.connectorGrid.Get(x, y)
			if c != nil {
				if c.region1.id == oldRegion.id {
					c.region1 = newRegion
				}
				if c.region2.id == oldRegion.id {
					c.region2 = newRegion
				}
			}
		}
	}
}
