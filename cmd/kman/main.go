package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kowala-tech/kman"
	"github.com/sanity-io/litter"
)

var parseGo = flag.Bool("go", false, "Parse Go files")
var parseMd = flag.Bool("md", false, "Parse Markdown files")

func main() {

	flag.Parse()

	var assemblers []kman.Assembler

	if *parseGo {
		assemblers = append(assemblers, kman.NewGoAssemblerFromLocalFilesystem())
	}

	if *parseMd {
		assemblers = append(assemblers, kman.NewMarkdownAssemblerFromLocalFilesystem())
	}

	docker := kman.NewDefaultDocumenter(kman.NewDefaultSorter(), assemblers...)

	doc, err := docker.Document()

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	litter.Dump(doc)
}
