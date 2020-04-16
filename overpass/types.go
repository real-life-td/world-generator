package overpass

type overpassResult struct {
	Version  string
	Elements []*way
}

type latLon struct {
	Lat float64
	Lon float64
}

type tags struct {
	Highway  string
	Building string
}

type way struct {
	Id       int
	Bounds   [4]int
	Nodes    []int
	Geometry []*latLon
	Tags     *tags
}
