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

	expectedRoads := []*world.Road{
		world.NewRoad(makeId(0, world.RoadType), world.NewNode(makeId(0, world.NodeType), 0, 0)),
		world.NewRoad(makeId(1, world.RoadType), world.NewNode(makeId(1, world.NodeType), 0, 50)),
		world.NewRoad(makeId(2, world.RoadType), world.NewNode(makeId(2, world.NodeType), 50, 50)),
		world.NewRoad(makeId(3, world.RoadType), world.NewNode(makeId(3, world.NodeType), 100, 50)),
		world.NewRoad(makeId(4, world.RoadType), world.NewNode(makeId(4, world.NodeType), 100, 100)),
		world.NewRoad(makeId(5, world.RoadType), world.NewNode(makeId(5, world.NodeType), 100, 0)),
	}

	expectedRoads[0].InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{expectedRoads[1]}})
	expectedRoads[1].InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{expectedRoads[0], expectedRoads[2]}})
	expectedRoads[2].InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{expectedRoads[1], expectedRoads[3]}})
	expectedRoads[3].InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{expectedRoads[2], expectedRoads[4], expectedRoads[5]}})
	expectedRoads[4].InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{expectedRoads[3]}})
	expectedRoads[5].InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{expectedRoads[3]}})


	roads, err := convertRoads(metadata, roadElements)
	require.NotNil(t, roads)
	require.NoError(t, err)

	// The order of the road connections could be different without breaking anything. This doesn't check that currently
	require.ElementsMatch(t, expectedRoads, roads)
}
