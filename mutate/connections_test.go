package mutate

import (
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/math/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitBuildingConnections(t *testing.T) {
	r1 := world.NewRoad(0, world.NewNode(1, 30, 0))
	r2 := world.NewRoad(2, world.NewNode(3, 40, 0))
	r3 := world.NewRoad(4, world.NewNode(5, 50, 0))
	r4 := world.NewRoad(6, world.NewNode(7, 60, 0))
	r5 := world.NewRoad(8, world.NewNode(9, 70, 0))
	r6 := world.NewRoad(10, world.NewNode(11, 80, 0))
	r7 := world.NewRoad(12, world.NewNode(13, 90, 0))
	r8 := world.NewRoad(14, world.NewNode(15, 100, 0))
	r9 := world.NewRoad(16, world.NewNode(17, 110, 0))
	r10 := world.NewRoad(18, world.NewNode(19, 120, 0))
	r11 := world.NewRoad(20, world.NewNode(21, 130, 0))
	r12 := world.NewRoad(22, world.NewNode(23, 0, 60))
	r13 := world.NewRoad(24, world.NewNode(25, 0, 50))

	r1.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r2}})
	r2.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r1, r3}})
	r3.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r2, r4}})
	r4.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r3, r5}})
	r5.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r4, r6}})
	r6.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r5, r7}})
	r7.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r6, r8}})
	r8.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r7, r9}})
	r9.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r8, r10}})
	r10.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r9, r11}})
	r11.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r10}})
	r12.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r13}})
	r13.InitOperation(&world.RoadInitOperation{NewConnections: []*world.Road{r12}})

	b := world.NewBuilding(26, []*world.Node{
		world.NewNode(27, 0, 0),
		world.NewNode(28, 20, 0),
		world.NewNode(29, 20, 20),
		world.NewNode(30, 0, 20),
	})

	container := world.NewContainer(nil, []*world.Road{r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12, r13}, []*world.Building{b})

	InitBuildingConnections(container, 2.0)
	expectedConnections := []*world.Connection{
		world.NewConnection(r1, 10.0, primitives.NewPoint(20, 0)),
		world.NewConnection(r13, 30.0, primitives.NewPoint(0, 20)),
	}
	require.ElementsMatch(t, expectedConnections, b.Connections())

	InitBuildingConnections(container, 3.0)
	require.Equal(t, len(b.Connections()), 3)
}
