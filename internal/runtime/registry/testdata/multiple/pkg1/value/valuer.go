package value

type Valuer interface {
	Value() string
}

type ImplValuer struct{}

func (i *ImplValuer) Value() string {
	return "Value from pkg1.value"
}
