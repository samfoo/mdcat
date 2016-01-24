package main

import (
	"flag"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	input, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalf("unable to read from file: %s", args[0])
	}

	renderer := &Console{}
	extensions := blackfriday.EXTENSION_STRIKETHROUGH
	output := blackfriday.Markdown(input, renderer, extensions)

	os.Stdout.Write(output)
}
