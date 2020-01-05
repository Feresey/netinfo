package ifstat

import (
	"bytes"
	"os"
	"testing"
)

func Test_fileWithOffset_Read(t *testing.T) {
	file, err := os.Create("testfile")
	f := (*fileWithOffset)(file)
	if err != nil {
		t.Error(err)
		return
	}
	data := []byte("abracadabra")
	_, err = file.Write(data)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(file.Name())

	tests := []struct {
		name    string
		data    []byte
		want    int
		wantErr bool
	}{
		{
			data: make([]byte, len(data)),
			want: len(data),
		},
		{
			data: make([]byte, len(data)),
			want: len(data),
		},
		{
			data: make([]byte, len(data)),
			want: len(data),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.Read(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileWithOffset.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileWithOffset.Read() = %v, want %v", got, tt.want)
			}
			if !bytes.Equal(tt.data, data) {
				t.Errorf("Received wrong data.\nGiven:\n%s\nExpected:\n%s", string(tt.data), string(data))
			}
		})
	}
	f.Close()
}
