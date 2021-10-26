package inreach_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/activity/inreach"
)

func TestFault(t *testing.T) {
	a := assert.New(t)
	f := &inreach.Fault{Message: "foobar"}
	a.Error(f)
	a.Equal("foobar", f.Error())
}

func TestEvents(t *testing.T) {
	a := assert.New(t)
	events := inreach.Events()
	a.NotNil(events)
	a.Len(events, 18)
	a.Contains(events, 13)
	a.Equal(events[1], "Emergency initiated from device")
}
