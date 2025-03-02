package common

type Set[T comparable] struct {
	items map[T]struct{}
}

func SetOf[T comparable](items ...T) Set[T] {
	var s Set[T]
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s *Set[T]) Add(items ...T) {
	if s.items == nil {
		s.items = make(map[T]struct{})
	}

	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

func (s Set[T]) Contains(item T) bool {
	if len(s.items) == 0 {
		return false
	}

	_, ok := s.items[item]
	return ok
}
func (s Set[T]) Slices() []T {
	if len(s.items) == 0 {
		return []T{}
	}

	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}
