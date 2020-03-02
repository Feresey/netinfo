package ifstat

import (
	"io"
	"time"
)

type readPair struct {
	rx string
	tx string
}

type IfStat struct {
	Path  []readPair
	Delay time.Duration
	Out   io.Writer
}

type pair struct {
	rx int
	tx int
}
