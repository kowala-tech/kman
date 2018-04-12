package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kowala-tech/kman"
	"github.com/spf13/afero"
)

var (
	parseGo      = flag.Bool("go", false, "Parse Go files")
	parseMd      = flag.Bool("md", true, "Parse Markdown files")
	templatePath = flag.String("theme", "themes/kman", "Theme path")
	outputPath   = flag.String("output", "public", "Public assets output path")
	httpAddress  = flag.String("http", "", "Serve http on a given address (for example, :8080)")
)

func main() {

	flag.Parse()

	if err := build(); err != nil {
		log.Fatal(err)
	}

	if *httpAddress != "" {

		server := http.FileServer(http.Dir(*outputPath))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			if err := build(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				server.ServeHTTP(w, r)
			}
		})

		log.Printf("Serving documentation on %s\n", *httpAddress)
		log.Fatal(http.ListenAndServe(*httpAddress, nil))
	}
}

func build() error {

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
		return fmt.Errorf("Error 01: %s", err)
	}

	renderer := kman.NewRendererAce(
		afero.NewOsFs(),
		*templatePath,
		*outputPath,
	)

	if err := renderer.Render(doc); err != nil {
		return fmt.Errorf("Error 02: %s", err)
	}

	return nil
}
