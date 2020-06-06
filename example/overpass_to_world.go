package main

import (
	"errors"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/convert"
	"github.com/real-life-td/world-generator/mutate"
	"github.com/real-life-td/world-generator/overpass"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type queryParams struct {
	lat1, lon1, lat2, lon2 float64
}

type traversalColor byte
const (
	white traversalColor = iota
	grey
	black
)

func main() {
	log.Println("Starting server at localhost:8080")
	http.Handle("/", http.HandlerFunc(example))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error creating server", err)
	}
}

func getLatLonParam(url *url.URL, name string) (result float64, err error) {
	value := url.Query().Get(name)

	if value == "" {
		return math.NaN(), errors.New("missing required parameter: " + name + "=<number>")
	}

	result, err = strconv.ParseFloat(value, 65)
	return
}

func getQueryParams(url *url.URL) (params *queryParams, err error) {
	lat1, err := getLatLonParam(url, "lat1")
	if err != nil {
		return nil, err
	}

	lon1, err := getLatLonParam(url, "lon1")
	if err != nil {
		return nil, err
	}

	lat2, err := getLatLonParam(url, "lat2")
	if err != nil {
		return nil, err
	}

	lon2, err := getLatLonParam(url, "lon2")
	if err != nil {
		return nil, err
	}

	return &queryParams{lat1: lat1, lon1: lon1, lat2: lat2, lon2: lon2}, nil
}

func example(w http.ResponseWriter, req *http.Request) {
	// Get latitude and longitude parameters
	params, err := getQueryParams(req.URL)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintln(w, "Invalid query parameters: "+err.Error())
		return
	}

	println(time.Now().UnixNano())
	// Execute an overpass query
	result, err := overpass.ExecuteQuery(params.lat1, params.lon1, params.lat2, params.lon2)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintln(w, "Internal error when executing Overpass query: "+err.Error())
		return
	}
	println(time.Now().UnixNano())

	// Convert the result into a world object
	metadata := world.NewMetadata(1000, 1000, params.lat1, params.lon1, params.lat2, params.lon2)
	world, err := convert.Convert(metadata, result)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintln(w, "Internal error when executing converting: "+err.Error())
		return
	}
	println(time.Now().UnixNano())

	mutate.InitBuildingConnections(world, 2.0)

	println(time.Now().UnixNano())

	// Translate the world into SVG
	w.Header().Set("Content-Type", "image/svg+xml")

	s := svg.New(w)
	s.Start(metadata.Width(), metadata.Height())
	renderContainer(s, world)
	s.End()
}

func renderContainer(s *svg.SVG, container *world.Container) {
	renderRoads(s, container.Roads())

	for _, b := range container.Buildings() {
		renderBuilding(s, b)
	}
}

func renderRoads(s *svg.SVG, roads []*world.Road) {
	visited := make(map[world.Id]traversalColor)

	// Uses a breadth first search to render the road network. This implementation ensures that each road segment is
	// rendered only once.
	traversalRender := func (start *world.Road) {
		toVisit := []*world.Road{start}
		for len(toVisit) > 0 {
			cur := toVisit[0]
			toVisit = toVisit[1:]
			visited[cur.Id()] = black

			for _, connected := range cur.Connections() {
				color := visited[connected.Id()]

				if color == white || color == grey {
					s.Line(cur.X(), cur.Y(), connected.X(), connected.Y(), "stroke-width:3;stroke:rgb(0,0,255);stroke-linecap:round;")
				}

				if color == white {
					visited[connected.Id()] = grey
					toVisit = append(toVisit, connected)
				}
			}
		}
	}

	// Ensure that all roads are rendered if they are not connected to others
	for _, r := range roads {
		if visited[r.Id()] == white {
			traversalRender(r)
		}
	}
}

func renderBuilding(s *svg.SVG, building *world.Building) {
	x := make([]int, 0, len(building.Points()))
	y := make([]int, 0, len(building.Points()))

	for _, b := range building.Points() {
		x = append(x, b.X())
		y = append(y, b.Y())
	}

	s.Polygon(x, y)

	for _, c := range building.Connections() {
		s.Circle(c.Road().X(), c.Road().Y(), 4, "stroke-width:1;stroke:rgb(0,255,0);fill:none")
		s.Line(c.Road().X(), c.Road().Y(), c.PointOnBuilding().X(), c.PointOnBuilding().Y(), "stroke-width:1;stroke:rgb(0,255,0);")
	}
}
