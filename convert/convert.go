package convert

import (
	"errors"
	"github.com/real-life-td/game-core/world"
	"github.com/real-life-td/world-generator/overpass"
)

type elementType int
const (
	BUILDING elementType = iota
	HIGHWAY
)

func Convert(result *overpass.Result) (world *world.Container, err error) {
	return nil, nil
}

func classify(e *overpass.Way) (t elementType, err error) {
	if e.Tags.Building != "" {
		if e.Tags.Highway != "" {
			return -1, errors.New("element is both building and highway type")
		}

		return BUILDING, nil
	} else if e.Tags.Highway != "" {
		return HIGHWAY, nil
	}

	return -1, errors.New("element does not have type")
}
