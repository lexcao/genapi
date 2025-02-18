package annotation

import (
	"errors"
	"fmt"
)

type InterfaceAnnotations struct {
	BaseURL BaseURL
	Headers []Header
}

type MethodAnnotations struct {
	RequestLine RequestLine
	Queries     []Query
	Headers     []Header
}

func ParseInterfaceAnnotations(comments []string) InterfaceAnnotations {
	var annotations InterfaceAnnotations

	parseComments(
		comments,
		&annotations.BaseURL,
		&annotations.Headers,
	)

	return annotations
}

func ParseMethodAnnotations(comments []string) MethodAnnotations {
	var annotations MethodAnnotations

	parseComments(
		comments,
		&annotations.RequestLine,
		&annotations.Queries,
		&annotations.Headers,
	)

	return annotations
}

// parseComments into typed Annotation
// the inputs must be pointers
func parseComments(comments []string, inputs ...any) {
	for _, comment := range comments {
		annotation, err := parse(comment)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				fmt.Printf("[error] parsing annotation: %s\n", err)
			}
			continue
		}

		for _, input := range inputs {
			if err := typed(annotation, input); err != nil {
				if !errors.Is(err, errSkipTyped) {
					fmt.Printf("[error] parsing annotation value @%s: %s\n", annotation.Name, err)
				}
			}
		}
	}
}
