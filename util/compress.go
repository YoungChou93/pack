package util

import (
	"os"
	"archive/tar"
	"io"
	"compress/gzip"
	"path"
)

func Untar(source string,targetdir string)error{
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gzReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		target := path.Join(targetdir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, 0755)
			if err != nil {
				return err
			}

			setAttrs(target, header)
			break

		case tar.TypeReg:
			w, err := os.Create(target)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, tarReader)
			if err != nil {
				return err
			}
			w.Close()

			setAttrs(target, header)
			break
		case tar.TypeLink:
			os.Link(header.Linkname,target)
			setAttrs(target, header)
			break
		case tar.TypeSymlink:
			os.Symlink(header.Linkname,target)
			setAttrs(target, header)
			break
		default:

			break
		}
	}

	return nil

}

func setAttrs(target string, header *tar.Header) {
	os.Chmod(target, os.FileMode(header.Mode))
	os.Chtimes(target, header.AccessTime, header.ModTime)
}
