package overpass

type Result struct {
	Version  string
	Elements []*Way
}

type LatLon struct {
	Lat float64
	Lon float64
}

type Tags struct {
	Highway  string
	Building string
}

type Way struct {
	Id       uint64
	Bounds   [4]int
	Nodes    []uint64
	Geometry []*LatLon
	Tags     *Tags
}
