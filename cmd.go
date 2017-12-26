package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RangelReale/gocompar"
	"github.com/RangelReale/gocompar/clike"
	"github.com/RangelReale/trapi"
	"github.com/RangelReale/trapi2raml/generator"
)

// trapi2raml <source-code-path> <output-filename.raml>
func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "trapi2raml <source-code-path> <output-filename.raml>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// parse parameters
	if len(flag.Args()) < 2 {
		fmt.Fprint(os.Stderr, "Invalid number of parameters\n")
		flag.Usage()
		os.Exit(1)
	}

	srcPath := flag.Arg(0)
	dstFile := flag.Arg(1)

	if fi, err := os.Stat(srcPath); err != nil {
		fmt.Fprint(os.Stderr, "Error opening source path: %v\n", err)
		os.Exit(1)
	} else if !fi.IsDir() {
		fmt.Fprint(os.Stderr, "Source path is not a directory: %s\n", srcPath)
		os.Exit(1)
	}

	// create output file
	file, err := os.Create(dstFile)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error creating destination file: %v\n", err)
		os.Exit(1)
	}

	// create c-like comments parser
	clike_parser := gocompar.NewParser(gcp_clike.NewParser(), &gcp_clike.Filter_Golang{})

	// create trapi parser
	parser := trapi.NewParser(clike_parser)

	// parse requested directory
	parser.AddDir(srcPath)
	err = parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing source code directory %s: %v\n", srcPath, err)
		os.Exit(1)
	}

	// generate RAML file
	gen := trapi2ramlgen.NewGenerator()
	err = gen.Generate(parser, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating RAML file: %v\n", err)
		os.Exit(1)
	}

	// warnings
	if gen.Warnings != nil && len(gen.Warnings) > 0 {
		for _, w := range gen.Warnings {
			fmt.Printf("WARNING: %v\n", w)
		}
	}

	fmt.Printf("File %s generated successfully\n", dstFile)
}
