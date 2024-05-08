package util

import (
	"errors"
	"fmt"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/act3-ai/asce/data/schema/pkg/selectors"
)

func TestParseSourceURI(t *testing.T) {
	type args struct {
		uri string
	}

	sel, err := selectors.Parse([]string{"partkey!=value1,mykey=value2", "partkey2=45"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		args    args
		want    digest.Digest
		want1   selectors.LabelSelectorSet
		wantErr error
	}{
		{"generic", args{"http://www.google.com"}, "", nil, nil},
		{"relative", args{"www.google.com"}, "", nil, errors.New("sources URI must not be relative")},
		{
			"bottle scheme",
			args{"bottle:sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9"},
			digest.Digest("sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9"),
			selectors.Everything(),
			nil,
		},
		{
			"bottle hash",
			args{"hash://sha256/05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?type=application/vnd.act3-ace.bottle.config.v1%2Bjson"},
			digest.Digest("sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9"),
			selectors.Everything(),
			nil,
		},
		{
			"bottle scheme selector",
			args{"bottle:sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?selector=partkey!=value1,mykey=value2&selector=partkey2=45"},
			digest.Digest("sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9"),
			sel,
			nil,
		},
		{
			"bottle hash selector",
			args{"hash://sha256/05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?type=application/vnd.act3-ace.bottle.config.v1%2Bjson&selector=partkey!=value1,mykey=value2&selector=partkey2=45"},
			digest.Digest("sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9"),
			sel,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseSourceURI(tt.args.uri)
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
			assert.Equal(t, tt.want1, got1)
		})
	}
}
