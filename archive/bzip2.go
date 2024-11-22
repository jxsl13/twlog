package archive

import (
	"compress/bzip2"
	"os"
)

func WalkTarBzip2(file *os.File, walkFunc WalkFunc) error {
	r := bzip2.NewReader(file)
	return WalkTar(r, walkFunc)
}
