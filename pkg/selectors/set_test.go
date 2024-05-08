package selectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestParseSelectors(t *testing.T) {
	makeRequirement := func(key string, op selection.Operator, vals []string, opts ...field.PathOption) labels.Requirement {
		req, err := labels.NewRequirement(key, op, vals, opts...)
		require.NoError(t, err)
		return *req
	}

	single := labels.NewSelector()
	single = single.Add(makeRequirement("a", selection.Equals, []string{"b"}))

	multi1 := labels.NewSelector()
	multi1 = multi1.Add(
		makeRequirement("a", selection.Equals, []string{"b"}),
		makeRequirement("c", selection.NotIn, []string{"back", "front"}),
	)
	multi2 := labels.NewSelector()
	multi2 = multi2.Add(makeRequirement("x", selection.DoesNotExist, []string{}))

	type args struct {
		selectors []string
	}
	tests := []struct {
		name    string
		args    args
		want    LabelSelectorSet
		wantErr error
	}{
		{
			"everything",
			args{nil},
			Everything(),
			nil,
		},
		{
			"empty",
			args{[]string{}},
			Nothing(),
			nil,
		},
		{
			"single",
			args{[]string{"a=b"}},
			LabelSelectorSet{single},
			nil,
		},
		{
			"multi",
			args{[]string{"a=b,c notin (front, back)", "!x"}},
			LabelSelectorSet{multi1, multi2},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.selectors)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
				return
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLabelSelectorSet_Matches(t *testing.T) {
	selectorSet, err := Parse([]string{"a=b,c notin (front, back)", "x"})
	require.NoError(t, err)

	type args struct {
		l labels.Set
	}
	tests := []struct {
		name string
		ls   *LabelSelectorSet
		args args
		want bool
	}{
		{
			"match",
			&selectorSet,
			args{labels.Set(map[string]string{"z": "99", "a": "b"})},
			true,
		},
		{
			"nomatch",
			&selectorSet,
			args{labels.Set(map[string]string{"z": "99"})},
			false,
		},
		{
			"notin",
			&selectorSet,
			args{labels.Set(map[string]string{"a": "b", "c": "other"})},
			true,
		},
		{
			name: "empty",
			ls:   &selectorSet,
			args: args{labels.Set(map[string]string{})},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ls.Matches(tt.args.l)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLabelSelectorSet_MatchAny(t *testing.T) {
	selectorSet, err := Parse([]string{"a=b,c notin (front, back)", "x"})
	require.NoError(t, err)

	type args struct {
		l []labels.Labels
	}
	tests := []struct {
		name string
		lss  *LabelSelectorSet
		args args
		want bool
	}{
		{
			"match",
			&selectorSet,
			args{[]labels.Labels{
				labels.Set(map[string]string{"z": "99", "a": "b"}),
			}},
			true,
		},
		{
			"nomatch",
			&selectorSet,
			args{[]labels.Labels{
				labels.Set(map[string]string{"z": "99"}),
			}},
			false,
		},
		{
			"notin",
			&selectorSet,
			args{[]labels.Labels{
				labels.Set(map[string]string{"a": "b", "c": "other"}),
				labels.Set(map[string]string{}),
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.lss.MatchAny(tt.args.l))

			// twice
			assert.Equal(t, tt.want, tt.lss.MatchAny(append(tt.args.l, tt.args.l...)))
		})
	}
}

func TestLabelSelectorSet_MatchAll(t *testing.T) {
	selectorSet, err := Parse([]string{"a=b,c notin (front, back)", "x"})
	require.NoError(t, err)

	type args struct {
		l []labels.Labels
	}
	tests := []struct {
		name string
		lss  *LabelSelectorSet
		args args
		want bool
	}{
		{
			"match",
			&selectorSet,
			args{[]labels.Labels{
				labels.Set(map[string]string{"z": "99", "a": "b"}),
			}},
			true,
		},
		{
			"nomatch",
			&selectorSet,
			args{[]labels.Labels{
				labels.Set(map[string]string{"z": "99"}),
			}},
			false,
		},
		{
			"notin",
			&selectorSet,
			args{[]labels.Labels{
				labels.Set(map[string]string{"a": "b", "c": "other"}),
			}},
			true,
		},
		{
			name: "nope",
			lss:  &selectorSet,
			args: args{[]labels.Labels{
				labels.Set(map[string]string{"z": "99", "a": "b"}), // matches
				labels.Set(map[string]string{}),                    // does not match
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.lss.MatchAll(tt.args.l))

			// twice
			assert.Equal(t, tt.want, tt.lss.MatchAll(append(tt.args.l, tt.args.l...)))
		})
	}
}

func TestLabelsFromSets(t *testing.T) {
	tests := []struct {
		name string
		s    []labels.Set
		want []labels.Labels
	}{
		{
			name: "basic",
			s: []labels.Set{
				labels.Set(map[string]string{"z": "99", "a": "b"}),
				labels.Set(map[string]string{"1": "2"}),
			},
			want: []labels.Labels{
				labels.Set(map[string]string{"z": "99", "a": "b"}),
				labels.Set(map[string]string{"1": "2"}),
			},
		},
		{
			name: "nil",
			s:    nil,
			want: nil,
		},
		{
			name: "empty",
			s:    []labels.Set{},
			want: []labels.Labels{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LabelsFromSets(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
