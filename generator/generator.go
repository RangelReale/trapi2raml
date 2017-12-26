package trapi2ramlgen

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/RangelReale/trapi"
	"sort"
)

type Generator struct {
	ContentType string
	Warnings    []error
}

func NewGenerator() *Generator {
	return &Generator{
		ContentType: "application/json",
	}
}

func (g *Generator) Generate(parser *trapi.Parser, out io.Writer) error {

	ww := newWrapWriter(out)

	ww.writeLine(0, "#%RAML 1.0")
	ww.writeLine(0, fmt.Sprintf("mediaType: %s", g.ContentType))
	ww.writeLine(0, "types:")

	// sort datatype map keys
	/*
		dtkeys := make([]string, 0)
		for dtk, _ := range parser.DataTypes {
			dtkeys = append(dtkeys, dtk)
		}
		sort.Strings(dtkeys)

		for _, dtk := range dtkeys {
			dt := parser.DataTypes[dtk]
			if !dt.BuiltIn {
				ww.writeLine(1, fmt.Sprintf("%s:", dtk))
				ww.writeType(1, dt)
			}
		}
	*/

	// api defines
	for _, def := range parser.ApiDefines {
		ww.writeLine(1, fmt.Sprintf("%s:", def.Name))
		ww.writeTypeDefine(1, def.DataType)
	}

	if ww.Err() != nil {
		return ww.Err()
	}

	// apis
	apilist := parser.BuildApiList()
	err := g.writeApi(ww, 0, "", apilist)
	if err != nil {
		return err
	}

	if ww.Err() != nil {
		return ww.Err()
	}

	g.Warnings = ww.warnings

	return nil
}

func (g *Generator) writeApi(ww *wrapWriter, ident int, curpath string, apilist *trapi.ApiList) error {

	apipath := path.Join(curpath, apilist.Path)
	is_first := apipath == "/"

	newident := ident
	if !is_first {
		// replace <param> with {param}
		param_repl := strings.NewReplacer("<", "{", ">", "}")
		ww.writeLine(ident, fmt.Sprintf("/%s:", param_repl.Replace(apilist.Path)))
		newident = ident + 1
	}

	for _, api := range apilist.Apis {
		ww.writeLine(ident+1, fmt.Sprintf("%s:", strings.ToLower(api.Method)))
		if api.Description != "" {
			ww.writeLine(ident+2, fmt.Sprintf("description: %s", api.Description))
		}
		// headers
		if api.Headers != nil && api.Headers.List != nil && len(api.Headers.List) > 0 {
			ww.writeLine(ident+2, "headers:")
			ww.writeHeaders(ident+3, api.Headers)
		}
		// uri params
		if uparams, ok := api.Params[trapi.PARAMTYPE_URI]; ok {
			ww.writeLine(ident+2, "uriParameters:")
			for _, po := range uparams.Order {
				p := uparams.List[po]
				ww.writeLine(ident+3, fmt.Sprintf("%s:", p.Name))
				ww.writeType(ident+3, p.DataType)
			}

		}
		// query params
		if qparams, ok := api.Params[trapi.PARAMTYPE_QUERY]; ok {
			ww.writeLine(ident+2, "queryParameters:")
			for _, po := range qparams.Order {
				p := qparams.List[po]
				ww.writeLine(ident+3, fmt.Sprintf("%s:", p.Name))
				ww.writeType(ident+3, p.DataType)
			}

		}

		if api.Responses != nil && len(api.Responses.List) > 0 {
			//
			// responses
			//
			ww.writeLine(ident+2, "responses:")
			// sort codes
			qcodelist := make([]string, 0)
			for qcode, _ := range api.Responses.List {
				qcodelist = append(qcodelist, qcode)
			}
			sort.Strings(qcodelist)
			for _, qcode := range qcodelist {
				qrespbodylist := api.Responses.List[qcode]
				ww.writeLine(ident+3, fmt.Sprintf("%s:", qcode))
				if len(qrespbodylist) > 0 {
					ww.writeLine(ident+4, "body:")
					check_repeat := make(map[string]bool)
					for _, qrespbody := range qrespbodylist {
						contenttype := qrespbody.ContentType
						if contenttype == "-" {
							contenttype = g.ContentType
						}
						if _, is_repeat := check_repeat[contenttype]; is_repeat {
							return fmt.Errorf("Duplicate content-type '%s' for responde code '%s' on api '%s %s'", contenttype, qcode, api.Method, apipath)
						}
						check_repeat[contenttype] = true

						// content type
						ww.writeLine(ident+5, fmt.Sprintf("%s:", contenttype))
						// data type
						ww.writeTypeDefine(ident+6, qrespbody.ApiResponse.DataType)
						// headers
						if qrespbody.ApiResponse.Headers != nil && qrespbody.ApiResponse.Headers.List != nil && len(qrespbody.ApiResponse.Headers.List) > 0 {
							ww.writeLine(ident+7, "headers:")
							ww.writeHeaders(ident+8, qrespbody.ApiResponse.Headers)
						}
						// examples
						if qrespbody.ApiResponse.Examples != nil && len(qrespbody.ApiResponse.Examples) > 0 {
							ww.writeExamples(ident+7, qrespbody.ApiResponse.Examples)
						}
					}
				}
			}
		}
	}
	for _, asi := range apilist.SubItems {
		err := g.writeApi(ww, newident, apipath, asi)
		if err != nil {
			return err
		}
	}

	return nil
}
