package v1

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"

	"gitlab.com/act3-ai/asce/data/schema/pkg/mediatype"
	val "gitlab.com/act3-ai/asce/data/schema/pkg/validation"
)

/*
func TestPart_ValidateWithContext(t *testing.T) {
	dgst := digest.Digest("sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0")
	manifest := &ocispec.Manifest{
		Layers: []ocispec.Descriptor{ // This is the wrong number of layers but it allows us to test everything
			{
				MediaType: mediatype.MediaTypeLayerTarGzip,
			},
			{
				MediaType: mediatype.MediaTypeLayerTarGzip,
			},
			{
				MediaType: mediatype.MediaTypeLayerZstd,
			},
			{
				MediaType: mediatype.MediaTypeLayer,
			},
		},
	}
	ctx := context.Background()
	ctxManifest := val.ContextWithManifest(ctx, manifest)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		part    Part
		args    args
		wantErr error
	}{
		{"valid", Part{"foo/bar", 13, dgst, nil}, args{ctx}, nil},
		{"bad digest", Part{"foo/bar", 13, "sha256:deedbeef", nil}, args{ctx}, errors.New("digest: invalid checksum digest length.")}, //nolint
		{"manifest dir", Part{"dogs/", 13, dgst, nil}, args{ctxManifest}, nil},
		{"manifest dir no found", Part{"other", 13, dgst, nil}, args{ctxManifest}, errors.New("name: part \"other\" not in the manifest.")},           //nolint
		{"manifest dir missing trailing slash", Part{"a/archive", 13, dgst, nil}, args{ctxManifest}, errors.New("name: must have a trailing slash.")}, //nolint
		{"manifest file", Part{"cat-file", 13, dgst, nil}, args{ctxManifest}, nil},
		{"manifest file with trailing slash", Part{"a/file/", 13, dgst, nil}, args{ctxManifest}, errors.New("name: must not have a trailing slash.")}, //nolint
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.part
			err := p.ValidateWithContext(tt.args.ctx)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
*/

func TestBottle_Validate(t *testing.T) {
	assert := assert.New(t)
	bottle := testBottle()

	bottle.Labels = map[string]string{
		"key with space": "werd?",
	}
	bottle.Metrics = []Metric{
		{
			Name:  "Here is a valid metric",
			Value: "-0.34",
		},
		{
			Name:  "Not a flow metric",
			Value: "dog",
		},
	}

	err := bottle.Validate()
	assert.Error(err)
	assert.Equal(err.Error(), "labels: [[]: Invalid value: \"key with space\": name part must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyName',  or 'my.name',  or '123-abc', regex used for validation is '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]'), []: Invalid value: \"werd?\": a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')]; metrics: (1: (value: must be a floating point number.).).")
}

func TestBottle_ValidateWithContext(t *testing.T) {
	dgst1 := digest.Digest("sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0")
	dgst2 := digest.Digest("sha256:8dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0")
	manifest := &ocispec.Manifest{
		Layers: []ocispec.Descriptor{ // This is the wrong number of layers but it allows us to test everything
			{
				MediaType: mediatype.MediaTypeLayerTarGzip,
			},
			{
				MediaType: mediatype.MediaTypeLayer,
			},
		},
	}
	ctx := context.Background()
	ctxManifest := val.ContextWithManifest(ctx, manifest)

	validBottle := NewBottle()
	validBottle.Parts = []Part{
		{"dogs/", 13, dgst1, nil},
		{"cats", 20, dgst2, nil},
	}

	invalidBottle := NewBottle()
	invalidBottle.Parts = []Part{
		{"dogs", 13, dgst1, nil},
		{"cats", 20, dgst2, nil},
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		bottle  Bottle
		args    args
		wantErr error
	}{
		{"valid with manifest", validBottle, args{ctxManifest}, nil}, //nolint
		{"invalid part name with manifest", invalidBottle, args{ctxManifest}, errors.New("parts: part 'dogs' (index 0) is an archive thus it must have a trailing slash.")}, //nolint
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.bottle
			err := b.ValidateWithContext(tt.args.ctx)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
