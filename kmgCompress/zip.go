package kmgCompress

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgFile"
)

func ZipUncompressFromBytesToDir(zipB []byte, dir string, trimPrefix string) (err error) {
	buf := bytes.NewReader(zipB)
	reader, err := zip.NewReader(buf, int64(len(zipB)))
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	for _, file := range reader.File {

		fullPath := filepath.Join(dir, strings.TrimPrefix(file.Name, trimPrefix))
		if file.FileInfo().IsDir() {
			err = kmgFile.Mkdir(fullPath)
			if err != nil {
				kmgErr.LogErrorWithStack(err)
				return
			}
			continue
		}
		err = kmgFile.MkdirForFile(fullPath)
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			return
		}
		rc, err := file.Open()
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			return err
		}
		f, err := os.Create(fullPath)
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			rc.Close()
			return err
		}
		_, err = io.Copy(f, rc)
		rc.Close()
		f.Close()
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			return err
		}
	}
	return nil
}
