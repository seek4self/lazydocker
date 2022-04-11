package docker

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestImages(t *testing.T) {
	tests := []struct {
		name string
		option string
		want []types.ImageSummary
	}{
		// TODO: Add test cases.
		{"1", "",[]types.ImageSummary{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Images(tt.option); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Images() = %+v, want %v", got, tt.want)
			}
		})
	}
}
