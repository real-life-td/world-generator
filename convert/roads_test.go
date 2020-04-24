package convert

import (
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertRoads(t *testing.T) {
	_, err := convertRoads(nil, []*overpass.Way{})
	require.Error(t, err, "nil metadata should error")

	metadata := world.NewMetadata(100, 100, 0.0, 0.0, 1.0, 1.0)

	_, err = convertRoads(metadata, nil)
	require.Error(t, err, "nil road elements should error")

	roadElements := []*overpass.Way{
		{
			Id:     0,
			Bounds: &overpass.Bounds{MinLat: 0, MinLon: 0, MaxLat: 0.5, MaxLon: 0.5},
			Nodes:  []uint64{0, 1, 2},
			Geometry: []*overpass.LatLon{
				{0.0, 0.0},
				{0.5, 0.0},
				{0.5, 0.5},
			},
			Tags: &overpass.Tags{Highway: "primary"},
		},
		{
			Id:     1,
			Bounds: &overpass.Bounds{MinLat: 0.5, MinLon: 0.5, MaxLat: 1.0, MaxLon: 1.0},
			Nodes:  []uint64{2, 3, 4},
			Geometry: []*overpass.LatLon{
				{0.5, 0.5},
				{0.5, 1.0},
				{1.0, 1.0},
			},
			Tags: &overpass.Tags{Highway: "primary"},
		},
		{
			Id:     2,
			Bounds: &overpass.Bounds{MinLat: 0.0, MinLon: 1.0, MaxLat: 0.5, MaxLon: 1.0},
			Nodes:  []uint64{5, 3},
			Geometry: []*overpass.LatLon{
				{0, 1.0},
				{0.5, 1.0},
			},
			Tags: &overpass.Tags{Highway: "primary"},
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
	node3 := world.NewNode(makeId(3, world.NodeType), 100, 50)
	node4 := world.NewNode(makeId(4, world.NodeType), 100, 100)
	node5 := world.NewNode(makeId(5, world.NodeType), 100, 0)

	expectedRoads := []*world.Road{
		world.NewRoad(1, node0, node1, 1),
		world.NewRoad(1, node1, node2, 1),
		world.NewRoad(1, node2, node3, 1),
		world.NewRoad(1, node3, node4, 1),
		world.NewRoad(1, node5, node3, 1),
	}

	roads, err := convertRoads(metadata, roadElements)
	require.NotNil(t, roads)
	require.NoError(t, err)

	// Check that every road was given a unique id
	usedIds := make([]world.Id, 0, len(roads))
	for _, road := range roads {
		require.Equal(t, world.RoadType, road.Id().Type())
		require.NotContains(t, usedIds, road.Id(), "ids must be unique")
		usedIds = append(usedIds, road.Id())
	}

	// Replace each road with a road with the same values but with ids all equal to 1
	for i, road := range roads {
		roads[i] = world.NewRoad(1, road.Node1(), road.Node2(), road.Cost())
	}

	// With the new ids the roads should match the expected roads
	require.ElementsMatch(t, expectedRoads, roads)
}
