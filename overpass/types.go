package overpass

type overpassResult struct {
	version string
	elements []*way
}

type latLon struct {
	lat float64
	lon float64
}

type tags struct {
	highway string
	building string
}

type way struct {
	id int
	bounds [4]int
	nodes []int
	geometry []*latLon
	tags *tags
}