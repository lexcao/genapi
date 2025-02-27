package annotation

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

type Variable string

func (v Variable) IsVariable() bool {
	return strings.HasPrefix(string(v), "{") && strings.HasSuffix(string(v), "}")
}

func (v Variable) Escape() string {
	if !v.IsVariable() {
		return string(v)
	}

	return string(v)[1 : len(string(v))-1]
}

func (v Variable) String() string {
	return string(v)
}

type BaseURL struct {
	Value string
}

func (b BaseURL) name() string {
	return "BaseURL"
}

func (b BaseURL) from(annotation Annotation) (any, error) {
	if len(annotation.Values) == 0 {
		return BaseURL{}, errors.New("URL not found")
	}

	return BaseURL{Value: annotation.Values[0]}, nil
}

type Header struct {
	Key    string
	Values []Variable
}

func (h Header) name() string {
	return "Header"
}

func (h Header) from(annotation Annotation) (any, error) {
	if len(annotation.Values) == 0 {
		return Header{}, errors.New("key not found")
	}
	if len(annotation.Values) == 1 {
		return Header{}, errors.New("value not found")
	}

	var values []Variable
	for _, value := range annotation.Values[1:] {
		values = append(values, Variable(value))
	}

	return Header{Key: annotation.Values[0], Values: values}, nil
}

type Query struct {
	Key   string
	Value Variable
}

func (q Query) name() string {
	return "Query"
}

func (q Query) from(annotation Annotation) (any, error) {
	if len(annotation.Values) == 0 {
		return Query{}, errors.New("key not found")
	}
	if len(annotation.Values) == 1 {
		return Query{}, errors.New("value not found")
	}
	if len(annotation.Values) > 2 {
		return Query{}, errors.New("too many values")
	}

	return Query{Key: annotation.Values[0], Value: Variable(annotation.Values[1])}, nil
}

type RequestLine struct {
	Method string
	Path   string

	params []Variable
}

func (r RequestLine) name() string {
	return "RequestLine"
}

func (r RequestLine) PathParams() []Variable {
	if r.params == nil {
		r.params = []Variable{}
		for _, part := range strings.Split(r.Path, "/") {
			v := Variable(part)
			if v.IsVariable() {
				r.params = append(r.params, v)
			}
		}
	}
	return r.params
}

func (r RequestLine) from(annotation Annotation) (any, error) {
	aliasMatch := !strings.EqualFold(annotation.Name, r.name())
	if aliasMatch {
		if len(annotation.Values) == 0 {
			return RequestLine{}, errors.New("path not found")
		}
		if len(annotation.Values) > 1 {
			return RequestLine{}, errors.New("too many values")
		}
		return RequestLine{Method: annotation.Name, Path: annotation.Values[0]}, nil
	}

	if len(annotation.Values) == 0 {
		return RequestLine{}, errors.New("method not found")
	}
	if len(annotation.Values) == 1 {
		return RequestLine{}, errors.New("path not found")
	}
	return RequestLine{Method: annotation.Values[0], Path: annotation.Values[1]}, nil
}

func (r RequestLine) alias() []string {
	return []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
}

type annotatable interface {
	name() string
	from(annotation Annotation) (any, error)
}

type annotationAlias interface {
	alias() []string
}

var annotatableType = typeFor[annotatable]()

var errSkipTyped = errors.New("skip typed")

func typed(annotation Annotation, input any) error {
	t := reflect.TypeOf(input)

	if t.Kind() != reflect.Ptr {
		return errors.New("input must be a pointer")
	}

	v := reflect.ValueOf(input)
	ptrToType := t.Elem()
	ptrToValue := v.Elem()

	if ptrToType.Kind() == reflect.Slice {
		itemValue := reflect.New(ptrToType.Elem())
		if err := typed(annotation, itemValue.Interface()); err != nil {
			return err
		}
		ptrToValue.Set(reflect.Append(ptrToValue, itemValue.Elem()))
		return nil
	}

	// check if the dereferenced type implements annotatable
	if ptrToType.Implements(annotatableType) {
		anno := ptrToValue.Interface().(annotatable)

		aliasMatch := func() bool {
			if aliases, ok := anno.(annotationAlias); ok {
				for _, alias := range aliases.alias() {
					if strings.EqualFold(annotation.Name, alias) {
						return true
					}
				}
			}
			return false
		}

		if strings.EqualFold(annotation.Name, anno.name()) || aliasMatch() {
			newValue, err := anno.from(annotation)
			if err != nil {
				return err
			}
			ptrToValue.Set(reflect.ValueOf(newValue))
			return nil
		}
	}

	return errSkipTyped
}

// typeFor returns the [Type] that represents the type argument T.
func typeFor[T any]() reflect.Type {
	var v T
	if t := reflect.TypeOf(v); t != nil {
		return t // optimize for T being a non-interface kind
	}
	return reflect.TypeOf((*T)(nil)).Elem() // only for an interface kind
}
