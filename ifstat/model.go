package ifstat

import (
	"io"
	"os"
	"time"
)

type InterfaceNotExists string

type fileWithOffset os.File

type readPair struct {
	rx io.ReadCloser
	tx io.ReadCloser
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
