package archive

import (
	"archive/zip"
	"os"
)

func WalkZip(file *os.File, fileSize int64, walkFunc WalkFunc) error {
	zfs, err := zip.NewReader(file, fileSize)
	if err != nil {
		return err
	}

	for _, f := range zfs.File {
		err = walkZipFile(f, walkFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func walkZipFile(f *zip.File, walkFunc WalkFunc) error {
	zFile, err := f.Open()
	if err != nil {
		err = walkFunc(f.Name, f.FileInfo(), nil, err)
	} else {
		defer zFile.Close()
		err = walkFunc(f.Name, f.FileInfo(), zFile, err)
	}
	return err
}
