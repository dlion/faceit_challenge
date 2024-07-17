package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUserFilter(t *testing.T) {
	t.Run("Generate a BSON from a filter", func(t *testing.T) {
		filterBuilder := NewFilterBuilder()
		country := "UK"
		email := "test@test.com"
		userFilter := filterBuilder.
			ByCountry(&country).
			ByEmail(&email).
			Build()

		assert.Equal(t, bson.M{"country": "UK", "email": "test@test.com"}, userFilter.ToBSON())
	})
}
