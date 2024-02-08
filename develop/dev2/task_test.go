package main

import (
	"errors"
	"testing"
)

func TestUnpacking(t *testing.T) {
	tsts := []struct {
		name    string
		arg     string
		want    string
		wantErr error
	}{
		{
			name:    "number after letter",
			arg:     "a4bc2d5e",
			want:    "aaaabccddddde",
			wantErr: nil,
		},
		{
			name:    "string without \\, \\\\ and numbers",
			arg:     "abcd",
			want:    "abcd",
			wantErr: nil,
		},
		{
			name:    "symbol - with number",
			arg:     "-5",
			want:    "-----",
			wantErr: nil,
		},
		{
			name:    "numbers in string with \\",
			arg:     "qwe\\4\\5",
			want:    "qwe45",
			wantErr: nil,
		},
		{
			name:    "number after number with \\",
			arg:     "qwe\\45",
			want:    "qwe44444",
			wantErr: nil,
		},
		{
			name:    "number after \\\\",
			arg:     "qwe\\\\5",
			want:    "qwe\\\\\\\\\\",
			wantErr: nil,
		},
		{
			name:    "just \\",
			arg:     "\\",
			want:    "",
			wantErr: nil,
		},
		{
			name:    "empty string",
			arg:     "",
			want:    "",
			wantErr: nil,
		},
		{
			name:    "just \\\\",
			arg:     "\\\\",
			want:    "\\",
			wantErr: nil,
		},
		{
			name:    "just number after \\",
			arg:     "\\4",
			want:    "4",
			wantErr: nil,
		},
		{
			name:    "incorrect string with first number",
			arg:     "45asdf",
			want:    "",
			wantErr: ErrorIncorrectString,
		},
		{
			name:    "incorrect string with only one number",
			arg:     "5",
			want:    "",
			wantErr: ErrorIncorrectString,
		},
		{
			name:    "incorrect string with double number after not number letter",
			arg:     "a45",
			want:    "",
			wantErr: ErrorIncorrectString,
		},
	}

	for _, tt := range tsts {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackString(tt.arg)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UnpackString() error: %s, want err: %s", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnpackString() result: %s, want result: %s", got, tt.want)
			}
		})
	}
}
