package pkg2

type Valuer interface {
	Value() string
}

type ImplValuer struct{}

func (i *ImplValuer) Value() string {
	return "Value from pkg2"
}
