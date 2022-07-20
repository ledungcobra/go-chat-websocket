package tracer

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buff bytes.Buffer
	tracer := New(&buff)

	if tracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		tracer.Trace("Hello World")
	}
}
