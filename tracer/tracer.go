package tracer

import (
	"fmt"
	"io"
	time2 "time"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(i ...interface{}) {
	time := time2.Now().Format("2006-01-02 15:04:05.999999")
	fmt.Fprintf(t.out, time+" TRACE: %s\n", fmt.Sprint(i...))
}

func New(out io.Writer) Tracer {
	return &tracer{out: out}
}
