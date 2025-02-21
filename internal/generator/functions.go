package generator

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
)

func parseQueries(queries []annotation.Query) string {
	values := url.Values{}

	for _, query := range queries {
		values.Add(query.Key, query.Value.String())
	}

	result := fmt.Sprintf("%#v", values)

	for _, query := range queries {
		result = replaceVariable(result, query.Value)
	}

	return result
}

func parseHeaders(headers []annotation.Header) string {
	values := http.Header{}

	for _, header := range headers {
		for _, value := range header.Values {
			values.Add(header.Key, value.String())
		}
	}

	result := fmt.Sprintf("%#v", values)

	for _, header := range headers {
		for _, value := range header.Values {
			result = replaceVariable(result, value)
		}
	}

	return result
}

func parsePathParams(method model.Method) string {
	values := map[string]string{}

	for _, param := range method.Annotations.RequestLine.PathParams() {
		values[param.Escape()] = param.String()
	}

	result := fmt.Sprintf("%#v", values)

	for _, param := range method.Annotations.RequestLine.PathParams() {
		result = replaceVariable(result, param)
	}

	return result
}

// input - is a string value of a map
// value - is a variable param like {owner}
// example:
// input: `map[string]string{"owner": "{owner}", "repo": "{repo}"}`
// ouput: `map[string]string{"owner": owner, "repo": repo}`
func replaceVariable(input string, value annotation.Variable) string {
	if !value.IsVariable() {
		return input
	}

	v := value.Escape()
	return strings.ReplaceAll(input, `"{`+v+`}"`, v)
}
