package ifstat

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

type fakeBuffer struct {
	data []byte
}

func (B *fakeBuffer) Read(b []byte) (int, error) {
	if B.data == nil {
		return 0, io.EOF
	}
	copy(b, B.data)
	return len(B.data), nil
}

func (B *fakeBuffer) Close() error {
	B.data = nil
	return nil
}

func (B *fakeBuffer) IsClosed() bool { return B.data == nil }

func TestIfStat(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat := &IfStat{
				Path: []readPair{
					{&fakeBuffer{tt.rx}, &fakeBuffer{tt.tx}},
				},
			}
			p := stat.MustRead()
			if !reflect.DeepEqual(tt.want, p) {
				t.Errorf("result does not match.\nExpected:\n%#v\nGiven:\n%#v", tt.want, p)
			}
		})
	}

	t.Run("readDetached", func(t *testing.T) {
		buf := &fakeBuffer{[]byte("123\n")}
		stat := &IfStat{
			Path:  []readPair{{buf, buf}},
			Delay: time.Millisecond,
			Out:   os.Stdout,
		}
		wait := make(chan struct{})
		c, stop := stat.readDetached()
		go func() {
			want := pair{123, 123}
			for p := range c {
				if !reflect.DeepEqual(p, want) {
					t.Errorf("result does not match.\nExpected:\n%#v\nGiven:\n%#v", want, p)
					break
				}
			}
			wait <- struct{}{}
		}()
		time.Sleep(20 * time.Millisecond)
		stop()
		<-wait
		_, ok := <-c
		if ok {
			t.Error("channel must be closed")
		}
		if !buf.IsClosed() {
			t.Error("buffer must be closed")
		}
	})

	t.Run("Run", func(t *testing.T) {
		stat := IfStat{Delay: time.Millisecond, Out: ioutil.Discard}
		stop := stat.Run()
		time.Sleep(20 * time.Millisecond)
		stop()
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

type errorReader []byte

func (errorReader) Read([]byte) (int, error) { return 0, InterfaceNotExists("") }
func (errorReader) Close() error             { return nil }

func Test_getInt(t *testing.T) {
	tests := []struct {
		name string
		r    offsetReader
		want int
	}{
		{
			name: "one",
			r:    &fakeBuffer{[]byte("1\n")},
			want: 1,
		},
		{
			name: "not a number",
			r:    &fakeBuffer{[]byte("abc\n")},
			want: 0,
		},
		{
			name: "no newline",
			r:    &fakeBuffer{[]byte("123")},
			want: 0,
		},
		{
			name: "no data",
			r:    &errorReader{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInt(tt.r); got != tt.want {
				t.Errorf("getInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
