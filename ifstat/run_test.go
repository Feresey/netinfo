package ifstat

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func writeFiles(rxFile, txFile *os.File, rx, tx []byte) {
	_, err := rxFile.WriteAt(rx, 0)
	if err != nil {
		panic(err)
	}
	_, err = txFile.WriteAt(tx, 0)
	if err != nil {
		panic(err)
	}
}

func TestIfStat(t *testing.T) {
	dir := os.TempDir() + "/" + uuid.New().String() + "/"
	err := os.Mkdir(dir, 0777)
	if err != nil {
		if !os.IsExist(err) {
			t.Error(err)
			return
		}
	}

	rx, err := os.Create(dir + "rx_bytes")
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := os.Create(dir + "tx_bytes")
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		rx   []byte
		tx   []byte
		want pair
	}{
		{
			name: "normal",
			rx:   []byte("123\n"),
			tx:   []byte("456\n"),
			want: pair{123, 456},
		},
		{
			name: "rx corrupted",
			rx:   []byte("aaa\n"),
			tx:   []byte("456\n"),
			want: pair{0, 456},
		},
		{
			name: "tx corrupted",
			rx:   []byte("123\n"),
			tx:   []byte("bbb\n"),
			want: pair{123, 0},
		},
	}

	stat := &IfStat{Path: []string{dir}, Delay: time.Millisecond, Out: os.Stdout}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeFiles(rx, tx, tt.rx, tt.tx)
			p := stat.MustRead()
			if !reflect.DeepEqual(tt.want, p) {
				t.Errorf("result does not match.\nExpected:\n%#v\nGiven:\n%#v", tt.want, p)
			}
		})
	}

	t.Run("readDetached", func(t *testing.T) {
		writeFiles(rx, tx, []byte("123\n"), []byte("123\n"))
		c, cancel := stat.readDetached()
		go func() {
			want := pair{123, 123}
			for p := range c {
				if !reflect.DeepEqual(p, want) {
					t.Errorf("result does not match.\nExpected:\n%#v\nGiven:\n%#v", want, p)
				}
			}
		}()
		time.Sleep(20 * time.Millisecond)
		cancel()
		_, ok := <-c
		if ok {
			t.Error("channel must be closed")
		}
	})

	t.Run("Run", func(t *testing.T) {
		cancel := stat.Run()
		time.Sleep(2 * time.Millisecond)
		cancel()
	})

	_ = rx.Close()
	_ = tx.Close()
	_ = os.RemoveAll(dir)
	t.Run("no files", func(t *testing.T) {
		p := stat.MustRead()
		if !reflect.DeepEqual(pair{}, p) {
			t.Errorf("result does not match.\nExpected:\n%#v\nGiven:\n%#v", pair{}, p)
		}
	})
}

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
