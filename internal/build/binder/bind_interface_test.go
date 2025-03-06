package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
	"github.com/stretchr/testify/require"
)

func TestBindInterface(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		iface := &model.Interface{
			Annotations: annotation.InterfaceAnnotations{
				BaseURL: annotation.BaseURL{
					Value: "https://api.example.com",
				},
				Headers: []annotation.Header{
					{
						Key:    "Content-Type",
						Values: []annotation.Variable{"application/json"},
					},
				},
			},
		}

		err := BindInterface(iface)
		require.NoError(t, err)
		require.Equal(t, `genapi.Config{
	BaseURL: "https://api.example.com",
	Header: http.Header{
		"Content-Type": []string{
			"application/json",
		},
	},
}`, iface.Bindings.Config)
		require.ElementsMatch(t, []string{`"net/http"`}, iface.Imports.Slices())
	})

	t.Run("Empty", func(t *testing.T) {
		iface := &model.Interface{}
		err := BindInterface(iface)
		require.NoError(t, err)
		require.Equal(t, `genapi.Config{}`, iface.Bindings.Config)
	})
}
