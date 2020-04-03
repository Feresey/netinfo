package main

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

func Test_prettyPrint(t *testing.T) {
	tests := []struct {
		name    string
		args    int
		wantRes string
	}{
		{
			args:    1,
			wantRes: "    1.0B ",
		},
		{
			args:    1024,
			wantRes: "    1.0Kb",
		},
		{
			args:    2 * 1025,
			wantRes: "    2.0Kb",
		},
		{
			args:    5 * 1024 * 1024,
			wantRes: "    5.0Mb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := prettyPrint(tt.args); gotRes != tt.wantRes {
				t.Errorf("prettyPrint() = %q, want %q", gotRes, tt.wantRes)
			}
		})
	}
}

func TestIfStat_runDetached(t *testing.T) {
	out := bytes.NewBuffer(nil)
	stat := &IfStat{Delay: time.Second, Out: out}
	data := make(chan pair)
	go stat.runDetached(data)

	item := pair{1, 2}
	for i := 0; i < 10; i++ {
		data <- item
		item.rx++
		item.tx++
	}
	close(data)
	line := prettyPrint(1) + " " + prettyPrint(1)
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		if scanner.Text() != line {
			t.Errorf("result does not match.\nExpected:\n%q\nGiven:\n%q", line, scanner.Text())
		}
	}
}
