package validation

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/act3-ai/bottle-schema/pkg/mediatype"
)

// manifestKey is how we find the Manifest in a context.Context.
type manifestKey struct{}

// ManifestFromContext returns the database instance customized for this request
func ManifestFromContext(ctx context.Context) *ocispec.Manifest {
	if v := ctx.Value(manifestKey{}); v != nil {
		return v.(*ocispec.Manifest)
	}
	// panic("manifest missing from context")
	return nil
}

// ContextWithManifest addes the provides manifest to the context (for validation purposes)
func ContextWithManifest(ctx context.Context, manifest *ocispec.Manifest) context.Context {
	return context.WithValue(ctx, manifestKey{}, manifest)
}

// ValidateManifest validates a bottle Manifest for correctness
func ValidateManifest(m ocispec.Manifest) error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SchemaVersion, validation.Required, validation.In(2)),
		validation.Field(&m.MediaType, validation.Required, validation.In(ocispec.MediaTypeImageManifest)),
		validation.Field(&m.Config, configDescriptor),
		// TODO ensure the number of layers is correct
		validation.Field(&m.Layers, validation.Each(layerDescriptor)),
	)
}

// TODO we could also add the bottle to the manifest's context.  Then we can verify that the number of layers in the manifest matches the bottle.  Also that the media types match properly for parts.

func validateConfigDescriptor(config ocispec.Descriptor) error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.MediaType, configMediaType),
		validation.Field(&config.Digest, validation.Required, IsDigest),
		validation.Field(&config.Size, validation.Min(0)),
	)
}

var configDescriptor = validation.By(func(value any) error {
	return validateConfigDescriptor(value.(ocispec.Descriptor))
})

func validateLayerDescriptor(d ocispec.Descriptor) error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.MediaType, validation.Required, layerMediaType),
		validation.Field(&d.Digest, validation.Required, IsDigest),
		validation.Field(&d.Size, validation.Min(0)),
		// validation.Field(&d.URLs, validation.Empty),
		validation.Field(&d.Platform, validation.Empty),
	)
}

var layerDescriptor = validation.By(func(value any) error {
	return validateLayerDescriptor(value.(ocispec.Descriptor))
})

var layerMediaType = validation.By(func(value any) error {
	if !mediatype.IsLayer(value.(string)) {
		return errors.New("invalid layer media type")
	}
	return nil
})

var configMediaType = validation.By(func(value any) error {
	if !mediatype.IsBottleConfig(value.(string)) {
		return errors.New("invalid bottle config media type")
	}
	return nil
})
