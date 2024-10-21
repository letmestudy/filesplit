package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func NewMergeCommand() *cobra.Command {
	c := &cobra.Command{
		Use:          "merge [output_file_path] [input_dir]",
		Short:        "Merge split files into a single file",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("usage: filesplit merge <output_file_path> <input_dir>")
			}

			return mergeFiles(args[0], args[1])
		},
	}

	return c
}

func mergeFiles(outputFilePath string, inputDir string) error {
	if _, err := os.Stat(inputDir); err != nil {
		return err
	}

	fmt.Println("merge files in: ", inputDir)

	files, err := filepath.Glob(filepath.Join(inputDir, "*.part*"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no files to merge")
	}

	sort.Slice(files, func(i, j int) bool {
		iext := strings.Replace(filepath.Ext(files[i]), ".part", "", -1)
		jext := strings.Replace(filepath.Ext(files[j]), ".part", "", -1)
		return cast.ToInt(iext) < cast.ToInt(jext)
	})

	fmt.Println("output to: ", outputFilePath)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for idx, filePath := range files {
		fmt.Printf("start progress: %d/%d, file %s\n", idx+1, len(files), filePath)
		inputFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		if _, err := io.Copy(outputFile, inputFile); err != nil {
			return err
		}
	}
	fmt.Println("merge files successfully!")

	return nil
}
