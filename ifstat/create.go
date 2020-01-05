package ifstat

import (
	"bufio"
	"log"
	"os"
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
		s := "/sys/class/net/" + ifname + "/statistics/"
		rx, _ := os.Open(s + "rx_bytes")
		tx, _ := os.Open(s + "tx_bytes")
		path = append(path,
			readPair{
				(*fileWithOffset)(rx),
				(*fileWithOffset)(tx),
			})
	}

	return &IfStat{Path: path, Delay: time.Second, Out: os.Stdout}
}

func checkIfaceExists(name string) error {
	if name == "" {
		return InterfaceNotExists(name)
	}
	file, _ := os.Open("/proc/net/dev")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), name+":") {
			return nil
		}
	}

	return InterfaceNotExists(name)
}
