package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Unzip(zipFileName string, targetDir string) error {
	readerCloser, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer readerCloser.Close()
	var decodeName string
	for _, f := range readerCloser.File {
		if f.Flags == 0 {
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			decodeName = f.Name
		}
		fpath := filepath.Join(targetDir, decodeName)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Zip(path string, target string) error {
	dirPath := filepath.Dir(target)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
	}
	if err != nil {
		return err
	}
	err = zipFiles(path, "", target, nil)
	return err
}

func zipFiles(path string, subPath string, target string, wzip *zip.Writer) error {

	files, err := ioutil.ReadDir(path + "/" + subPath)
	if err == nil {
		if subPath == "" {
			zipfile, err := os.Create(target)
			if err != nil {
				return err
			}
			defer zipfile.Close()
			wzip = zip.NewWriter(zipfile)
			defer wzip.Close()
		}
		for _, file := range files {
			if file.IsDir() {
				err = zipFiles(path, subPath+file.Name()+"/", target, wzip)
				if err != nil {
					return err
				}
			} else {
				a := strings.Replace(os.Args[0], filepath.Join(os.Args[0], ".."), "", 1)
				a = strings.Replace(a, "/", "", 1)
				a = strings.Replace(a, "\\", "", 1)
				if file.Name() == a {
					continue
				}

				header := &zip.FileHeader{
					Name:   subPath + file.Name(),
					Flags:  1 << 11,
					Method: zip.Deflate,
				}
				item, err := wzip.CreateHeader(header)
				if err != nil {
					return err
				}
				datas, _ := ioutil.ReadFile(path + "/" + subPath + file.Name())
				_, err = item.Write(datas)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return err
}
