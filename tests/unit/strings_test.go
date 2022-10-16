package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/util"
)

func TestRemoveDuplicateWhitespace(t *testing.T) {
	t.Run("should remove duplicate whitespace from string", func(t *testing.T) {
		t.Parallel()
		
		s := "Hello              World!"
		assert.Equal(t, "Hello World!", util.RemoveDuplicateWhitespace(s))
	})
}
