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
			got, _ := New[Doer]()
			assert.Equal(t, "done", got.Do())
		})

		t.Run("Config", func(t *testing.T) {
			Register[Doer, implDoer](Config{Hello: "world"})
			got, config := New[Doer]()
			assert.Equal(t, "done", got.Do())
			assert.Equal(t, Config{Hello: "world"}, config)
		})

		t.Run("Pointer", func(t *testing.T) {
			Register[Doer, *implDoer]()
			got, _ := New[Doer]()
			assert.Equal(t, "done", got.Do())
		})

		t.Run("Package", func(t *testing.T) {
			Register[simple.Valuer, *simple.ImplValuer]()
			got, _ := New[simple.Valuer]()
			assert.Equal(t, "Value from simple", got.Value())
		})
	})

	t.Run("Multiple", func(t *testing.T) {
		t.Run("SameName", func(t *testing.T) {
			Register[pkg1.Valuer, *pkg1.ImplValuer]()
			Register[pkg2.Valuer, *pkg2.ImplValuer]()

			got, _ := New[pkg1.Valuer]()
			assert.Equal(t, "Value from pkg1", got.Value())

			got, _ = New[pkg2.Valuer]()
			assert.Equal(t, "Value from pkg2", got.Value())
		})

		t.Run("SamePath", func(t *testing.T) {
			Register[pkg1_value.Valuer, *pkg1_value.ImplValuer]()
			Register[pkg2_value.Valuer, *pkg2_value.ImplValuer]()

			got, _ := New[pkg1_value.Valuer]()
			assert.Equal(t, "Value from pkg1.value", got.Value())

			got, _ = New[pkg2_value.Valuer]()
			assert.Equal(t, "Value from pkg2.value", got.Value())
		})
	})

	t.Run("PanicsWithHelpfulMessage", func(t *testing.T) {
		type NotImpl interface {
			NotImplemented() string
		}
		assert.PanicsWithValue(t, `
genapi: no registration found for interface github.com/lexcao/genapi/internal/runtime/registry.NotImpl

This usually means:
1. You forgot to run 'go generate' on your API package
2. The generated *.gen.go file wasn't imported
3. There's a bug in code generation

Run: go generate ./...
`, func() {
			New[NotImpl]()
		})
	})
}
