package zwift_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/activity/zwift"
)

func TestExporter(t *testing.T) {
	a := assert.New(t)
	client, err := zwift.NewClient()
	a.NoError(err)
	a.NotNil(client)
	a.NotNil(client.Exporter())
}
