package main

import (
	"crypto/md5"
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

func NewVerifyCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "verify [path1] [path2]",
		Short:        "Verify split files",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("usage: filesplit verify <path1> <path2>")
			}

			err := verifyFile(args[0], args[1])
			if err != nil {
				fmt.Printf("compare %s with %s with err: %s", args[0], args[1], err.Error())
			}
			return err
		},
	}
}

func isDir(file string) (bool, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return true, nil
	}
	return false, nil
}

func verifyFile(path1, path2 string) error {
	md5Func := calculateMD5
	if ok, err := isDir(path1); err == nil {
		if ok {
			fmt.Printf("path %s is directory\n", path1)
			md5Func = computeDirectoryMD5
		}
	} else {
		return err
	}

	md5File1, err := md5Func(path1)
	if err != nil {
		return err
	}
	fmt.Printf("path %s MD5 value: %s\n", path1, md5File1)

	md5Func = calculateMD5
	if ok, err := isDir(path2); err == nil {
		if ok {
			fmt.Printf("path %s is directory\n", path2)
			md5Func = computeDirectoryMD5
		}
	} else {
		return err
	}

	md5File2, err := md5Func(path2)
	if err != nil {
		return err
	}
	fmt.Printf("path %s MD5 value: %s\n", path2, md5File2)
	if md5File1 == md5File2 {
		fmt.Println("Files are same")
		return nil
	}

	fmt.Println("Files are different")
	return fmt.Errorf("MD5 of %s and %s are different", path1, path2)
}

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func computeDirectoryMD5(dirPath string) (string, error) {
	hasher := md5.New()

	matchers, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return "", err
	}

	sort.Slice(matchers, func(i, j int) bool {
		iext := strings.Replace(filepath.Ext(matchers[i]), ".part", "", -1)
		jext := strings.Replace(filepath.Ext(matchers[j]), ".part", "", -1)
		return cast.ToInt(iext) < cast.ToInt(jext)
	})

	for _, m := range matchers {
		fmt.Printf("start to deal file %s\n", m)
		fi, err := os.Stat(m)
		if err != nil {
			return "", err
		}
		if fi.IsDir() {
			fmt.Printf("path %s is directory, skip ...\n", m)
			continue
		}
		file, err := os.Open(m)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
