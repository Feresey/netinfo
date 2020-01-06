package ifstat

import (
	"io"
	"os"
	"time"
)

type InterfaceNotExists string

type offsetReader interface {
	Read([]byte) (int, error)
	Close() error
}

type fileWithOffset os.File

type readPair struct {
	rx offsetReader
	tx offsetReader
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
