package main

import (
	"fmt"
	"os"
	"time"
)

const (
	KB int = 1 << (10 * (iota + 1))
	MB
)

func getInt(r string) (res int) {
	file, err := os.Open(r)
	if err != nil {
		return 0
	}
	fmt.Fscan(file, &res)
	file.Close()
	return
}

func (i *IfStat) MustRead() (res pair) {
	for _, ifpath := range i.Path {
		res.rx += getInt(ifpath.rx)
		res.tx += getInt(ifpath.tx)
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
				// for _, p := range I.Path {
				// 	p.rx.Close()
				// }
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

		_, _ = I.Out.Write([]byte(prettyPrint(rxInt) + " " + prettyPrint(txInt) + "\n"))

		prev = curr
	}
}

func (I *IfStat) Run() func() {
	last, done := I.readDetached()
	go I.runDetached(last)
	return done
}
