// **********************************************************************
//    Copyright (c) 2017 Henry Seurer
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
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {

	// Parse Command line.
	//
	var libraries string
	var sources string

	flag.StringVar(&libraries, "libraries", "", "Location of the libraries file")
	flag.StringVar(&sources, "source", "", "Directory location of the go files to parse.")
	update := flag.Bool("update", false, "Update libraries from sources")
	install := flag.Bool("install", false, "GO GET libraries")
	showVersion := flag.Bool("version", false, "Show version")
	showHelp := flag.Bool("help", false, "Show Help")
	flag.Parse()

	if *showVersion == true {
		ShowVersion()
	}

	if *showHelp == true {
		ShowHelp()
	}

	if *update == true {
		UpdateLibraries(sources, libraries)
	}

	if *install == true {
		InstallLibraries(libraries)
	}
}

func ShowVersion() {
	fmt.Println("go-getter Version 20220510.1")
}

func ShowHelp() {
	flag.PrintDefaults()
}

func UpdateLibraries(sources string, libraries string) {

	fmt.Println(sources)
	fmt.Println(libraries)
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
	// If no parameters simply look for the libraries.txt file in the gogetter's directory
	if len(libraries) == 0 {
		libraries = path.Join(path.Dir(os.Args[0]), "libraries.txt")
	} else {
		// ok, if they just passed in a directory, then assume the file name is libraries.txt
		if IsDirectory(libraries) {
			libraries = path.Join(libraries, "libraries.txt")
		}
	}

	log.Printf("Installing libraries from: '%s'\n", libraries)

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
			log.Printf("go get %s\n", library)

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
