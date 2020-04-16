package overpass

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteQuery(t *testing.T) {
	// Passing invalid latitude or longitude should result in an error
	invalidCoords := func(lat1, lon1, lat2, lon2 float64) {
		_, err := ExecuteQuery(lat1, lon1, lat2, lon2)
		require.Error(t, err)
	}

	invalidCoords(-91, 0, 0, 0)
	invalidCoords(91, 0, 0, 0)
	invalidCoords(math.NaN(), 0, 0, 0)
	invalidCoords(math.Inf(-1), 0, 0, 0)
	invalidCoords(math.Inf(1), 0, 0, 0)
	invalidCoords(0, -181, 0, 0)
	invalidCoords(0, 181, 0, 0)
	invalidCoords(0, math.NaN(), 0, 0)
	invalidCoords(0, math.Inf(-1), 0, 0)
	invalidCoords(0, math.Inf(1), 0, 0)
	invalidCoords(0, 0, -91, 0)
	invalidCoords(0, 0, 91, 0)
	invalidCoords(0, 0, math.NaN(), 0)
	invalidCoords(0, 0, math.Inf(-1), 0)
	invalidCoords(0, 0, math.Inf(1), 0)
	invalidCoords(0, 0, 0, -181)
	invalidCoords(0, 0, 0, 181)
	invalidCoords(0, 0,0, math.NaN())
	invalidCoords(0, 0,0, math.Inf(-1))
	invalidCoords(0, 0,0, math.Inf(1))

	testData := &overpassResult{
		Version: "test",
		Elements: []*way{
			{
				Id: 1,
				Bounds: [4]int{2, 3, 4, 5},
				Nodes: []int{6, 7},
				Geometry: []*latLon{
					{8.0, 9.0},
				},
				Tags: &tags{
					Highway: "primary",
					Building: "yes",
				},
			},
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(testData)
		require.NoError(t, err)

		_, err = fmt.Fprint(w, string(body))
		require.NoError(t, err)
	}))
	overpassEndpoint = testServer.URL
	defer testServer.Close()

	resp, err := ExecuteQuery(0, 0, 0, 0)
	require.NoError(t, err)
	require.Equal(t, testData, resp)
}
