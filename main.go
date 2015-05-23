// chris 051215 Simple tee implementation.

// Tee reads data from standard in, copies it to standard out, and
// additionally copies it to the zero or more paths specified as
// command-line arguments.
//
//	Usage: tee [-append] [path ...]
//	  -append=false: Append the output to the files rather than overwriting them.
//	  path: Zero or more files to which to additionally copy standard in.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const bufsize = 4096

func main() {
	optappend := flag.Bool("append", false,
		"Append the output to the files rather than overwriting them.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-append] [path ...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  path: Zero or more files to which to additionally copy standard in.\n")
	}
	flag.Parse()

	// Open output files in addition to standard out.
	out := []*os.File{os.Stdout}
	fileflag := os.O_WRONLY | os.O_CREATE
	if *optappend {
		fileflag |= os.O_APPEND
	} else {
		fileflag |= os.O_TRUNC
	}
	for _, path := range flag.Args() {
		file, err := os.OpenFile(path, fileflag, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		out = append(out, file)
	}

	// Buffer for repeated reads on standard in.
	buf := make([]byte, bufsize)

	// Read from standard in (until eof), and copy data to output files.
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if n == 0 {
				// eof
				break
			}
			log.Fatal(err)
		}
		data := buf[:n] // Write same data to each output file.
		for _, file := range out {
			_, err := file.Write(data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
