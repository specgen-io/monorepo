package imports

import "github.com/specgen-io/specgen/v2/sources"

func GenerateImports(w *sources.Writer, imports []string) {
	for _, imp := range imports {
		w.Line(`import %s;`, imp)
	}
}
