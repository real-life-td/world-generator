package convert

import (
	"errors"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
)

func convertRoads(metadata *world.Metadata, roadElements []*overpass.Way) (roads []*world.Road, err error) {
	if roadElements == nil {
		return nil, errors.New("roadElements cannot be nil")
	}

	toGameCoords, _, err := world.CreateConverters(metadata)
	if err != nil {
		return nil, err
	}

	roads = make([]*world.Road, 0, len(roadElements))
	placedNodes := make(map[uint64]*world.Node)

	var nextRoadId uint64 = 0

	for _, e := range roadElements {
		// Road segments will go from this node to the next in the array
		var prevNode *world.Node
		for i, nodeId := range e.Nodes {
			// Check if the node has already been created and create it if not
			p := placedNodes[nodeId]

			if p == nil {
				id, err := world.NewId(nodeId, world.NodeType)
				if err != nil {
					return nil, err
				}

				lat, lon := e.Geometry[i].Lat, e.Geometry[i].Lon
				x, y := toGameCoords(lat, lon)
				p = world.NewNode(id, x, y)
				placedNodes[nodeId] = p
			}

			// The will be no previous node to connect to on the first loop
			if i != 0 {
				id, err := world.NewId(nextRoadId, world.RoadType)
				nextRoadId++
				if err != nil {
					return nil, err
				}

				roads = append(roads, world.NewRoad(id, prevNode, p, 1))
			}

			prevNode = p
		}
	}

	return roads, nil
}
