package wasmfile

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestExpression(t *testing.T) {
	numbers := []int64{
		100, 1000, 1000000, -100, 27381092,
	}

	b := make([]byte, 0)
	for _, n := range numbers {
		b = AppendSleb128(b, n)
	}

	// Now decode

	for _, expected := range numbers {
		n, l := DecodeSleb128(b)
		assert.Equal(t, n, expected)
		b = b[l:]
	}

	/*
		require.NoError(t, err)

		assert.Equal()
	*/
}
