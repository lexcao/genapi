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

type Config struct {
	Hello string
}

func TestRegistry(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		t.Run("Local", func(t *testing.T) {
			Register[Doer, implDoer]()
			got, _, err := New[Doer]()
			assert.NoError(t, err)
			assert.Equal(t, "done", got.Do())
		})

		t.Run("Config", func(t *testing.T) {
			Register[Doer, implDoer](Config{Hello: "world"})
			got, config, err := New[Doer]()
			assert.NoError(t, err)
			assert.Equal(t, "done", got.Do())
			assert.Equal(t, Config{Hello: "world"}, config)
		})

		t.Run("Pointer", func(t *testing.T) {
			Register[Doer, *implDoer]()
			got, _, err := New[Doer]()
			assert.NoError(t, err)
			assert.Equal(t, "done", got.Do())
		})

		t.Run("Package", func(t *testing.T) {
			Register[simple.Valuer, *simple.ImplValuer]()
			got, _, err := New[simple.Valuer]()
			assert.NoError(t, err)
			assert.Equal(t, "Value from simple", got.Value())
		})
	})

	t.Run("Multiple", func(t *testing.T) {
		t.Run("SameName", func(t *testing.T) {
			Register[pkg1.Valuer, *pkg1.ImplValuer]()
			Register[pkg2.Valuer, *pkg2.ImplValuer]()

			got, _, err := New[pkg1.Valuer]()
			assert.NoError(t, err)
			assert.Equal(t, "Value from pkg1", got.Value())

			got, _, err = New[pkg2.Valuer]()
			assert.NoError(t, err)
			assert.Equal(t, "Value from pkg2", got.Value())
		})

		t.Run("SamePath", func(t *testing.T) {
			Register[pkg1_value.Valuer, *pkg1_value.ImplValuer]()
			Register[pkg2_value.Valuer, *pkg2_value.ImplValuer]()

			got, _, err := New[pkg1_value.Valuer]()
			assert.NoError(t, err)
			assert.Equal(t, "Value from pkg1.value", got.Value())

			got, _, err = New[pkg2_value.Valuer]()
			assert.NoError(t, err)
			assert.Equal(t, "Value from pkg2.value", got.Value())
		})
	})

	t.Run("ErrorOnMissingRegistration", func(t *testing.T) {
		type NotImpl interface {
			NotImplemented() string
		}
		_, _, err := New[NotImpl]()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no registration found for interface")
	})
}
