package cyclinganalytics_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/activity/cyclinganalytics"
)

func TestDateTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		date cyclinganalytics.Datetime
	}{
		{
			name: "success",
			date: cyclinganalytics.Datetime{Time: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			t.Parallel()
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			a.NoError(enc.Encode(tt.date))
			dec := json.NewDecoder(&buf)
			var dst cyclinganalytics.Datetime
			a.NoError(dec.Decode(&dst))
			a.Equal(tt.date, dst)
		})
	}
}
