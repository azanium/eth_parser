package utils

import (
	"testing"
)

func TestHexToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{
			name:    "valid hex with 0x prefix",
			input:   "0x1a",
			want:    26,
			wantErr: false,
		},
		{
			name:    "valid hex without prefix",
			input:   "1a",
			want:    26,
			wantErr: false,
		},
		{
			name:    "zero value",
			input:   "0x0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "large hex number",
			input:   "0xff",
			want:    255,
			wantErr: false,
		},
		{
			name:    "invalid hex string",
			input:   "0xzz",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressToHex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "validate address",
			input: "0x123456789abcdef123456789abcdef123456789a",
			want:  "0x000000000000000000000000123456789abcdef123456789abcdef123456789a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddressToHex(tt.input); got != tt.want {
				t.Errorf("AddressToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
