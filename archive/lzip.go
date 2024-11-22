package archive

import (
	"os"

	"github.com/sorairolake/lzip-go"
)

func WalkTarLz(file *os.File, walkFunc WalkFunc) error {
	r, err := lzip.NewReader(file)
	if err != nil {
		return err
	}

	return WalkTar(r, walkFunc)
}
