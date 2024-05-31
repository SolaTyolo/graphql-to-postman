package introspection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntrospect(t *testing.T) {
	testUrl := "https://swapi-graphql.netlify.app/.netlify/functions/index"
	_, err := Introspect(testUrl)
	assert.Nil(t, err)
}
