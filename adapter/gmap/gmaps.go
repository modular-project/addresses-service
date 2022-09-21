package gmaps

import (
	"context"
	"fmt"

	"github.com/modular-project/address-service/model"
	"googlemaps.github.io/maps"
)

type gMapService struct {
	c *maps.Client
}

func NewGMapService(key string) (gMapService, error) {
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		return gMapService{}, fmt.Errorf("NewClient: %w", err)
	}
	return gMapService{c}, nil
}

func (gm gMapService) GeoCode(ctx context.Context, add string) (model.Location, error) {
	gr := maps.GeocodingRequest{
		Address: add,
	}
	res, err := gm.c.Geocode(ctx, &gr)
	if err != nil {
		return model.Location{}, fmt.Errorf("geocode: %w", err)
	}
	if res == nil {
		return model.Location{}, fmt.Errorf("geocode: nil response")
	}
	r := res[0]
	loc := model.Location{
		Type:        "Point",
		Coordinates: []float64{r.Geometry.Location.Lng, r.Geometry.Location.Lat},
	}
	return loc, nil
}
