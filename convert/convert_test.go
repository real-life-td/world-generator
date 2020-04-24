package convert

import (
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvert(t *testing.T) {
	metadata := world.NewMetadata(100, 100, 0.0, 0.0, 1.0, 1.0)

	result := overpass.Result{
		Elements: []*overpass.Way{
			{
				Id:     0,
				Bounds: [4]float64{0.0, 0.0, 0.5, 0.5},
				Nodes:  []uint64{0, 1},
				Geometry: []*overpass.LatLon{
					{0.0, 0.0},
					{0.5, 0.5},
				},
				Tags: &overpass.Tags{Highway: "primary"},
			},
			{
				Id:     1,
				Bounds: [4]float64{0.6, 0.6, 1.0, 1.0},
				Nodes:  []uint64{2, 3, 4, 5},
				Geometry: []*overpass.LatLon{
					{0.6, 0.6},
					{1.0, 0.6},
					{1.0, 1.0},
					{0.6, 1.0},
				},
				Tags: &overpass.Tags{Building: "yes"},
			},
		},
	}

	makeId := func(baseId uint64, idType world.Type) world.Id {
		id, err := world.NewId(baseId, idType)
		require.NoError(t, err)
		return id
	}

	expectedRoads := []*world.Road{
		world.NewRoad(
			makeId(0, world.RoadType),
			world.NewNode(makeId(0, world.NodeType), 0, 0),
			world.NewNode(makeId(1, world.NodeType), 50, 50),
			1),
	}

	expectedBuildings := []*world.Building{
		world.NewBuilding(makeId(1, world.BuildingType), []*world.Node{
			world.NewNode(makeId(2, world.NodeType), 60, 60),
			world.NewNode(makeId(3, world.NodeType), 60, 100),
			world.NewNode(makeId(4, world.NodeType), 100, 100),
			world.NewNode(makeId(5, world.NodeType), 100, 60),
		}),
	}

	expectedContainer := world.NewContainer(metadata, expectedRoads, expectedBuildings)

	container, err := Convert(metadata, &result)
	require.NoError(t, err)
	require.Equal(t, expectedContainer, container)
}

func TestClassify(t *testing.T) {
	test := func(tags overpass.Tags, expected elementType, shouldError bool) {
		w := overpass.Way{
			Id:       0,
			Bounds:   [4]float64{0, 0, 0, 0},
			Nodes:    nil,
			Geometry: nil,
			Tags:     &tags,
		}

		r, err := classify(&w)
		if shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, expected, r)
		}
	}

	test(overpass.Tags{Building: "", Highway: ""}, -1, true)
	test(overpass.Tags{Building: "yes", Highway: ""}, BuildingType, false)
	test(overpass.Tags{Building: "", Highway: "primary"}, HighwayType, false)
	test(overpass.Tags{Building: "yes", Highway: "primary"}, -1, true)
}
