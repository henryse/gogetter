# Introduction

This is a simple application I use to install all of the go libraries I need for my applications.

    go build gogetter.go
    go install
    
    OR
    
    make install

The command line parameters are as follows:

      --help string
        Show this help
      --libraries string
        Directory location of the libraries.txt file OR the libraries file to use.

If you do not specify a "--libraries" parameter it will assume the local directory.  If you enter a directory path it will look for a libraries.txt file and finally if you enter a fully qualified path to a text file it use use that file.

For example to load a libraries.txt in a directory you can do the following:
 
    gogetter --libraries=/this/is/a/directory
    
For a specific file: 

    gogetter --libraries=/this/is/a/file.txt

The "libraries.txt" file contains a list packages to install.

You can use # to comment out lines.