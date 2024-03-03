package gen

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		entryPoint string
		text       string
		want       interface{}
	}{
		{"+int", "Integer", "30", 30},
		{"-int", "Integer", "-30", -30},
		{"letter", "Letter", "a", []rune{'a'}},
		{"hex", "Hex", "0x0ff0", "0x0ff0"},
		{"charseq", "CharSeq", "'0x0ff0'", "0x0ff0"},
		{"ident", "Ident", "_testZ009$", "_testZ009$"},
		{"label", "Label", "_test:\n", "_test:"},
		{"line comment1", "Comment", "# sample \n", ""},
		{"line comment2", "Comment", "; sample \n", ""},
		{"line comment1", "Comment", "# sample", ""},
		{"line comment2", "Comment", "; sample", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse("", []byte(tt.text), Entrypoint(tt.entryPoint))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%v) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
