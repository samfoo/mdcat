package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()

	renderer := &Console{}
	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	if len(args) > 0 {
		input, err := ioutil.ReadFile(args[0])
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("mdcat: %s: unable to read from file\n", args[0]))
			os.Exit(1)
		}

		output := blackfriday.Markdown(input, renderer, extensions)
		os.Stdout.Write(output)
	} else {
		reader := bufio.NewReader(os.Stdin)

		var input []byte
		buffer := make([]byte, 2<<20)

		for {
			count, err := reader.Read(buffer)

			if count == 0 {
				break
			}

			if err != nil {
				os.Stderr.WriteString("mdcat: unable to read from pipe\n")
				os.Exit(1)
			}

			input = append(input, buffer...)
		}

		output := blackfriday.Markdown(input, renderer, extensions)
		os.Stdout.Write(output)
	}

}
