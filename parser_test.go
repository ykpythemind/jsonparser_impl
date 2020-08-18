package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input   string
		wantErr bool
	}{
		{
			input:   `{"hoge":"fuga","piyo": "12"}`,
			wantErr: false,
		},
		{
			input:   `}`,
			wantErr: true,
		},
		{
			input:   `{"f":{"p": "12"}}`,
			wantErr: false,
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
			t.Error("want error but no error")
			continue
		}
		if !tt.wantErr && err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		t.Logf("input %s, got %+v\n", tt.input, result)
	}
}
