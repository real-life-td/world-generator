package overpass

type Result struct {
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

type Bounds struct {
	MinLat, MinLon, MaxLat, MaxLon float64
}

type Way struct {
	Id       uint64
	Bounds   *Bounds
	Nodes    []uint64
	Geometry []*LatLon
	Tags     *Tags
}
