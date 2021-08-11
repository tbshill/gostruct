package main

import (
	"testing"
)

func TestToPublicName(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"HelloWorld", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"hello World", "HelloWorld"},
		{"hello World1", "HelloWorld1"},
		//{"1hello World", "HelloWorld1"}, // TODO(tbshill) fix this
		{"Hello World (yyyymmdd)", "HelloWorld"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			if output := toPublicName(test.input); output != test.expect {
				t.Errorf("expected: %s, got %s\n", test.expect, output)
			}
		})
	}
}
