package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	var accumulate int

	flag.IntVar(&accumulate, "t", 1000, "milliseconds to accumulate data to calculate speed")

	flag.Usage = func() {
		fmt.Printf(`Usage of %s
  t=1000: milliseconds to accumulate data to calculate speed
  all other arguments reads as interface names, such as eth0, enp2s3 and others`, os.Args[0])
	}

	flag.Parse()

	names := flag.Args()
	if len(names) == 0 {
		names = []string{"eno1", "wlan0"}
	}

	var (
		iface   = newStat(time.Duration(accumulate)*time.Millisecond, names...)
		cancel  = iface.Run()
		sigChan = make(chan os.Signal, 1)
		sig     os.Signal
	)

	signal.Notify(sigChan, os.Interrupt)

	sig = <-sigChan

	fmt.Println("\nCaught signal:", sig)
	cancel()
}
