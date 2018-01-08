package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RangelReale/gocompar"
	"github.com/RangelReale/gocompar/clike"
	"github.com/RangelReale/trapi"
	"github.com/RangelReale/trapi2raml/generator"
	"strings"
)

var tags = flag.String("tags", "", "Filter tags")

// trapi2raml <source-code-path>... <output-filename.raml>
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

	dstFile := flag.Arg(flag.NArg()-1)

	// create c-like comments parser
	clike_parser := gocompar.NewParser(gcp_clike.NewParser(), &gcp_clike.Filter_Golang{})

	// create trapi parser
	parser := trapi.NewParser(clike_parser)

	if tags != nil && *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			fmt.Printf("Tag: %s\n", t)
			parser.AddTag(t)
		}
	}

	// parse requested directories
	for si := 0; si < flag.NArg()-1; si++ {
		srcPath := flag.Arg(si)
		if fi, err := os.Stat(srcPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening source path: %v\n", err)
			os.Exit(1)
		} else if !fi.IsDir() {
			fmt.Fprintf(os.Stderr, "Source path is not a directory: %s\n", srcPath)
			os.Exit(1)
		}

		parser.AddDir(srcPath)
	}

	// create output file
	file, err := os.Create(dstFile)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error creating destination file: %v\n", err)
		os.Exit(1)
	}

	err = parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing source code directories: %v\n", err)
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
