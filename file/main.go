package file

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDir(src string, dst string) error {
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcFileInfo.Mode()); err != nil {
		return err
	}

	fileInfos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, fileInfo := range fileInfos {
		srcfp := path.Join(src, fileInfo.Name())
		dstfp := path.Join(dst, fileInfo.Name())
		if fileInfo.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				break
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				break
			}
		}
	}
	return nil
}
