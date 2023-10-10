package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func main() {
	// ZIP 压缩
	if err := compress("test.zip", []string{"./basic/compress/zip/file1.txt", "./basic/compress/zip/file2.txt"}); err != nil {
		fmt.Println(err)
		return
	}

	// ZIP 解压缩
	if err := extract("test.zip", "."); err != nil {
		fmt.Println(err)
		return
	}
}

// ZIP 压缩函数
func compress(zipFile string, files []string) error {
	// 创建一个新的 ZIP 文件
	newZipFile, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	// 实例化 zipWriter
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 添加文件到 ZIP 文件中
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		// 获取文件信息
		info, err := f.Stat()
		if err != nil {
			return err
		}

		// 创建一个文件头信息
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 添加文件到 ZIP 文件中
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, f); err != nil {
			return err
		}
	}
	fmt.Println("ZIP file created:", zipFile)
	return nil
}

// ZIP 解压缩函数
func extract(zipFile string, dest string) error {
	// 打开 ZIP 文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 解压缩所有文件
	for _, f := range r.File {
		path := fmt.Sprintf("%s/%s", dest, f.Name)

		// 如果是一个目录，则创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		// 创建文件
		if err = os.MkdirAll(path[:len(path)-len(f.Name)], os.ModePerm); err != nil {
			return err
		}
		newFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer newFile.Close()

		// 打开文件
		file, err := f.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// 写入文件内容
		if _, err = io.Copy(newFile, file); err != nil {
			return err
		}
	}
	fmt.Println("ZIP file extracted to:", dest)
	return nil
}
