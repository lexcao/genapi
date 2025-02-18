package annotation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testAnnotatable struct {
	Value string
}

func (t testAnnotatable) name() string {
	return "test"
}

func (t testAnnotatable) from(annotation Annotation) (any, error) {
	return testAnnotatable{Value: annotation.Values[0]}, nil
}

func TestTyped(t *testing.T) {
	t.Run("Struct", func(t *testing.T) {
		var s testAnnotatable
		err := typed(Annotation{Name: "test", Values: []string{"value"}}, &s)
		require.NoError(t, err)
		require.Equal(t, "value", s.Value)
	})

	t.Run("Slice", func(t *testing.T) {
		var s []testAnnotatable
		err := typed(Annotation{Name: "test", Values: []string{"value1"}}, &s)
		require.NoError(t, err)
		err = typed(Annotation{Name: "test", Values: []string{"value2"}}, &s)
		require.NoError(t, err)
		require.Equal(t, []testAnnotatable{{Value: "value1"}, {Value: "value2"}}, s)
	})
}
