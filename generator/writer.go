package trapi2ramlgen

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/RangelReale/trapi"
)

type wrapWriter struct {
	writer   io.Writer
	err      error
	replacer *strings.Replacer

	warnings []error
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
	w.writeTypeInternal(ident, dt, true, false, false)
}

func (w *wrapWriter) writeTypeDefine(ident int, dt *trapi.ApiDataType) {
	w.writeTypeInternal(ident, dt, true, false, true)
}

func (w *wrapWriter) writeTypeInternal(ident int, dt *trapi.ApiDataType, is_root bool, is_required bool, is_define bool) {

	/*
		if !isroot {
			w.writeLine(ident, fmt.Sprintf("%s:", dt.Name))
		}
	*/
	if dt.BuiltIn || dt.DataTypeName == "" || is_define {
		w.writeLine(ident+1, fmt.Sprintf("type: %s", TRType(dt.DataType)))
	} else {
		w.writeLine(ident+1, fmt.Sprintf("type: %s", dt.DataTypeName))
	}
	if dt.Description != "" {
		w.writeLine(ident+1, fmt.Sprintf("description: %s", dt.Description))
	}
	if !is_root && !is_required {
		w.writeLine(ident+1, "required: false")
	}
	if dt.DataType == trapi.DATATYPE_OBJECT && dt.Items != nil {
		w.writeLine(ident+1, "properties:")
		for _, iord := range dt.ItemsOrder {
			w.writeLine(ident+2, fmt.Sprintf("%s:", dt.Items[iord].FieldName))
			w.writeTypeInternal(ident+2, dt.Items[iord].ApiDataType, false, dt.Items[iord].Required, false)
		}
	}
	w.writeExamples(ident+1, dt.Examples)
	/*
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
	*/
}

func (w *wrapWriter) writeHeaders(ident int, hl *trapi.ApiHeaderList) {

	if hl.List != nil && len(hl.List) > 0 {
		for _, qheadname := range hl.Order {
			qhead := hl.List[qheadname]
			if len(qhead) > 0 {
				w.writeLine(ident, fmt.Sprintf("%s:", qheadname))
				if qhead[0].Description != "" {
					w.writeLine(ident+1, fmt.Sprintf("description: %s", qhead[0].Description))
				}
				w.writeTypeInternal(ident, qhead[0].DataType, true, false, true)
			}
		}
	}

}

func (w *wrapWriter) writeExamples(ident int, ex []*trapi.ApiExample) {
	if ex != nil && len(ex) > 0 {
		if len(ex) == 1 {
			w.writeLine(ident, "example:")
			w.writeLineMultiline(ident+1, w.unidentTypedText(ex[0].ContentType, ex[0].Text))
		} else {
			w.writeLine(ident, "examples:")
			for ect, e := range ex {
				w.writeLine(ident+1, fmt.Sprintf("example%d: |", ect))
				w.writeLineMultiline(ident+2, w.unidentTypedText(e.ContentType, e.Text))
			}
		}
	}
}
