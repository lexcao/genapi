package annotation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Annotation
		wantErr error
	}{
		{
			input:   "",
			want:    Annotation{},
			wantErr: ErrNotFound,
		},
		{
			input:   "   ",
			want:    Annotation{},
			wantErr: ErrNotFound,
		},
		{
			input:   "NotAnnotation",
			want:    Annotation{},
			wantErr: ErrNotFound,
		},
		{
			input:   "@",
			want:    Annotation{},
			wantErr: ErrInvalidFormat{Message: "empty annotation name", Source: ""},
		},
		{
			input:   "@()",
			want:    Annotation{},
			wantErr: ErrInvalidFormat{Message: "empty annotation name", Source: "@()"},
		},
		{
			input: `@Test()`,
			want:  Annotation{Name: "Test"},
		},
		{
			input: `@Test("Value1", "Value2")`,
			want:  Annotation{Name: "Test", Values: []string{"Value1", "Value2"}},
		},
		{
			input: `@Test("Value with spaces")`,
			want:  Annotation{Name: "Test", Values: []string{"Value with spaces"}},
		},
		{
			input: `@Test("Value1","Value2")`,
			want:  Annotation{Name: "Test", Values: []string{"Value1", "Value2"}},
		},
		{
			input: `@Test("Value1",    "Value2"   )`,
			want:  Annotation{Name: "Test", Values: []string{"Value1", "Value2"}},
		},
		{
			input: `  @Test("Value1")`,
			want:  Annotation{Name: "Test", Values: []string{"Value1"}},
		},
		{
			input:   `@Test("unclosed`,
			want:    Annotation{},
			wantErr: ErrInvalidFormat{Message: "unclosed quote in value", Source: `@Test("unclosed`},
		},
		{
			input:   `@Test(Value1)`,
			want:    Annotation{},
			wantErr: ErrInvalidFormat{Message: "value must be quoted", Source: "Value1"},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := parse(test.input)
			if test.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, test.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.want, got)
			}
		})
	}
}
