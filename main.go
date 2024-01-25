package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "", true)
}

func printDir(out io.Writer, path string, printFiles bool, prefix string, isLastDir bool) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	var allFiles []os.FileInfo

	for _, file := range files {
		if !printFiles && !file.IsDir() {
			continue
		}
		allFiles = append(allFiles, file)
	}

	sortFilesAndDirs(allFiles)

	for i, file := range allFiles {
		isLastInLevel := i == len(allFiles)-1
		if printFiles || file.IsDir() {
			printElement(out, path, file.Name(), printFiles, prefix, isLastInLevel)
		}
		if file.IsDir() {
			subdirPath := path + string(os.PathSeparator) + file.Name()
			newPrefix := prefix + calculatePrefix(isLastInLevel)
			printDir(out, subdirPath, printFiles, newPrefix, isLastInLevel)
		}
	}

	return nil
}

func printElement(out io.Writer, rootPath, name string, printFiles bool, prefix string, isLastInLevel bool) {
	fmt.Fprint(out, prefix)
	if isLastInLevel {
		fmt.Fprint(out, "└───")
	} else {
		fmt.Fprint(out, "├───")
	}

	fmt.Fprint(out, name)

	if printFiles {
		fileInfo, err := os.Stat(rootPath + string(os.PathSeparator) + name)
		if err == nil && !fileInfo.IsDir() {
			size := fileInfo.Size()
			sizeStr := "empty"
			if size > 0 {
				sizeStr = fmt.Sprintf("%db", size)
			}
			fmt.Fprintf(out, " (%s)", sizeStr)
		}
	}

	fmt.Fprintln(out)
}

func calculatePrefix(isLastInLevel bool) string {
	if isLastInLevel {
		return "\t"
	}
	return "│\t"
}

func sortFilesAndDirs(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}
