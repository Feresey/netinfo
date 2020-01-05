package ifstat

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

const (
	KB int = 1 << (10 * (iota + 1))
	MB
)

func getInt(path string) int {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	if len(raw) < 2 || raw[len(raw)-1] != '\n' {
		log.Printf("File %q corrupted", path)
		return 0
	}
	res, err := strconv.Atoi(string(raw[:len(raw)-1]))
	if err != nil {
		log.Println(err)
	}

	return res
}

func (i *IfStat) MustRead() (res pair) {
	for _, ifpath := range i.Path {
		res.rx += getInt(ifpath + "rx_bytes")
		res.tx += getInt(ifpath + "tx_bytes")
	}
	return
}

func prettyPrint(speed int) (res string) {
	switch {
	case speed < KB:
		res += fmt.Sprintf("%7.1fB ", float32(speed))
	case speed >= KB && speed < MB:
		res += fmt.Sprintf("%7.1fKb", float32(speed)/float32(KB))
	case speed >= MB:
		res += fmt.Sprintf("%7.1fMb", float32(speed)/float32(MB))
	}
	return
}

func (I *IfStat) readDetached() (<-chan pair, func()) {
	last := make(chan pair)
	done := make(chan struct{})
	ticker := time.NewTicker(I.Delay)
	go func() {
		last <- I.MustRead()
		last <- I.MustRead()
		for {
			select {
			case <-ticker.C:
				last <- I.MustRead()
			case <-done:
				close(last)
				return
			}
		}
	}()
	return last, func() { done <- struct{}{} }
}

func (I *IfStat) runDetached(last <-chan pair) {
	prod := int(time.Second / I.Delay)
	prev := <-last
	for curr := range last {
		rxInt := (curr.rx - prev.rx) * prod
		txInt := (curr.tx - prev.tx) * prod

		fmt.Fprintln(I.Out, prettyPrint(rxInt)+" "+prettyPrint(txInt))

		prev = curr
	}
}

func (I *IfStat) Run() func() {
	last, done := I.readDetached()
	go I.runDetached(last)
	return done
}
