package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func makeUniq(slice []string) []string {
	if len(slice) <= 1 {
		return slice
	}
	write := 1

	for read := 1; read < len(slice); read++ {
		if slice[read] != slice[read-1] {
			slice[write] = slice[read]
			write++
		}
	}
	return slice[:write]
}

func NewStat(names ...string) *IfStat {
	sort.Strings(names)
	names = makeUniq(names)
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

	return &IfStat{Path: path, Delay: time.Second, Out: os.Stdout}
}

func checkIfaceExists(name string) error {
	if name == "" {
		return errors.New(name)
	}
	file, _ := os.Open("/proc/net/dev")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), name+":") {
			return nil
		}
	}

	return errors.New(name)
}
