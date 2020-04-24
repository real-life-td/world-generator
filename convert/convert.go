package convert

import (
	"errors"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
)

type elementType int

const (
	BuildingType elementType = iota
	HighwayType
)

func Convert(meta *world.Metadata, result *overpass.Result) (w *world.Container, err error) {
	buildingElements, highwayElements := separate(result.Elements)

	buildings, err := convertBuildings(meta, buildingElements)
	if err != nil {
		return nil, err
	}

	roads, err := convertRoads(meta, highwayElements)
	if err != nil {
		return nil, err
	}

	return world.NewContainer(meta, roads, buildings), nil
}

func separate(elements []*overpass.Way) (buildings, highways []*overpass.Way) {
	buildingElements := make([]*overpass.Way, 0, 10)
	highwayElements := make([]*overpass.Way, 0, len(elements)) // Most elements will end up being roads

	for _, e := range elements {
		t, err := classify(e)
		if err != nil {
			// Overpass data might not be perfect so just ignore this error
			continue
		}

		switch t {
		case BuildingType:
			buildingElements = append(buildingElements, e)
		case HighwayType:
			highwayElements = append(highwayElements, e)
		}
	}

	return buildingElements, highwayElements
}

func classify(e *overpass.Way) (t elementType, err error) {
	if e.Tags.Building != "" {
		if e.Tags.Highway != "" {
			return -1, errors.New("element is both building and highway type")
		}

		return BuildingType, nil
	} else if e.Tags.Highway != "" {
		return HighwayType, nil
	}

	return -1, errors.New("element does not have type")
}
