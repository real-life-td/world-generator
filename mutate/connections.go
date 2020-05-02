package mutate

import (
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/math/primitives"
	"github.com/real-life-td/math/raycast"
	"math"
)

type closestPoint struct {
	road *world.Road // the road that this closest point was computed for
	distance float64
}

func InitBuildingConnections(container *world.Container, targetAvgConnections float64) {
	initBuildingConnections(container, 4, targetAvgConnections)
}

func initBuildingConnections(container *world.Container, initialConnectDistance int, targetAvgConnections float64) {
	for _, b := range container.Buildings() {
		connectionBounds := b.Bounds().Expand(initialConnectDistance, initialConnectDistance, initialConnectDistance, initialConnectDistance)
		withinBounds := withinBounds(connectionBounds, container.Roads())

		closeEnough := make(map[world.Id]*closestPoint) // Map to make the next step more efficient
		for _, r := range withinBounds {
			closest := closestDistance(b, r)
			if int(math.Round(closest.distance)) < initialConnectDistance {
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
		initBuildingConnections(container, initialConnectDistance * 2, targetAvgConnections)
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

func closestDistance(building *world.Building, r *world.Road) *closestPoint {
	closestDistance := math.MaxFloat64
	for i, bPoint := range building.Points() {
		var bNextPoint *world.Node
		if i == len(building.Points()) - 1 {
			// In the last iteration the next point is at the start of the slice
			bNextPoint = building.Points()[0]
		} else {
			bNextPoint = building.Points()[i + 1]
		}

		_, distance := raycast.ClosestPointTo(bPoint.Point, bNextPoint.Point, r.Point)
		closestDistance = math.Min(closestDistance, distance)
	}

	return &closestPoint{
		road: r,
		distance: closestDistance,
	}
}

// Returns only roads that are closer then any other road that is connected and connected by only other roads that are
// connected (within 5 roads). Because of the 5 road limit and the faster implementation used this function can return
// different results based on the order of the roads
func cullPaths(closeEnough map[world.Id]*closestPoint) []*world.Road {
	passing := make([]*world.Road, 0)

	visited := make(map[world.Id]bool)
	for _, c := range closeEnough {
		if visited[c.road.Id()] {
			continue
		}

		closest := c
		toVisit := append([]*world.Road(nil), c.road.Connections()...)
		curLevel := 0 // How many roads "jumps" from the original road c
		nextIncrement := 1 // How many roads to process in toVisit before the next level

		for len(toVisit) > 0 {
			nextIncrement--
			if nextIncrement == 0 {
				curLevel++
				if curLevel == 6 {
					break
				}

				nextIncrement = len(toVisit)
			}

			cur := toVisit[0]
			toVisit = toVisit[1:]

			curData, present := closeEnough[cur.Id()]
			if present {
				if curData.distance < closest.distance {
					closest = curData
				}

				for _, connected := range cur.Connections() {
					if !visited[connected.Id()] {
						visited[connected.Id()] = true
						toVisit = append(toVisit, connected)
					}
				}
			}
		}

		passing = append(passing, closest.road)
	}

	return passing
}

func averageNumberOfConnections(buildings []*world.Building) float64 {
	number := 0
	for _, b := range buildings {
		number += len(b.Connections())
	}

	return float64(number) / float64(len(buildings))
}