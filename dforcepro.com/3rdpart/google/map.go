package google

import (
	"context"
	"fmt"
	"strings"

	"googlemaps.github.io/maps"
)

type Map interface {
	Geocode()
	Matrix()
}

func Geocode(Address string) maps.LatLng {

	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyB1Th0z68puegij0hHqcMAm_bYnpMoiiTM"))
	if err != nil {
		fmt.Println(err)
	}
	r := &maps.GeocodingRequest{
		Address:  Address,
		Language: "zh-TW",
	}
	route, err := c.Geocode(context.Background(), r)
	if err != nil {
		fmt.Println(err)

	}
	return route[0].Geometry.Location

}
func Distance(or string, des string) int {

	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyB1Th0z68puegij0hHqcMAm_bYnpMoiiTM"))
	if err != nil {
		fmt.Println(err)
	}

	r := &maps.DistanceMatrixRequest{
		Language: "zh-TW",
	}

	r.Origins = strings.Split(or, "|")

	r.Destinations = strings.Split(des, "|")

	resp, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		fmt.Println(err)

	}
	return resp.Rows[0].Elements[0].Distance.Meters
}
