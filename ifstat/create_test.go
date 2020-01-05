package ifstat

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	// logger = log.New(os.Null)
	os.Exit(m.Run())
}

func TestUniq(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			name: "already uniq",
			in:   []string{"1", "2", "3"},
			out:  []string{"1", "2", "3"},
		},
		{
			name: "some equal",
			in:   []string{"1", "2", "2"},
			out:  []string{"1", "2"},
		},
		{
			name: "some equal",
			in:   []string{"1", "1", "2"},
			out:  []string{"1", "2"},
		},
		{
			name: "all equal",
			in:   []string{"2", "2", "2"},
			out:  []string{"2"},
		},
		{
			name: "nothing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := makeUniq(tt.in)
			if !reflect.DeepEqual(res, tt.out) {
				t.Errorf("result does not match.\nExpected:\n%#v\nGiven:\n%#v", tt.out, res)
			}
		})
	}
}

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
			name:   "eno1",
			ifname: "eno1",
		},
		{
			name:   "wlan0",
			ifname: "wlan0",
		},
		{
			name:    "err",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkIfaceExists(tt.ifname); (err != nil) != tt.wantErr {
				t.Errorf("checkIfaceExists() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
	_, _ = NewStat("eno1", "eno1", "", "error")
}
