package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	val, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069", val)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	// Test: Invalid field name
	headers = NewHeaders()
	data = []byte("        H©st: localhost:42069 \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// headers = NewHeaders()
	// data = []byte("Host: localhost:42069\r\n Set-Person: lane-loves-go\r\n Set-Person: prime-loves-zig \r\n\r\n")
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)
	// val, _ = headers.Get("Host")
	// assert.Equal(t, "localhost:42069,Set-Person: lane-loves-go,Set-Person: prime-loves-zig", val)
	// // assert.Equal(t, 23, n)
	// assert.False(t, done)
}
