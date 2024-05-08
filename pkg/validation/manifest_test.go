package validation

import (
	"fmt"
	"testing"

	"git.act3-ace.com/ace/data/schema/pkg/mediatype"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go"
	ocispecv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestValidateManifest(t *testing.T) {
	dgst1 := digest.Digest("sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0")
	dgst2 := digest.Digest("sha256:8dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0")
	validManifest := ocispecv1.Manifest{
		Versioned: ocispec.Versioned{SchemaVersion: 2},
		MediaType: ocispecv1.MediaTypeImageManifest,
		Config: ocispecv1.Descriptor{
			MediaType: mediatype.MediaTypeBottleConfig,
			Digest:    dgst1,
			Size:      102,
		},
		Layers: []ocispecv1.Descriptor{
			{
				MediaType: mediatype.MediaTypeLayerTar,
				Digest:    dgst2,
				Size:      100,
			},
		},
	}

	type args struct {
		m ocispecv1.Manifest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"valid", args{validManifest}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateManifest(tt.args.m)
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
