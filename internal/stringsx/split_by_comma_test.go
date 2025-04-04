package stringsx_test

import (
	"reflect"
	"testing"

	"github.com/jbduncan/go-containers/internal/stringsx"
)

func TestSplitByComma(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty string",
			args: args{
				s: "",
			},
			want: []string{},
		},
		{
			name: "one-element string",
			args: args{
				s: "1",
			},
			want: []string{"1"},
		},
		{
			name: "two-element string",
			args: args{
				s: "1, 2",
			},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := stringsx.SplitByComma(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitByComma(): got %v, want %v", got, tt.want)
			}
		})
	}
}
