package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input   string
		wantErr bool
	}{
		{
			input:   `{"hoge":"fuga", "piyo": 13}`,
			wantErr: false,
		},
		{
			input:   `0}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		reader := strings.NewReader(tt.input)
		parser, err := NewParser(reader)
		if err != nil {
			t.Fatal(err)
		}
		result, err := parser.Parse()
		if tt.wantErr && err == nil {
			t.Fatalf("want error but no error")
		}
		if err == nil {
			fmt.Printf("input %s, got %+v", tt.input, result)
		}
	}
}
