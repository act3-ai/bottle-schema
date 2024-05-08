package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsPortableFilename(t *testing.T) {

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"valid path", "foo.bar", true},
		{"valid special", "foo.-_bar", true},
		{"hidden", ".bar", false},
		{"directory path", "foo/bar", false},
		{"filename with space", "foo bar", false},
		{"empty filename", "", false},
		{"space as filename", " ", false},
		{"dot", ".", false},
		{"dot dot", "..", false},
		{"os path separator", "/", false},
		{"plus sign", "foo+bar", false},
		{"colon character", "foo:bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := IsPortableFilename(tt.path)
			assert.Equalf(t, tt.want, check, "IsPortableFilename(%v)", tt.path)
		})
	}
}

func Test_IsPortablePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"valid path", "dir/file/", true},
		{"valid path no trailing slash", "dir/file", true},
		{"separator only", "/", false},
		{"dot", "./foo/bar", false},
		{"absolute", "/foo/bar", false},
		{"dot dot", "../foo/bar", false},
		{"spaces", "dir/dataset 2020/items/", false},
		{"plus", "dir/foo+bar/items", false},
		{"colon", "dir:foo/bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsPortablePath(tt.path), "IsPortablePath(%v)", tt.path)
		})
	}
}
