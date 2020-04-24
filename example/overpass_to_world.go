package main

import (
	"errors"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/convert"
	"github.com/real-life-td/world-generator/overpass"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type queryParams struct {
	lat1, lon1, lat2, lon2 float64
}

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
		_, _ = fmt.Fprintln(w, "Invalid query parameters: " + err.Error())
		return
	}

	// Execute an overpass query
	result, err := overpass.ExecuteQuery(params.lat1, params.lon1, params.lat2, params.lon2)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintln(w, "Internal error when executing Overpass query: " + err.Error())
		return
	}

	// Convert the result into a world object
	metadata := world.NewMetadata(1000, 1000, params.lat1, params.lon1, params.lat2, params.lon2)
	world, err := convert.Convert(metadata, result)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprintln(w, "Internal error when executing converting: " + err.Error())
		return
	}

	// Translate the world into SVG
	w.Header().Set("Content-Type", "image/svg+xml")

	s := svg.New(w)
	s.Start(metadata.Width(),metadata.Height())
	renderContainer(s, world)
	s.End()
}

func renderContainer(s *svg.SVG, container *world.Container) {
	for _, r := range container.Roads() {
		renderRoad(s, r)
	}

	for _, b := range container.Buildings() {
		renderBuilding(s, b)
	}
}

func renderRoad(s *svg.SVG, road *world.Road) {
	s.Line(road.Node1().X(), road.Node1().Y(), road.Node2().X(), road.Node2().Y(), "stroke-width:3;stroke:rgb(0,0,0)")
}

func renderBuilding(s *svg.SVG, building *world.Building) {
	x := make([]int, 0, len(building.Points()))
	y := make([]int, 0, len(building.Points()))

	for _, b := range building.Points() {
		x = append(x, b.X())
		y = append(y, b.Y())
	}

	s.Polygon(x, y)
}