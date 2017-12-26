package trapi2raml

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/RangelReale/trapi"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(parser *trapi.Parser, out io.Writer) error {

	ww := newWrapWriter(out)

	ww.writeLine(0, "#%RAML 1.0")
	ww.writeLine(0, "mediaType: application/json")
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

	// apis
	apilist := parser.BuildApiList()
	fmt.Printf("%+v\n", apilist)
	g.writeApi(ww, 0, "", apilist)

	return nil
}

func (g *Generator) writeApi(ww *wrapWriter, ident int, curpath string, apilist *trapi.ApiList) {

	apipath := path.Join(curpath, apilist.Path)
	is_first := apipath == "/"

	newident := ident
	if !is_first {
		ww.writeLine(ident, fmt.Sprintf("/%s:", apilist.Path))
		newident = ident + 1
	}

	for _, api := range apilist.Apis {
		ww.writeLine(ident+1, fmt.Sprintf("%s:", strings.ToLower(api.Method)))
		if api.Description != "" {
			ww.writeLine(ident+2, fmt.Sprintf("description: %s", api.Description))
		}
		if uparams, ok := api.Params[trapi.PARAMTYPE_URI]; ok {
			ww.writeLine(ident+2, "uriParameters:")
			for _, po := range uparams.Order {
				p := uparams.List[po]
				ww.writeLine(ident+3, fmt.Sprintf("%s:", p.Name))
				ww.writeType(ident+3, p.DataType)
			}

		}
		if qparams, ok := api.Params[trapi.PARAMTYPE_QUERY]; ok {
			ww.writeLine(ident+2, "queryParameters:")
			for _, po := range qparams.Order {
				p := qparams.List[po]
				ww.writeLine(ident+3, fmt.Sprintf("%s:", p.Name))
				ww.writeType(ident+3, p.DataType)
			}

		}
	}
	for _, asi := range apilist.SubItems {
		g.writeApi(ww, newident, apipath, asi)
	}

}
