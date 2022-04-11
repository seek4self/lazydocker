package docker

import (
	"reflect"
	"testing"
)

func TestLogs(t *testing.T) {
	type args struct {
		container string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{"1", args{"gitlab-runner-share"}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Logs(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logs() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
