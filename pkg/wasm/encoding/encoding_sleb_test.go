package encoding

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

/**
 * Test our selb128 encoding is ok...
 *
 */
func TestEncodingSleb(t *testing.T) {
	numbers := []int64{
		10, 100, 1000, 1000000, -100, 27381092,
	}

	// Encode the numbers
	b := make([]byte, 0)
	for _, n := range numbers {
		b = AppendSleb128(b, n)
	}

	// Now decode them and assert they're as expected
	for _, expected := range numbers {
		n, l := DecodeSleb128(b)
		assert.Equal(t, n, expected)
		b = b[l:]
	}

	assert.Equal(t, len(b), 0)
}
