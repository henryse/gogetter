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
	"flag"
	"os"
	"path"
	"log"
	"bufio"
	"os/exec"
	"fmt"
)

func main() {

	var librariesLocation string

	flag.StringVar(&librariesLocation, "libraries", "", "Directory location of the libraries.txt file OR the libraries file to use.")
	boolPtr := flag.Bool("version", false, "a bool")
	flag.Parse()

	if *boolPtr == true{
		fmt.Println("Go Getter Version 1.0")
	} else {
		// If no parameters simply look for the libraries.txt file in the gogetter's directory
		if len(librariesLocation) == 0 {
			librariesLocation = path.Join(path.Dir(os.Args[0]), "libraries.txt")
		} else {
			// ok, if they just passed in a directory, then assume the file name is libraries.txt
			if IsDirectory(librariesLocation){
				librariesLocation = path.Join(librariesLocation, "libraries.txt")
			}
		}

		log.Printf("Installing libraries from: '%s'\n", librariesLocation)

		file, err := os.Open(librariesLocation)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			GoGet(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func IsDirectory(path string) (bool) {
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

// Heart of the tool, simply execute the "go get" command for each line of the libraries.txt file
func GoGet(library string) {
	if library[0] != '#' {
		log.Printf("go get %s\n", library)

		cmd := exec.Command("go", "get", library)
		stdout, err := cmd.Output()

		if err != nil {
			log.Fatal(err.Error())
			return
		}

		output := string(stdout)

		if len(output) >	 0 {
			log.Println(output)
		}
	}
}