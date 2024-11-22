package archive

import (
	"os"

	"github.com/klauspost/compress/zstd"
)

func WalkTarZstd(file *os.File, walkFunc WalkFunc) error {
	r, err := zstd.NewReader(file)
	if err != nil {
		return err
	}

	return WalkTar(r, walkFunc)
}
