package mutate

import (
	"container/heap"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/math/primitives"
	"github.com/real-life-td/math/raycast"
	"math"
)

type toRemove struct {
	building   *world.Building
	connection *world.Connection
	score      float64
}

type toRemoveMaxHeap []toRemove

func (r toRemoveMaxHeap) Len() int           { return len(r) }
func (r toRemoveMaxHeap) Less(i, j int) bool { return r[i].score > r[j].score }
func (r toRemoveMaxHeap) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *toRemoveMaxHeap) Push(x interface{}) { *r = append(*r, x.(toRemove)) }
func (r toRemoveMaxHeap) Peek() toRemove      { return r[len(r)-1] }

func (r *toRemoveMaxHeap) Pop() interface{} {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}

func InitBuildingConnections(container *world.Container, targetAvgConnections float64) {
	initBuildingConnections(container, 4, targetAvgConnections)
}

func initBuildingConnections(container *world.Container, initialConnectDistance int, targetAvgConnections float64) {
	for _, b := range container.Buildings() {
		connectionBounds := b.Bounds().Expand(initialConnectDistance, initialConnectDistance, initialConnectDistance, initialConnectDistance)
		withinBounds := withinBounds(connectionBounds, container.Roads())

		closeEnough := make(map[world.Id]*world.Connection) // Map to make the next step more efficient
		for _, r := range withinBounds {
			closest := closestConnection(b, r)
			if int(math.Round(closest.Distance())) < initialConnectDistance {
				closeEnough[r.Id()] = closest
			}
		}

		b.InitOperation(&world.BuildingInitOperation{
			NewConnections: cullPaths(closeEnough),
		})
	}

	// If the average number of connections for each building is low then recursively call this method with a connection
	// distance that is twice as large.
	avgConnections := averageNumberOfConnections(container.Buildings())
	if avgConnections < targetAvgConnections {
		initBuildingConnections(container, initialConnectDistance*2, targetAvgConnections)
		return
	} else {
		// Remove any extra connection points
		numToRemove := int(math.Round((avgConnections - targetAvgConnections) * float64(len(container.Buildings()))))
		cullWorstScoring(container.Buildings(), numToRemove)
	}
}

func withinBounds(bounds *primitives.Rectangle, roads []*world.Road) []*world.Road {
	possibleConnections := make([]*world.Road, 0)
	for _, r := range roads {
		if bounds.ContainsPoint(r.Node.Point) {
			possibleConnections = append(possibleConnections, r)
		}
	}

	return possibleConnections
}

func closestConnection(building *world.Building, r *world.Road) *world.Connection {
	closestDistance := math.MaxFloat64
	var closestPoint *primitives.Point

	for i, bPoint := range building.Points() {
		var bNextPoint *world.Node
		if i == len(building.Points())-1 {
			// In the last iteration the next point is at the start of the slice
			bNextPoint = building.Points()[0]
		} else {
			bNextPoint = building.Points()[i+1]
		}

		p, distance := raycast.ClosestPointTo(bPoint.Point, bNextPoint.Point, r.Point)

		if distance < closestDistance {
			closestDistance = distance
			closestPoint = p
		}
	}

	return world.NewConnection(r, closestDistance, closestPoint)
}

// Iterate through the map in it's arbitrary order. Returns only the connections that are closer than any other
// connection within 5 roads. This means that no two connections will be within 5 roads of another. Since map ordering
// is arbitrary the output of this function is generally not consistent.
func cullPaths(closeEnough map[world.Id]*world.Connection) []*world.Connection {
	passing := make([]*world.Connection, 0)

	visited := make(map[world.Id]bool)

	var breadthFirst func(cur *world.Road, distanceLimit int) *world.Connection
	breadthFirst = func(cur *world.Road, distanceLimit int) *world.Connection {
		closest := closeEnough[cur.Id()] // might be nil if the current road isn't connected to a building
		if distanceLimit == 0 {
			return closest
		}

		for _, connected := range cur.Connections() {
			if !visited[connected.Id()] {
				visited[connected.Id()] = true

				potentialClosest := breadthFirst(connected, distanceLimit-1)
				if potentialClosest == nil {
					continue
				}

				if closest == nil || potentialClosest.Distance() < closest.Distance() {
					closest = potentialClosest
				}
			}
		}

		return closest
	}

	for _, c := range closeEnough {
		if !visited[c.Road().Id()] {
			passing = append(passing, breadthFirst(c.Road(), 5))
		}
	}

	return passing
}

func cullWorstScoring(buildings []*world.Building, numToRemove int) {
	if numToRemove <= 0 {
		return
	}

	removeHeap := new(toRemoveMaxHeap)
	heap.Init(removeHeap)

	for _, b := range buildings {
		for i, score := range scoreConnections(b) {
			// First iteration
			if removeHeap.Len() == 0 {
				heap.Push(removeHeap, toRemove{building: b, score: score, connection: b.Connections()[i]})
			}

			if removeHeap.Peek().score < score {
				if removeHeap.Len() > numToRemove {
					heap.Pop(removeHeap)
				}

				heap.Push(removeHeap, toRemove{building: b, score: score, connection: b.Connections()[i]})
			}
		}
	}

	for removeHeap.Len() > 0 {
		cur := heap.Pop(removeHeap).(toRemove)
		cur.building.InitOperation(&world.BuildingInitOperation{ToRemoveConnections: []*world.Connection{cur.connection}})
	}
}

func scoreConnections(b *world.Building) []float64 {
	if len(b.Connections()) == 0 {
		return []float64{}
	}

	// Find the closest connection to the building
	var closest = b.Connections()[0]
	for _, c := range b.Connections()[1:] {
		if c.Distance() < closest.Distance() {
			closest = c
		}
	}

	scores := make([]float64, len(b.Connections()))
	// Score the nodes
	for i, c := range b.Connections() {
		// Closest connection has a perfect score
		if c == closest {
			scores[i] = 0
		}

		// distance from building x number of connections on building
		scores[i] = c.Distance() * float64(len(b.Connections()))
	}

	return scores
}

func averageNumberOfConnections(buildings []*world.Building) float64 {
	number := 0
	for _, b := range buildings {
		number += len(b.Connections())
	}

	return float64(number) / float64(len(buildings))
}
