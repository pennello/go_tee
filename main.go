// chris 051215 Simple tee implementation.

package main

import (
	"log"
	"os"
)

const bufsize = 4096

func main() {
	out := []*os.File{os.Stdout}
	for _, path := range os.Args[1:] {
		file, err := os.Create(path)
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
