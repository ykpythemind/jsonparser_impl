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
			input:   "{}",
			wantErr: false,
		},
		{
			input:   `}`,
			wantErr: true,
		},
		{
			input:   `0`,
			wantErr: false,
		},
		{
			input:   `"aiueo"`,
			wantErr: false,
		},
		{
			input:   `{"f":{"p": "12"}}`,
			wantErr: false,
		},
		{
			input:   `{"f":{"p": 1200}}`,
			wantErr: false,
		},
		{
			input:   `{"key":-10}`,
			wantErr: false,
		},
		{
			input:   `{"key":false}`,
			wantErr: false,
		},
		{
			input:   `{"key":true}`,
			wantErr: false,
		},
		{
			input:   `{"key":null}`,
			wantErr: false,
		},
		{
			input:   `{"key":["1", 2, {"a": "b"}]}`,
			wantErr: false,
		},
		{
			input:   `[1, 2, {"a": [4, 5, 6]}]`,
			wantErr: false,
		},
		{
			input:   `{"aaaa"::"aaa"}`,
			wantErr: true,
		},
		{
			input:   `{"00000"}`,
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
			t.Errorf("input %s want error but no error", tt.input)
			continue
		}
		if !tt.wantErr && err != nil {
			t.Errorf("input %s got unexpected error: %s", tt.input, err)
			continue
		}

		if !tt.wantErr {
			// if !reflect.DeepEqual(tt.expect, result) {
			// 	t.Logf("expect %+v\n", tt.expect)
			// 	t.Logf("result %+v\n", result)
			// 	t.Error("deep equal fail")
			// }

			// err = json.Unmarshal([]byte(tt.input), &parsedWithGo)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// if !reflect.DeepEqual(parsedWithGo, result) {
			// 	t.Logf("parsedwithgo %+v\n", parsedWithGo)
			// 	t.Logf("result %+v\n", result)
			// 	t.Error("parse is wrong")
			// }

			t.Logf("input %s, got %+v\n", tt.input, result)
		}
	}
}
