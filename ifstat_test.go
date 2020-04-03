package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_checkIfaceExists(t *testing.T) {
	tests := []struct {
		name    string
		ifname  string
		wantErr bool
	}{
		{
			name:   "lo",
			ifname: "lo",
		},
		{
			name:    "really not interface",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := checkIfaceExists(tt.ifname)

			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("new stat", func(t *testing.T) {
		newStat(time.Second, "eno1", "eno1", "", "error")
	})
}

func TestIfStat_runDetached(t *testing.T) {
	var (
		out       = new(bytes.Buffer)
		stat      = &ifStat{delay: time.Second, out: out}
		data      = make(chan pair)
		item      = pair{1, 1}
		line      = fmt.Sprintf(printFormat, item.rx, item.tx)
		lineCount = 10
		wantLines string
	)

	for i := 0; i < lineCount-2; i++ {
		wantLines += line
	}

	go stat.printDetached(data)
	defer close(data)

	for i := 0; i < lineCount; i++ {
		data <- item
		item.rx++
		item.tx++
	}

	assert.Equal(t, wantLines, out.String())
}

func TestByteSpeedString(t *testing.T) {
	tests := []struct {
		args    netSpeed
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
		assert.Equal(t, tt.wantRes, tt.args.String())
	}
}
