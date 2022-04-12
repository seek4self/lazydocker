package docker

import "testing"

func TestStats(t *testing.T) {
	type args struct {
		container string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"1", args{"stoic_driscoll"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Stats(tt.args.container)
		})
	}
}
