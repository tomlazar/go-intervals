package intervals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkResponse(assert *assert.Assertions, list []Person, err error) {
	assert.NotEqual(list, nil, "The List should not be nil")
	assert.Equal(err, nil, "Their should be no ERROR!")

	for _, item := range list {
		assert.NotEqual(item, nil, "None of the items in the list should be nil")
	}
}

func TestPersonService_List(t *testing.T) {
	assert := assert.New(t)

	c := NewClient(nil)
	list, err := c.PersonService.List(nil)

	checkResponse(assert, list, err)
}

func TestPersonService_ListOne(t *testing.T) {
	assert := assert.New(t)

	c := NewClient(nil)
	list, err := c.PersonService.List(&PersonOptions{Limit: 1})

	checkResponse(assert, list, err)

	assert.Len(list, 1, "There should only be one response")
}
