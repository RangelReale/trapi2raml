package trapi2raml

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/RangelReale/trapi"
	"github.com/davecgh/go-spew/spew"
)

type wrapWriter struct {
	writer   io.Writer
	err      error
	replacer *strings.Replacer
}

func newWrapWriter(writer io.Writer) *wrapWriter {
	return &wrapWriter{
		writer:   writer,
		replacer: strings.NewReplacer("\n", "", "\r", ""),
	}
}

func (w *wrapWriter) Err() error {
	return w.err
}

func (w *wrapWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err = w.writer.Write(p)
	if err != nil {
		w.err = err
	}

	return
}

func (w *wrapWriter) writeLine(ident int, text string) {
	if ident > 0 {
		io.WriteString(w, strings.Repeat("  ", ident))
	}
	io.WriteString(w, w.replacer.Replace(text))
	io.WriteString(w, "\n")
}

func (w *wrapWriter) writeLineMultiline(ident int, text string) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		w.writeLine(ident, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		w.err = err
	}
}

func (w *wrapWriter) writeType(ident int, dt *trapi.ApiDataType) {
	w.writeTypeInternal(ident, dt, true)
}

func (w *wrapWriter) writeTypeInternal(ident int, dt *trapi.ApiDataType, isroot bool) {

	/*
		if !isroot {
			w.writeLine(ident, fmt.Sprintf("%s:", dt.Name))
		}
	*/
	w.writeLine(ident+1, fmt.Sprintf("type: %s", TRType(dt.DataType)))
	w.writeLine(ident+1, fmt.Sprintf("debugtype: %s", spew.Sdump(dt)))
	if dt.Description != "" {
		w.writeLine(ident+1, fmt.Sprintf("description: %s", dt.Description))
	}
	if !isroot && !dt.Required {
		w.writeLine(ident+1, "required: false")
	}
	if dt.DataType == trapi.DATATYPE_OBJECT && dt.Items != nil {
		w.writeLine(ident+1, "properties:")
		for _, iord := range dt.ItemsOrder {
			w.writeLine(ident+2, fmt.Sprintf("%s:", iord))
			w.writeTypeInternal(ident+2, dt.Items[iord], false)
		}
	}
	if dt.Examples != nil && len(dt.Examples) > 0 {
		if len(dt.Examples) == 1 {
			w.writeLine(ident+1, "example:")
			w.writeLineMultiline(ident+2, unidentText(dt.Examples[0].Text))
		} else {
			w.writeLine(ident+1, "examples:")
			for ect, e := range dt.Examples {
				w.writeLine(ident+2, fmt.Sprintf("example%d: |", ect))
				w.writeLineMultiline(ident+3, unidentText(e.Text))
			}
		}
	}
}
