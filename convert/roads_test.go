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
			Id:       0,
			Bounds:   [4]float64{0, 0, 0.5, 0.5},
			Nodes:    []uint64{0, 1, 2},
			Geometry: []*overpass.LatLon{
				{0.0, 0.0},
				{0.5, 0.0},
				{0.5, 0.5},
			},
			Tags:     &overpass.Tags{Highway: "primary"},
		},
		{
			Id:       1,
			Bounds:   [4]float64{0.5, 0.5, 1.0, 1.0},
			Nodes:    []uint64{2, 3, 4},
			Geometry: []*overpass.LatLon{
				{0.5, 0.5},
				{0.5, 1.0},
				{1.0, 1.0},
			},
			Tags:     &overpass.Tags{Highway: "primary"},
		},
		{
			Id:       2,
			Bounds:   [4]float64{0.0, 1.0, 0.5, 1.0},
			Nodes:    []uint64{5, 3},
			Geometry: []*overpass.LatLon{
				{0, 1.0},
				{0.5, 1.0},
			},
			Tags:     &overpass.Tags{Highway: "primary"},
		},
	}

	node0 := world.NewNode(0, 0,0)
	node1 := world.NewNode(1, 0, 50)
	node2 := world.NewNode(2, 50, 50)
	node3 := world.NewNode(3, 100, 50)
	node4 := world.NewNode(4, 100, 100)
	node5 := world.NewNode(5, 100, 0)

	expectedRoads := []*world.Road{
		world.NewRoad(1, node0, node1, 1),
		world.NewRoad(1, node1, node2, 1),
		world.NewRoad(1, node2, node3, 1),
		world.NewRoad(1, node3, node4, 1),
		world.NewRoad(1, node5, node3, 1),
	}

	roadsEqual := func(r1, r2 *world.Road) bool {
		if r1.Cost() != r2.Cost() {
			return false
		}

		nodesEqual := func(n1, n2 *world.Node) bool {
			return n1.X() == n2.X() && n1.Y() == n2.Y()
		}

		// Nodes can be in either order
		if nodesEqual(r1.Node1(), r2.Node1()) {
			return nodesEqual(r1.Node2(), r2.Node2())
		} else {
			return nodesEqual(r1.Node1(), r2.Node2()) && nodesEqual(r1.Node2(), r2.Node1())
		}
	}

	roads, err := convertRoads(metadata, roadElements)
	require.NotNil(t, roads)
	require.NoError(t, err)

	require.Equal(t, len(expectedRoads), len(roads))
	for _, road := range roads {
		matchFound := false
		for _, expectedRoad := range expectedRoads {
			if roadsEqual(expectedRoad, road) {
				matchFound = true
				break
			}
		}

		require.True(t, matchFound, "No match found for road")
	}

	usedIds := make([]world.Id, 0, len(roads))
	for _, road := range roads {
		require.Equal(t, world.RoadType, road.Id().Type())
		require.NotContains(t, usedIds, road.Id(), "ids must be unique")
		usedIds = append(usedIds, road.Id())
	}
}
