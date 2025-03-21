// internal/ast/ast_string_test.go
package ast_test

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // プロジェクトのパスに合わせて修正
	"github.com/stretchr/testify/assert"
)

func TestExpToString(t *testing.T) {
	tests := []struct {
		name string
		exp  ast.Exp
		want string
	}{
		{
			name: "ImmExp with NumberFactor",
			exp:  &ast.ImmExp{Factor: &ast.NumberFactor{Value: 123}},
			want: "123",
		},
		{
			name: "ImmExp with HexFactor",
			exp:  &ast.ImmExp{Factor: &ast.HexFactor{Value: "0x42"}},
			want: "0x42",
		},
		{
			name: "ImmExp with IdentFactor",
			exp:  &ast.ImmExp{Factor: &ast.IdentFactor{Value: "ESP"}},
			want: "ESP",
		},
		{
			name: "MemoryAddrExp with AddExp",
			exp: &ast.MemoryAddrExp{
				Left: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp: &ast.ImmExp{Factor: &ast.IdentFactor{Value: "ESP"}},
					},
					Operators: []string{"+"},
					TailExps: []*ast.MultExp{
						&ast.MultExp{
							HeadExp: &ast.ImmExp{Factor: &ast.NumberFactor{Value: 4}},
						},
					},
				},
			},
			want: "[ESP + 4]",
		},
		{
			name: "MemoryAddrExp with SegmentExp",
			exp: &ast.MemoryAddrExp{
				Left: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp: &ast.ImmExp{
							Factor: &ast.IdentFactor{Value: "CS"},
						},
					},
				},
				Right: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp: &ast.ImmExp{Factor: &ast.HexFactor{Value: "0x20"}},
					},
				},
			},
			want: "[CS:0x20]",
		},
		{
			name: "AddExp",
			exp: &ast.AddExp{
				HeadExp: &ast.MultExp{
					HeadExp: &ast.ImmExp{Factor: &ast.IdentFactor{Value: "ESP"}},
				},
				Operators: []string{"+"},
				TailExps: []*ast.MultExp{
					&ast.MultExp{
						HeadExp: &ast.ImmExp{Factor: &ast.NumberFactor{Value: 4}},
					},
				},
			},
			want: "ESP + 4",
		},
		{
			name: "MultExp",
			exp: &ast.MultExp{
				HeadExp:   &ast.ImmExp{Factor: &ast.NumberFactor{Value: 4}},
				Operators: []string{"*"},
				TailExps: []*ast.ImmExp{
					&ast.ImmExp{Factor: &ast.IdentFactor{Value: "ESI"}},
				},
			},
			want: "4 * ESI",
		},
		{
			name: "MemoryAddrExp with DataType",
			exp: &ast.MemoryAddrExp{
				DataType: ast.Byte,
				Left: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp: &ast.ImmExp{Factor: &ast.IdentFactor{Value: "EAX"}},
					},
				},
			},
			want: "BYTE [EAX]",
		},
		{
			name: "SegmentExp with DataType",
			exp: &ast.SegmentExp{
				DataType: ast.Dword,
				Left: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp: &ast.ImmExp{Factor: &ast.IdentFactor{Value: "CS"}},
					},
				},
				Right: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp: &ast.ImmExp{Factor: &ast.HexFactor{Value: "0x20"}},
					},
				},
			},
			want: "DWORD CS:0x20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ast.ExpToString(tt.exp)
			assert.Equal(t, tt.want, actual)
		})
	}
}
