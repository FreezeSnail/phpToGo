package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

type functionObject struct {
	Name     string   `json:"name"`
	Args     []string `json:"args"`
	Rest     string   `json:"rest"`
	Query    string   `json:"query"`
	Contents []string `json:"contents"`
	Comment  string   `json:"comment"`
}

const (
	COMMENT int = 0
	FUNC        = 1
	OTHER       = 2
)

func getArgs() (string, string) {
	var (
		in, out string
	)

	if len(os.Args) == 3 {
		in = os.Args[1]
		out = os.Args[2]
	} else if len(os.Args) == 2 {
		in = os.Args[1]
		out = "."
	} else {
		in = "."
		out = "."
	}

	return in, out
}

func getFiles(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return files
}

func lineType(line string) int {
	if strings.Contains(line, "function") {
		return FUNC
	} else if strings.Contains(line, "/*") {
		return COMMENT
	} else {
		return OTHER
	}
}

func isRestful(line string) bool {
	if strings.Contains(line, "get") {
		return true
	}

	return false
}

func parseFile(file fs.FileInfo) []*functionObject {
	fmt.Println(file.Name())
	name := file.Name()
	if !strings.Contains(name, ".php") {
		return nil
	}
	f, err := os.Open(file.Name())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	var functions []*functionObject
	var currentFunction functionObject
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("current line:", line)

		switch lineType(line) {
		case COMMENT:
			currentFunction = functionObject{}
			functions = append(functions, &currentFunction)

			currentFunction.Comment += line + "\n"

			for !strings.Contains(line, "*/") {
				fmt.Println("current line:", line)

				scanner.Scan()
				line = scanner.Text()
				currentFunction.Comment += line + "\n"
			}

		case FUNC:
			currentFunction.Name = strings.Split(line, " ")[1]
			currentFunction.Args = strings.Split(strings.Split(line, " ")[2], ",")
			for line != "}" {
				fmt.Println("current line:", line)

				scanner.Scan()
				line = scanner.Text()

				if strings.Contains(line, "query") {
					currentFunction.Query = line
				} else {
					currentFunction.Contents = append(currentFunction.Contents, line)
				}
			}

			fmt.Println("CUR FUNC: ", currentFunction)
		}
	}

	return functions
}

func main() {
	in, out := getArgs()
	fmt.Println(in, out)
	files := getFiles(in)
	for _, file := range files {
		//fmt.Println(file.Name(), file.IsDir())
		funcs := parseFile(file)
		for _, fun := range funcs {
			e, err := json.MarshalIndent(*fun, "", "\t")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(e))
		}
	}
}
