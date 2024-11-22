package archive

import (
	"os"

	"github.com/bodgit/sevenzip"
)

func Walk7Zip(file *os.File, fileSize int64, walkFunc WalkFunc) error {
	zfs, err := sevenzip.NewReader(file, fileSize)
	if err != nil {
		return err
	}

	for _, f := range zfs.File {
		err = walk7ZipFile(f, walkFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func walk7ZipFile(f *sevenzip.File, walkFunc WalkFunc) error {
	zFile, err := f.Open()
	if err != nil {
		err = walkFunc(f.Name, f.FileInfo(), nil, err)
	} else {
		err = walkFunc(f.Name, f.FileInfo(), zFile, err)
		defer zFile.Close()
	}
	return err
}
