package main

import (
	"gset/internal/parser"
	"log"
	"os"
)

func main() {
	arr, err := parser.ParseFile("internal/testing/book.go")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, line := range arr {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

}
