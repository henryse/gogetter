// **********************************************************************
//    Copyright (c) 2017-2022 Henry Seurer
//
//   Permission is hereby granted, free of charge, to any person
//    obtaining a copy of this software and associated documentation
//    files (the "Software"), to deal in the Software without
//    restriction, including without limitation the rights to use,
//    copy, modify, merge, publish, distribute, sublicense, and/or sell
//    copies of the Software, and to permit persons to whom the
//    Software is furnished to do so, subject to the following
//    conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//    EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
//    OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//    NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//    HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
//    WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
//    FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//    OTHER DEALINGS IN THE SOFTWARE.
//
// **********************************************************************

// GoGetter is a simple application I use to load the necessary libraries that my
// projects require

package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	// Parse Command line.
	//
	var libraries string
	var sources string

	flag.StringVar(&libraries, "libraries", "", "Location of the libraries.txt file")
	flag.StringVar(&sources, "source", "", "Directory location of the go files to parse.")
	update := flag.Bool("update", false, "Update libraries from sources")
	install := flag.Bool("install", false, "GO GET libraries")
	showVersion := flag.Bool("version", false, "Show version")
	showHelp := flag.Bool("help", false, "Show Help")
	flag.Parse()

	// If no parameters simply look for the libraries.txt file in the gogetter's directory
	if len(libraries) == 0 {
		libraries = path.Join(path.Dir(os.Args[0]), "libraries.txt")
	} else {
		// ok, if they just passed in a directory, then assume the file name is libraries.txt
		if IsDirectory(libraries) {
			libraries = path.Join(libraries, "libraries.txt")
		}
	}

	log.Println("[INFO] libraries = ", libraries)
	log.Println("[INFO] sources = ", sources)
	log.Println("[INFO] update = ", *update)
	log.Println("[INFO] install = ", *install)

	if *showHelp == true {
		ShowHelp()
		return
	}

	if *showVersion == true {
		ShowVersion()
		return
	}

	if !Validate(sources, libraries, *update, *install) {
		log.Fatal("[ERROR] Invalid Parameters.")
	}

	if *update == true {
		UpdateLibraries(sources, libraries)
	}

	if *install == true {
		InstallLibraries(libraries)
	}
}

func Validate(sources string, libraries string, update bool, install bool) bool {
	if !install && !update {
		log.Println("[ERROR] Must specify --update, --install or both --update and --install")
		return false
	}

	if update {
		if len(sources) == 0 || len(libraries) == 0 {
			log.Println("[ERROR] --update requires both --sources and --libraries parameters")
			return false
		}
	}

	if install {
		if len(libraries) == 0 {
			log.Println("[ERROR] --install requires --libraries parameters")
			return false
		}
	}

	return true
}

func ShowVersion() {
	fmt.Println("[INFO] go-getter Version 20220510.1")
}

func ShowHelp() {
	flag.PrintDefaults()
}

func UpdateLibraries(sources string, libraries string) {
	log.Printf("[INFO] Updating library %s from source files at: %s\n", libraries, sources)

	_, err := os.Stat(sources)
	if os.IsNotExist(err) {
		log.Fatal("[ERROR] Folder does not exist.")
		return
	}

	_, err = os.Stat(libraries)
	if os.IsNotExist(err) {
		log.Fatal("[ERROR] Folder does not exist.")
		return
	}

	items, _ := ioutil.ReadDir(sources)
	var imports []string

	for _, item := range items {
		if !item.IsDir() {
			if strings.Compare(filepath.Ext(item.Name()), ".go") == 0 {
				imports = append(imports, ParseFile(filepath.Join(sources, item.Name()))...)
			}
		}
	}

	if len(imports) > 0 {
		WriteLibrariesFile(libraries, imports)
	}
}

func WriteLibrariesFile(libraries string, imports []string) {
	log.Printf("[INFO] Found imports, writing to %s\n", libraries)

	_, err := os.Stat(libraries)
	if os.IsNotExist(err) {
		log.Fatal("[ERROR] Folder does not exist.")
		return
	}

	file, err := os.OpenFile(libraries, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("[ERROR] Failed creating file: %s", err)
	}

	dataWriter := bufio.NewWriter(file)

	for _, data := range imports {
		_, _ = dataWriter.WriteString(data + "\n")
	}

	_ = dataWriter.Flush()
	_ = file.Close()
}

func ParseFile(fileName string) []string {
	var re = regexp.MustCompile(`^([a-z\d]+(-[a-z\d]+)*\.)+[a-z]{2,}$`)

	tokenSet := token.NewFileSet()
	parsed, err := parser.ParseFile(tokenSet, fileName, nil, parser.ImportsOnly)
	if err != nil {
		log.Fatal(err)
	}
	var imports []string
	for _, i := range parsed.Imports {
		p := strings.Trim(i.Path.Value, `"`)
		fmt.Print(p)

		splits := strings.Split(p, "/")

		for _, value := range splits {
			for range re.FindAllString(value, -1) {
				imports = append(imports, p)
			}
		}
	}

	return imports
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		switch mode := info.Mode(); {
		case mode.IsDir():
			return true
		case mode.IsRegular():
			return false
		}
	}

	return false
}

func InstallLibraries(libraries string) {
	log.Printf("[INFO] Installing library %s/libraries.txt\n", libraries)

	file, err := os.Open(libraries)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		library := scanner.Text()

		if library[0] != '#' {
			log.Printf("[INFO] go get %s\n", library)

			cmd := exec.Command("go", "get", library)
			stdout, err := cmd.Output()

			if err != nil {
				log.Fatal(err.Error())
				return
			}

			output := string(stdout)

			if len(output) > 0 {
				log.Println(output)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
