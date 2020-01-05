package ifstat

import (
	"os"
)

func (i InterfaceNotExists) Error() string {
	return "Interface " + string(i) + " does not exists"
}

func (f *fileWithOffset) Read(b []byte) (int, error) {
	return (*os.File)(f).ReadAt(b, 0)
}

func (f *fileWithOffset) Close() error {
	return (*os.File)(f).Close()
}
