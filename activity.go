package activity

import (
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-gpx"
)

const UserAgent = "github.com/bzimmer/activity"

type GPXer interface {
	// GPX returns a GPX instance
	GPX() (*gpx.GPX, error)
}

type GeoJSONer interface {
	// GeoJSON returns a GeoJSON instance
	GeoJSON() (*geojson.Feature, error)
}
