package convert

import (
	"github.com/real-life-td/world-generator/overpass"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClassify(t *testing.T) {
	test := func (tags overpass.Tags, expected elementType, shouldError bool) {
		w := overpass.Way{
			Id:       0,
			Bounds:   [4]int{0, 0, 0, 0},
			Nodes:    nil,
			Geometry: nil,
			Tags:     &tags,
		}

		r, err := classify(&w)
		if shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, expected, r)
		}
	}

	test(overpass.Tags{Building: "", Highway: ""}, -1, true)
	test(overpass.Tags{Building: "yes", Highway: ""}, BUILDING, false)
	test(overpass.Tags{Building: "", Highway: "primary"}, HIGHWAY, false)
	test(overpass.Tags{Building: "yes", Highway: "primary"}, -1, true)
}
