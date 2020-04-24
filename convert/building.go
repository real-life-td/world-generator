package convert

import (
	"errors"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
)

func convertBuildings(metadata *world.Metadata, buildingElements []*overpass.Way) (b []*world.Building, err error) {
	if buildingElements == nil {
		return nil, errors.New("building elements be nil")
	}

	toGameCoords, _, err := world.CreateConverters(metadata)
	if err != nil {
		return nil, err
	}

	buildings := make([]*world.Building, 0, len(buildingElements))

	for _, e := range buildingElements {
		id, err := world.NewId(e.Id, world.BuildingType)
		if err != nil {
			return nil, err
		}

		points := make([]*world.Node, 0, len(e.Nodes) - 1)
		for i, nodeId := range e.Nodes {
			id, err := world.NewId(nodeId, world.NodeType)
			if err != nil {
				return nil, err
			}

			x, y := toGameCoords(e.Geometry[i].Lat, e.Geometry[i].Lon)
			points = append(points, world.NewNode(id, x, y))
		}

		buildings = append(buildings, world.NewBuilding(id, points))
	}

	return buildings, nil
}