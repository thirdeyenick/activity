package activity

import (
	"time"

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

// A Named provides a minimal set of metadata about an entity
type Named struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Source string    `json:"source"`
	Date   time.Time `json:"date"`
}

// A Namer returns a Named instance
type Namer interface {
	Named() *Named
}
