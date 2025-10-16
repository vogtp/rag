package usercfg

import (
	"slices"
	"testing"
)

func TestSourceSystem_splitParts(t *testing.T) {
	tests := []struct {
		parts string // description of this test case
		want  []string
	}{
		{parts: "test1", want: []string{"test1"}},
		{parts: "test1,", want: []string{"test1"}},
		{parts: "test1 , ", want: []string{"test1"}},
		{parts: "test1  ", want: []string{"test1"}},
		{parts: "test1  ,,,", want: []string{"test1"}},
		{parts: "test1,test2", want: []string{"test1", "test2"}},
		{parts: "test1 , test2", want: []string{"test1", "test2"}},
		{parts: "test1,,test2", want: []string{"test1", "test2"}},
		{parts: "test1 ,, test2", want: []string{"test1", "test2"}},
		{parts: "test1 , , test2", want: []string{"test1", "test2"}},
	}
	for _, tt := range tests {
		t.Run(tt.parts, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s SourceSystem = SourceSystem{Parts: tt.parts}
			got := s.splitParts()
			if !slices.Equal(got, tt.want) {
				t.Errorf("parts %q splitting: want %#v got %#v ", tt.parts, tt.want, got)
			}
		})
	}
}
