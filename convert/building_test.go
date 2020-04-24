package convert

import (
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertBuildings(t *testing.T) {
	_, err := convertBuildings(nil, []*overpass.Way{})
	require.Error(t, err, "nil metadata should error")

	metadata := world.NewMetadata(100, 100, 0.0, 0.0, 1.0, 1.0)

	_, err = convertBuildings(metadata, nil)
	require.Error(t, err, "nil buildElement should error")

	buildingElements := []*overpass.Way{
		{
			Id:     0,
			Bounds: [4]float64{0.0, 0.0, 0.5, 0.5},
			Nodes:  []uint64{0, 1, 2, 3},
			Geometry: []*overpass.LatLon{
				{0.0, 0.0},
				{0.5, 0.0},
				{0.5, 0.5},
				{0.0, 0.5},
			},
			Tags: &overpass.Tags{Building: "yes"},
		},
		{
			Id:     1,
			Bounds: [4]float64{0.6, 0.6, 1.0, 1.0},
			Nodes:  []uint64{4, 5, 6},
			Geometry: []*overpass.LatLon{
				{0.6, 0.6},
				{0.6, 1.0},
				{1.0, 1.0},
			},
			Tags: &overpass.Tags{Building: "yes"},
		},
	}

	makeId := func(baseId uint64, idType world.Type) world.Id {
		id, err := world.NewId(baseId, idType)
		require.NoError(t, err)
		return id
	}

	node0 := world.NewNode(makeId(0, world.NodeType), 0, 0)
	node1 := world.NewNode(makeId(1, world.NodeType), 0, 50)
	node2 := world.NewNode(makeId(2, world.NodeType), 50, 50)
	node3 := world.NewNode(makeId(3, world.NodeType), 50, 0)
	node4 := world.NewNode(makeId(4, world.NodeType), 60, 60)
	node5 := world.NewNode(makeId(5, world.NodeType), 100, 60)
	node6 := world.NewNode(makeId(6, world.NodeType), 100, 100)

	expectedBuildings := []*world.Building{
		world.NewBuilding(makeId(0, world.BuildingType), []*world.Node{node0, node1, node2, node3}),
		world.NewBuilding(makeId(1, world.BuildingType), []*world.Node{node4, node5, node6}),
	}

	buildings, err := convertBuildings(metadata, buildingElements)
	require.NoError(t, err)
	require.ElementsMatch(t, buildings, expectedBuildings)
}
