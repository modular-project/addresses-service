package gmaps

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/modular-project/address-service/model"
)

func Test_gMapService_GeoCode(t *testing.T) {
	type args struct {
		ctx context.Context
		add string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Location
		wantErr bool
	}{
		{
			name: "OK API KEY",
			args: args{context.Background(), "Blvd. Gral. Marcelino García Barragán 1421, Olímpica, 44430 Guadalajara, Jal."},
			want: model.Location{Type: "Point", Coordinates: []float64{-103.3266212, 20.6545464}},
		},
	}
	key, ok := os.LookupEnv("GMAP_APIKEY")
	if !ok {
		t.Fatal("enviroment variable GMAP_APIKEY not found")
	}
	gms, err := NewGMapService(key)
	if err != nil {
		t.Fatalf("NewGMapService: %s", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gms.GeoCode(tt.args.ctx, tt.args.add)
			if (err != nil) != tt.wantErr {
				t.Errorf("gMapService.GeoCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gMapService.GeoCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
