package registry

import (
	"testing"

	"github.com/lexcao/genapi/internal/runtime/registry/testdata/multiple/pkg1"
	pkg1_value "github.com/lexcao/genapi/internal/runtime/registry/testdata/multiple/pkg1/value"
	"github.com/lexcao/genapi/internal/runtime/registry/testdata/multiple/pkg2"
	pkg2_value "github.com/lexcao/genapi/internal/runtime/registry/testdata/multiple/pkg2/value"
	"github.com/lexcao/genapi/internal/runtime/registry/testdata/simple"
	"github.com/stretchr/testify/assert"
)

type Doer interface {
	Do() string
}

type implDoer struct {
}

func (i *implDoer) Do() string {
	return "done"
}

func TestRegistry(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		t.Run("Local", func(t *testing.T) {
			Register[Doer, implDoer]()
			got := New[Doer]()
			assert.Equal(t, "done", got.Do())
		})

		t.Run("Pointer", func(t *testing.T) {
			Register[Doer, *implDoer]()
			got := New[Doer]()
			assert.Equal(t, "done", got.Do())
		})

		t.Run("Package", func(t *testing.T) {
			Register[simple.Valuer, *simple.ImplValuer]()
			got := New[simple.Valuer]()
			assert.Equal(t, "Value from simple", got.Value())
		})
	})

	t.Run("Multiple", func(t *testing.T) {
		t.Run("SameName", func(t *testing.T) {
			Register[pkg1.Valuer, *pkg1.ImplValuer]()
			Register[pkg2.Valuer, *pkg2.ImplValuer]()

			got := New[pkg1.Valuer]()
			assert.Equal(t, "Value from pkg1", got.Value())

			got = New[pkg2.Valuer]()
			assert.Equal(t, "Value from pkg2", got.Value())
		})

		t.Run("SamePath", func(t *testing.T) {
			Register[pkg1_value.Valuer, *pkg1_value.ImplValuer]()
			Register[pkg2_value.Valuer, *pkg2_value.ImplValuer]()

			got := New[pkg1_value.Valuer]()
			assert.Equal(t, "Value from pkg1.value", got.Value())

			got = New[pkg2_value.Valuer]()
			assert.Equal(t, "Value from pkg2.value", got.Value())
		})
	})

	t.Run("Panic", func(t *testing.T) {
		type NotImpl interface {
			NotImplemented() string
		}
		assert.Panics(t, func() {
			New[NotImpl]()
		})
	})
}
