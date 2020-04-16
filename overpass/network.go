package overpass

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
)

var overpassEndpoint = "https://overpass-api.de/api/interpreter"
const query = "[bbox:%f,%f,%f,%f];(way[highway];way[building];);out geom;"

func call(query string) (body io.ReadCloser, err error) {
	req, err := http.NewRequest("GET", overpassEndpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("data", query)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return resp.Body, nil
}

func makeQuery(lat1, lon1, lat2, lon2 float64) string {
	if lat1 > lat2 {
		lat1, lat2 = lat2, lat1
	}

	if lon1 > lon2 {
		lon1, lon2 = lon2, lon1
	}

	return fmt.Sprintf(query, lat1, lon1, lat2, lon2)
}

func checkCoordinates(lat, lon float64) error {
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return errors.New("latitude must be non-infinite and not nan")
	}

	if -90.0 > lat || lat > 90.0 {
		return errors.New("latitude must be in range -90 to 90 inclusive")
	}

	if math.IsNaN(lon) || math.IsInf(lon, 0) {
		return errors.New("longitude must be non-infinite and not nan")
	}

	if -180.0 > lon || lon > 180.0 {
		return errors.New("longitude must be in range -180 to 180 inclusive")
	}

	return nil
}

func ExecuteQuery(lat1, lon1, lat2, lon2 float64) (result *overpassResult, err error) {
	err = checkCoordinates(lat1, lon1)
	if err != nil {
		return nil, err
	}

	err = checkCoordinates(lat2, lon2)
	if err != nil {
		return nil, err
	}

	query := makeQuery(lat1, lon1, lat2, lon2)
	body, err := call(query)
	if err != nil {
		return nil, err
	}

	result = new(overpassResult)
	err = json.NewDecoder(body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}