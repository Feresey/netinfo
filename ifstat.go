package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	kb = 1 << (10 * (iota + 1))
	mb
)

type readPair struct {
	rx string
	tx string
}

type (
	netSpeed int

	pair struct {
		rx netSpeed
		tx netSpeed
	}
)

const printFormat = "%v %v\n"

type ifStat struct {
	path  []readPair
	delay time.Duration
	out   io.Writer
}

func (speed netSpeed) String() (res string) {
	switch {
	case speed < kb:
		res += fmt.Sprintf("%7.1fB ", float32(speed))
	case speed >= kb && speed < mb:
		res += fmt.Sprintf("%7.1fKb", float32(speed)/float32(kb))
	case speed >= mb:
		res += fmt.Sprintf("%7.1fMb", float32(speed)/float32(mb))
	}

	return
}

func newStat(accumulate time.Duration, names ...string) *ifStat {
	path := make([]readPair, 0, len(names))

	for _, ifname := range names {
		if err := checkIfaceExists(ifname); err != nil {
			log.Println(err)
			continue
		}

		prefix := filepath.Join("/sys/class/net", ifname, "statistics")

		path = append(path,
			readPair{
				rx: filepath.Join(prefix, "rx_bytes"),
				tx: filepath.Join(prefix, "tx_bytes"),
			})
	}

	return &ifStat{
		path:  path,
		delay: accumulate,
		out:   os.Stdout,
	}
}

func checkIfaceExists(name string) error {
	check, err := regexp.Compile(`^\s*` + name + ": .*$")
	if err != nil {
		return err
	}

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if check.MatchString(scanner.Text()) {
			return nil
		}
	}

	return fmt.Errorf("interface not exists: %s", name)
}

func (i *ifStat) Run() func() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		last        = i.readDetached(ctx)
	)

	go i.printDetached(last)

	return cancel
}

func (i *ifStat) readDetached(ctx context.Context) (status <-chan pair) {
	var (
		last   = make(chan pair)
		ticker = time.NewTicker(i.delay)
	)

	go func() {
		last <- i.read()

		for {
			select {
			case <-ticker.C:
				last <- i.read()
			case <-ctx.Done():
				close(last)
				return
			}
		}
	}()

	return last
}

func (i *ifStat) printDetached(last <-chan pair) {
	var (
		prod = netSpeed(time.Second / i.delay)
		prev = <-last
	)

	for curr := range last {
		var (
			rxInt = (curr.rx - prev.rx) * prod
			txInt = (curr.tx - prev.tx) * prod
		)

		_, err := fmt.Fprintf(i.out, printFormat, rxInt, txInt)
		if err != nil {
			log.Print(err)
		}

		prev = curr
	}
}

func mustGetInt(r string) (res netSpeed) {
	file, err := os.Open(r)
	if err != nil {
		return
	}

	_, err = fmt.Fscan(file, &res)
	if err != nil {
		log.Print(err)
		return
	}

	err = file.Close()
	if err != nil {
		log.Print(err)
	}

	return
}

func (i *ifStat) read() (res pair) {
	for _, ifpath := range i.path {
		res.rx += mustGetInt(ifpath.rx)
		res.tx += mustGetInt(ifpath.tx)
	}

	return
}
