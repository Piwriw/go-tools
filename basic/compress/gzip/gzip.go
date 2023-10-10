package main

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

func compressFile(inputFile, outputFile string) error {
	// 打开输入文件
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	// 创建输出文件
	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	// 创建gzip写入器
	gzipWriter := gzip.NewWriter(output)
	defer gzipWriter.Close()

	// 将输入文件内容写入gzip写入器
	_, err = io.Copy(gzipWriter, input)
	if err != nil {
		return err
	}

	return nil
}

func decompressFile(inputFile, outputFile string) error {
	// 打开输入文件
	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer input.Close()

	// 创建gzip读取器
	gzipReader, err := gzip.NewReader(input)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 创建输出文件
	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	// 将gzip读取器内容写入输出文件
	_, err = io.Copy(output, gzipReader)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// 压缩文件
	err := compressFile("./basic/compress/gzip/input.txt", "compressed.gz")
	if err != nil {
		log.Fatal(err)
	}

	// 解压缩文件
	err = decompressFile("compressed.gz", "./basic/compress/gzip/output.txt")
	if err != nil {
		log.Fatal(err)
	}
}
