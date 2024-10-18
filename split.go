package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func NewSplitCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "split [input_file_path] [chunk_size] [output_dir]",
		Short:        "Split a file into chunks",
		SilenceUsage: true,
		Example: `
	filesplit split input.txt 1024M output
	
	Note: 
		<chun_size> support [KB, MB, GB], max 10GB`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return errors.New("usage: filesplit split <input_file_path> <chunk_size> <output_dir>")
			}

			chunkSize, err := humanize.ParseBytes(args[1])
			if err != nil {
				return err
			}
			if chunkSize > 10*humanize.GByte {
				chunkSize = 10 * humanize.GByte
			}

			return splitFile(args[0], int64(chunkSize), args[2])
		},
	}
}

func splitFile(filePath string, chunkSize int64, outputDir string) error {
	fmt.Println("read file: ", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fmt.Println("file size: ", humanize.Bytes(uint64(fileInfo.Size())))
	fmt.Println("create output dir: ", outputDir)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	fmt.Println("chunk size: ", humanize.Bytes(uint64(chunkSize)))
	fmt.Println("start split file...")
	var chunkNum int64
	totalProgress := fileInfo.Size()/chunkSize + 1
	cnt := 0
	maxBuffSize := 1 * humanize.GByte
	filePos := 0
loop:
	for {
		chunkFilePath := filepath.Join(outputDir, fmt.Sprintf("%s.part%d", fileInfo.Name(), chunkNum))
		chunkFile, err := os.Create(chunkFilePath)
		if err != nil {
			return err
		}

		bytesWritten := int64(0)
		cnt += 1
		fmt.Printf("start progress: %d/%d ...\n", cnt, totalProgress)

		for bytesWritten < chunkSize {
			buffer := make([]byte, maxBuffSize) // 创建一个与 chunkSize 相同大小的缓冲区
			n, err := file.Read(buffer)

			if n > 0 {
				if bytesWritten+int64(n) > chunkSize {
					n = int(chunkSize - bytesWritten) // 确保不超过 chunkSize
				}
				if _, err := chunkFile.Write(buffer[:n]); err != nil {
					chunkFile.Close()
					return err
				}
				bytesWritten += int64(n)
				filePos += n
				file.Seek(int64(filePos), 0)
			}

			if err == io.EOF {
				break loop // 文件读取结束
			}
			if err != nil {
				chunkFile.Close()
				return err
			}
		}

		chunkFile.Close()
		if bytesWritten == 0 {
			break // 没有更多数据可写，结束循环
		}

		chunkNum++
	}

	fmt.Println("split file done!")

	return nil
}
