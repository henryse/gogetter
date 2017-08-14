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
ifdef VERSION
	project_version:=$(VERSION)
else
	project_version:=$(shell git rev-parse --short=8 HEAD)
endif

ifdef PROJECT_NAME
	project_name:=$(PROJECT_NAME)
else
	project_name:=$(shell basename $(CURDIR))
endif

version:
	@echo [INFO] [version]
	@echo [INFO]    Go Makefile Version 1.0
	@echo

settings: version
	@echo [INFO] [settings]
	@echo [INFO]    project_version=$(project_version)
	@echo [INFO]    project_name=$(project_name)
	@echo

help: settings
	@printf "\e[1;34m[INFO] [information]\e[00m\n\n"
	@echo [INFO] This make process supports the following targets:
	@echo [INFO]    clean       - clean up and targets in project
	@echo [INFO]    build       - build both the project and Docker image
	@echo [INFO]    push        - push image to repository
	@echo
	@echo [INFO] The script supports the following parameters:
	@echo [INFO]    VERSION      - version to tag docker image wth, default value is the git hash
	@echo [INFO]    PROJECT_NAME - project name, default is git project name
	@echo

libraries:
	@printf "\e[1;34m[INFO] Installing  libraries\e[00m\n\n"
	go get github.com/google/uuid
	go get gopkg.in/yaml.v2

build: settings libraries
	@printf "\e[1;34m[INFO] Building gogetter\e[00m\n\n"
	go build gogetter.go

clean: settings
	@printf "\e[1;34m[INFO] Cleaning up\e[00m\n\n"
	rm gogetter

push: settings
	@printf "\e[1;34m[INFO] Nothing to push\e[00m\n\n"

install: settings build
	@printf "\e[1;34m[INFO] Installing application\e[00m\n\n"
	go install
