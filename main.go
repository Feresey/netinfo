package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Feresey/netinfo/ifstat"
)

func main() {
	accumulate := flag.Int("t", 1000, "milliseconds to accumulate data to calculate speed")
	flag.Parse()

	names := flag.Args()
	if len(names) == 0 {
		names = []string{"eno1", "wlan0"}
	}

	iface, err := ifstat.NewStat(names...)
	if err != nil {
		panic(err)
	}

	iface.Delay = time.Duration(*accumulate) * time.Millisecond

	cancel := iface.Run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	fmt.Println("\nCaught signal:", sig)
	cancel()
}
