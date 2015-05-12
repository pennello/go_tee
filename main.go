// chris 051215 Simple tee implementation.

package main

import (
	"flag"
	"log"
	"os"
)

const bufsize = 4096

func main() {
	appendopt := flag.Bool("append", false,
		"Append the output to the files rather than overwriting them.")
	flag.Parse()
	out := []*os.File{os.Stdout}
	fileflag := os.O_WRONLY | os.O_CREATE
	if *appendopt {
		fileflag |= os.O_APPEND
	} else {
		fileflag |= os.O_TRUNC
	}
	for _, path := range os.Args[1:] {
		file, err := os.OpenFile(path, fileflag, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		out = append(out, file)
	}
	buf := make([]byte, bufsize)
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
