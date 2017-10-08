package ui

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var printTests = []struct {
	in   []interface{}
	out  string
	desc string
}{
	{
		[]interface{}{"this is a message"},
		"-> this is a message\n",
		"prints single string",
	},
	{
		[]interface{}{"this", "is", "a", "message"},
		"-> this is a message\n",
		"prints multiple strings",
	},
}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	originalOut := Out
	Out = &buf
	defer func() { Out = originalOut }()
	for _, tt := range printTests {
		Print(tt.in...)
		assert.Equal(t, tt.out, buf.String())
		buf.Reset()
	}
}

var printErrorTests = []struct {
	in   []interface{}
	out  string
	desc string
}{
	{
		[]interface{}{"this is an error"},
		"-> this is an error\n",
		"prints single string error",
	},
	{
		[]interface{}{"this", "is", "an", "error"},
		"-> this is an error\n",
		"prints multiple strings error",
	},
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	originalOut := Out
	Out = &buf
	defer func() { Out = originalOut }()
	for _, tt := range printErrorTests {
		Print(tt.in...)
		assert.Equal(t, tt.out, buf.String())
		buf.Reset()
	}
}
