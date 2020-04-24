package overpass

import (
	"encoding/json"
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
	invalidCoords(0, 0, 0, math.NaN())
	invalidCoords(0, 0, 0, math.Inf(-1))
	invalidCoords(0, 0, 0, math.Inf(1))

	testData := &Result{
		Version: "test",
		Elements: []*Way{
			{
				Id:     1,
				Bounds: [4]float64{2, 3, 4, 5},
				Nodes:  []uint64{6, 7},
				Geometry: []*LatLon{
					{8.0, 9.0},
				},
				Tags: &Tags{
					Highway:  "primary",
					Building: "yes",
				},
			},
		},
	}

	// Create a test server that will work
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(testData)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(body)
		require.NoError(t, err)
	}))
	overpassEndpoint = testServer.URL

	resp, err := ExecuteQuery(0, 0, 0, 0)
	require.NoError(t, err)
	require.Equal(t, testData, resp)

	testServer.Close()

	// Create a test server that returns the 500 status code
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte("500 - test error"))
		require.NoError(t, err)
	}))
	overpassEndpoint = testServer.URL

	resp, err = ExecuteQuery(0, 0, 0, 0)
	require.Error(t, err)

	testServer.Close()

	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err = w.Write([]byte("This is not valid JSON"))
		require.NoError(t, err)
	}))
	overpassEndpoint = testServer.URL

	resp, err = ExecuteQuery(0, 0, 0, 0)
	require.Error(t, err)

	testServer.Close()
}
