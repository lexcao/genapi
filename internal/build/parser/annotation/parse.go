package annotation

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrNotFound = errors.New("annotation not found")

type ErrInvalidFormat struct {
	Message string
	Source  string
}

func (e ErrInvalidFormat) Error() string {
	return fmt.Sprintf("invalid annotation format: %s (source: %s)", e.Message, e.Source)
}

// parse parses the annotation format for the given value
// Annotation format: @Name("Value1", "Value2", ...)
// - StartWith `@`
// - Followed by `Name("Value1", "Value2", ...)`
// - Name should be characters and case-insensitive
// - Value should be quoted string `"Value"`
// - Value is optional array list, sperated by comma
func parse(input string) (Annotation, error) {
	r := strings.NewReader(input)

	// find first `@`
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return Annotation{}, ErrNotFound
		}
		if ch == '@' {
			break
		}
	}

	// parse Name
	var name strings.Builder
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return Annotation{}, ErrInvalidFormat{
				Message: "empty annotation name",
				Source:  name.String(),
			}
		}
		if ch == '(' {
			break
		}
		name.WriteRune(ch)
	}

	if name.Len() == 0 {
		return Annotation{}, ErrInvalidFormat{
			Message: "empty annotation name",
			Source:  input,
		}
	}

	// Read annotation values
	var values []string
	var currentValue strings.Builder
	inQuote := false

	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			if inQuote {
				return Annotation{}, ErrInvalidFormat{
					Message: "unclosed quote in value",
					Source:  input,
				}
			}
			break
		}

		switch ch {
		case '"':
			inQuote = !inQuote
			if !inQuote && currentValue.Len() > 0 {
				values = append(values, currentValue.String())
				currentValue.Reset()
			}
		case ',':
			if !inQuote {
				if currentValue.Len() > 0 {
					return Annotation{}, ErrInvalidFormat{
						Message: "value must be quoted",
						Source:  currentValue.String(),
					}
				}
				// Skip comma and any following whitespace
				for {
					ch, _, err := r.ReadRune()
					if err == io.EOF {
						break
					}
					if ch != ' ' {
						if err := r.UnreadRune(); err != nil {
							return Annotation{}, err
						}
						break
					}
				}
				continue
			}
			fallthrough
		default:
			if inQuote {
				currentValue.WriteRune(ch)
			} else if ch != ' ' && ch != ')' {
				currentValue.WriteRune(ch)
				for {
					ch, _, err := r.ReadRune()
					if err == io.EOF || ch == ')' || ch == ',' {
						return Annotation{}, ErrInvalidFormat{
							Message: "value must be quoted",
							Source:  currentValue.String(),
						}
					}
					if ch != ' ' {
						currentValue.WriteRune(ch)
					}
				}
			}
		}
	}

	return Annotation{
		Name:   name.String(),
		Values: values,
	}, nil
}

type Annotation struct {
	Name   string
	Values []string
}
