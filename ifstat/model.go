package ifstat

import (
	"io"
	"time"
)

type InterfaceNotExists string

type IfStat struct {
	Path  []string
	Delay time.Duration
	Out   io.Writer
}

type pair struct {
	rx int
	tx int
}
