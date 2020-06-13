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
	placedRoads := make(map[uint64]*world.Road)

	for _, e := range roadElements {
		// Road segments will go from this node to the next in the array
		var prevRoad *world.Road
		for i, nodeId := range e.Nodes {
			// check if a road has already been placed for the node
			r := placedRoads[nodeId]

			if r == nil {
				roadNodeId, err := world.NewId(nodeId, world.NodeType)
				if err != nil {
					return nil, err
				}

				roadId, err := world.NewId(nodeId, world.RoadType)
				if err != nil {
					return nil, err
				}

				lat, lon := e.Geometry[i].Lat, e.Geometry[i].Lon
				x, y := toGameCoords(lat, lon)

				r = world.NewRoad(roadId, world.NewNode(roadNodeId, x, y))
				placedRoads[nodeId] = r
				roads = append(roads, r)
			}

			// The will be no previous road to connect to on the first loop
			if prevRoad != nil {
				prevRoad.InitOperation(&world.RoadInitOperation{
					AdditionalConnections: []*world.Road{r},
				})

				r.InitOperation(&world.RoadInitOperation{
					AdditionalConnections: []*world.Road{prevRoad},
				})
			}

			prevRoad = r
		}
	}

	return roads, nil
}
