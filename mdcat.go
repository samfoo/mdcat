package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"
	"github.com/russross/blackfriday"
)

func main() {
	flag.Parse()

	args := flag.Args()

	renderer := &Console{
		tty: isatty.IsTerminal(os.Stdout.Fd()),
	}
	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			input, err := ioutil.ReadFile(args[i])
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("mdcat: %s: unable to read from file\n", args[i]))
				os.Exit(1)
			}

			renderer.base = filepath.Dir(args[i])
			output := blackfriday.Markdown(input, renderer, extensions)
			os.Stdout.Write(output)
		}
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
