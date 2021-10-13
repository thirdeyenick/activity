package strava_test

import (
	"testing"

	"github.com/bzimmer/activity"
	"github.com/bzimmer/activity/strava"
	"github.com/stretchr/testify/assert"
)

func TestUploadDone(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	u := &strava.Upload{ID: 10}
	a.Equal(activity.UploadID(10), u.Identifier())
	a.False(u.Done())
	a.True((&strava.Upload{Error: "error"}).Done())
	a.True((&strava.Upload{ActivityID: 1234567890}).Done())
}
