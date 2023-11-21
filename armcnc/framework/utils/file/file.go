/**
 ******************************************************************************
 * @file    file.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package FileUtils

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathMkdir(path string) (bool, error) {
	err := os.Mkdir(path, 0666)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func PathMkdirAll(path string) (bool, error) {
	err := os.MkdirAll(path, 0666)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func WriteFile(data string, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return err
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var buffer bytes.Buffer
	chunk := make([]byte, 512*1024)
	for {
		n, err := reader.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		buffer.Write(chunk[:n])
	}
	return buffer.Bytes(), nil
}

func Unzip(src string, dest string, length int) bool {
	status := true
	reader, err := zip.OpenReader(src)
	if err != nil {
		return false
	}
	defer reader.Close()

	for _, file := range reader.File {
		parts := strings.Split(file.Name, "/")
		if len(parts) > length {
			status = false
			break
		}
		if strings.Contains(file.Name, "machine") || strings.Contains(file.Name, "launch") {
			filePath := path.Join(dest, file.Name)
			if file.FileInfo().IsDir() {
				os.MkdirAll(filePath, os.ModePerm)
			} else {
				if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
					status = false
					break
				}
				inFile, err := file.Open()
				if err != nil {
					status = false
					break
				}
				outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				if err != nil {
					inFile.Close()
					status = false
					break
				}
				if _, err = io.Copy(outFile, inFile); err != nil {
					outFile.Close()
					inFile.Close()
					status = false
					break
				}
				outFile.Close()
				inFile.Close()
			}
		} else {
			status = false
			break
		}
	}

	return status
}

func Zip(src string, dest string) bool {
	status := true
	zipFile, err := os.Create(dest)
	if err != nil {
		status = false
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, src)
		if info.IsDir() {
			header.Name += "/"
		}
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}

		return err
	})

	if err != nil {
		status = false
	}

	return status
}

func ZipFiles(src []string, dest string) bool {
	status := true
	zipFile, err := os.Create(dest)
	if err != nil {
		status = false
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range src {
		if err = fileToZip(zipWriter, file); err != nil {
			status = false
		}
	}

	return status
}

func fileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
