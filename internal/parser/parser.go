package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
)

const EMBEDED_INTERFACE = "genapi.Interface"

// ParseFile parse the given file and return the interface definition
func ParseFile(filename string) ([]model.Interface, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var interfaces []model.Interface

	ast.Inspect(file, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok {
			for _, spec := range decl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						if hasEmbededInterface(iface, EMBEDED_INTERFACE) {
							interfaces = append(interfaces, parseInterface(
								parseInterfaceParams{
									File:        file,
									Declaration: decl,
									TypeSpec:    typeSpec,
									Interface:   iface,
								},
							))
						}
						return false
					}
				}
			}
		}
		return true
	})

	return interfaces, nil
}

func hasEmbededInterface(iface *ast.InterfaceType, typ string) bool {
	for _, field := range iface.Methods.List {
		if _, ok := field.Type.(*ast.SelectorExpr); ok {
			return types.ExprString(field.Type) == typ
		}
	}
	return false
}

type parseInterfaceParams struct {
	File        *ast.File          // for package
	Declaration *ast.GenDecl       // for comments
	TypeSpec    *ast.TypeSpec      // for interface type
	Interface   *ast.InterfaceType // for interface methods
}

func parseInterface(params parseInterfaceParams) model.Interface {
	result := model.Interface{
		Name:    params.TypeSpec.Name.Name,
		Package: params.File.Name.Name,
	}

	result.Annotations = annotation.ParseInterfaceAnnotations(commentsString(params.Declaration.Doc))

	for _, method := range params.Interface.Methods.List {
		if fn, ok := method.Type.(*ast.FuncType); ok {
			parsed := parseMethod(method, fn)
			parsed.Interface = result.Name
			result.Methods = append(result.Methods, parsed)
		}
	}

	return result
}

func parseMethod(method *ast.Field, fn *ast.FuncType) model.Method {
	result := model.Method{
		Name: method.Names[0].Name,
	}

	result.Annotations = annotation.ParseMethodAnnotations(commentsString(method.Doc))
	result.Params = parseParams(fn.Params)
	result.Results = parseParams(fn.Results)

	return result
}

func parseParams(params *ast.FieldList) []model.Param {
	if params == nil {
		return nil
	}

	var result []model.Param
	for _, param := range params.List {
		var name string
		if len(param.Names) > 0 {
			name = param.Names[0].Name
		}

		result = append(result, model.Param{
			Name: name,
			Type: types.ExprString(param.Type),
		})
	}
	return result
}

func commentsString(doc *ast.CommentGroup) []string {
	if doc == nil {
		return nil
	}

	result := make([]string, len(doc.List))
	for i, comment := range doc.List {
		result[i] = comment.Text
	}
	return result
}
