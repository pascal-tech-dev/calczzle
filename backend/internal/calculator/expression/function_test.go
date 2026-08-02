package expression

import "testing"

func TestIsSupportedFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   string
		want bool
	}{
		{name: "sqrt supported", fn: "sqrt", want: true},
		{name: "uppercase not accepted here", fn: "SQRT", want: false},
		{name: "unknown function", fn: "foo", want: false},
		{name: "empty name", fn: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsSupportedFunction(tt.fn); got != tt.want {
				t.Fatalf("IsSupportedFunction(%q) = %v, want %v", tt.fn, got, tt.want)
			}
		})
	}
}
