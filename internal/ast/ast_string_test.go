// internal/ast/ast_string_test.go
package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpToString(t *testing.T) {
	tests := []struct {
		name string
		exp  Exp
		want string
	}{
		{
			name: "ImmExp with NumberFactor",
			exp:  &ImmExp{Factor: &NumberFactor{Value: 123}},
			want: "123",
		},
		{
			name: "ImmExp with HexFactor",
			exp:  &ImmExp{Factor: &HexFactor{Value: "0x42"}},
			want: "0x42",
		},
		{
			name: "ImmExp with IdentFactor",
			exp:  &ImmExp{Factor: &IdentFactor{Value: "ESP"}},
			want: "ESP",
		},
		{
			name: "MemoryAddrExp with AddExp",
			exp: &MemoryAddrExp{
				Left: &AddExp{
					HeadExp: &MultExp{
						HeadExp: &ImmExp{Factor: &IdentFactor{Value: "ESP"}},
					},
					Operators: []string{"+"},
					TailExps: []*MultExp{
						&MultExp{
							HeadExp: &ImmExp{Factor: &NumberFactor{Value: 4}},
						},
					},
				},
			},
			want: "[ESP + 4]",
		},
		{
			name: "MemoryAddrExp with SegmentExp",
			exp: &MemoryAddrExp{
				Left: &AddExp{
					HeadExp: &MultExp{
						HeadExp: &ImmExp{
							Factor: &IdentFactor{Value: "CS"},
						},
					},
				},
				Right: &AddExp{
					HeadExp: &MultExp{
						HeadExp: &ImmExp{Factor: &HexFactor{Value: "0x20"}},
					},
				},
			},
			want: "[CS:0x20]",
		},
		{
			name: "AddExp",
			exp: &AddExp{
				HeadExp: &MultExp{
					HeadExp: &ImmExp{Factor: &IdentFactor{Value: "ESP"}},
				},
				Operators: []string{"+"},
				TailExps: []*MultExp{
					&MultExp{
						HeadExp: &ImmExp{Factor: &NumberFactor{Value: 4}},
					},
				},
			},
			want: "ESP + 4",
		},
		{
			name: "MultExp",
			exp: &MultExp{
				HeadExp:   &ImmExp{Factor: &NumberFactor{Value: 4}},
				Operators: []string{"*"},
				TailExps: []Exp{
					&ImmExp{Factor: &IdentFactor{Value: "ESI"}},
				},
			},
			want: "4 * ESI",
		},
		{
			name: "MemoryAddrExp with DataType",
			exp: &MemoryAddrExp{
				DataType: Byte,
				Left: &AddExp{
					HeadExp: &MultExp{
						HeadExp: &ImmExp{Factor: &IdentFactor{Value: "EAX"}},
					},
				},
			},
			want: "BYTE [EAX]",
		},
		{
			name: "SegmentExp with DataType",
			exp: &SegmentExp{
				DataType: Dword,
				Left: &AddExp{
					HeadExp: &MultExp{
						HeadExp: &ImmExp{Factor: &IdentFactor{Value: "CS"}},
					},
				},
				Right: &AddExp{
					HeadExp: &MultExp{
						HeadExp: &ImmExp{Factor: &HexFactor{Value: "0x20"}},
					},
				},
			},
			want: "DWORD CS:0x20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ExpToString(tt.exp)
			assert.Equal(t, tt.want, actual)
		})
	}
}
