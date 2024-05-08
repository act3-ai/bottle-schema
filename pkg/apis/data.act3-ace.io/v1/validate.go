package v1

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"gitlab.com/act3-ai/asce/data/schema/pkg/mediatype"
	"gitlab.com/act3-ai/asce/data/schema/pkg/util"
	val "gitlab.com/act3-ai/asce/data/schema/pkg/validation"
)

// Validate Part
func (p Part) Validate() error {
	return p.ValidateWithContext(context.Background())
}

// ValidateWithContext Part using ozzo-validation and k8s label validation
// uses the "manifest" if provided in the context
func (p Part) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &p,
		validation.Field(&p.Name, validation.Required, val.IsRelativePath, val.IsPortablePath),
		// zero is a valid value so we cannot use the validation.Required test
		validation.Field(&p.Size, validation.Min(0)),
		validation.Field(&p.Digest, validation.Required, val.IsDigest),
		validation.Field(&p.Labels, val.KubernetesLabels),
	)
}

// Validate Source using ozzo-validation
func (s Source) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.URI, validation.Required /*, is.URL*/),
	)
}

// Validate Author using ozzo-validation
func (a Author) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Email, validation.Required, is.EmailFormat),
		validation.Field(&a.URL /*, is.URL*/),
	)
}

// Validate PublicArtifact using ozzo-validation
func (a PublicArtifact) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.MediaType, validation.Required, val.IsMediaType),
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Path, validation.Required, val.IsRelativePath, val.IsPortablePath),
		validation.Field(&a.Digest, validation.Required, val.IsDigest),
	)
}

// Validate Metric using ozzo-validation
func (m Metric) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Value, validation.Required, is.Float),
	)
}

// validateMetrics ensures that the metrics names are unique
func validateMetrics(metrics []Metric) error {
	// Metrics.Name is unique
	metricNames := make(map[string]bool, len(metrics))
	for _, m := range metrics {
		if _, exists := metricNames[m.Name]; exists {
			return fmt.Errorf("metric name '%s' is not unique", m.Name)
		}
		metricNames[m.Name] = true
	}
	return nil
}

// validatePartsWithContext ensures that the part name is unique and that no part is a prefix of any other part
func validatePartsWithContext(ctx context.Context, parts []Part) error {
	if manifest := val.ManifestFromContext(ctx); manifest != nil {
		nLayers := len(manifest.Layers)
		nParts := len(parts)
		if nLayers != nParts {
			return fmt.Errorf("number of parts (%d) is not equal to the number of layers (%d)", nParts, nLayers)
		}

		for i, p := range parts {
			layer := manifest.Layers[i]
			if mediatype.IsArchived(layer.MediaType) {
				if !strings.HasSuffix(p.Name, "/") {
					return fmt.Errorf("part '%s' (index %d) is an archive thus it must have a trailing slash", p.Name, i)
				}
			} else {
				if strings.HasSuffix(p.Name, "/") {
					return fmt.Errorf("part '%s' (index %d) is not an archive thus it must not have a trailing slash", p.Name, i)
				}
			}
		}
	}

	// Parts.Name is unique
	partNames := make(map[string]struct{}, len(parts))
	for i, p := range parts {
		if _, exists := partNames[p.Name]; exists {
			return fmt.Errorf("part name '%s' is not unique", p.Name)
		}
		partNames[p.Name] = struct{}{}

		// check that no part.Name is a prefix of any other part.Name
		// The following is not allowed:
		// foo
		// foo/bar
		// foo/dog
		for j, otherPart := range parts {
			if i == j {
				continue
			}
			if util.IsPathPrefix(otherPart.Name, p.Name) {
				return fmt.Errorf("part '%s' is invalid because it is a prefix of part '%s'", p.Name, otherPart.Name)
			}
		}
	}
	return nil
}

// TODO might also want to check that the number of layers equals the number of parts

// validatePublicArtifacts validates public artifacts path is unique and that each artifact belongs to a single part
func validatePublicArtifacts(b Bottle) error {
	// PublicArtifacts.Path is unique
	artifactPaths := make(map[string]struct{}, len(b.Parts))
	for _, a := range b.PublicArtifacts {
		if _, exists := artifactPaths[a.Path]; exists {
			return fmt.Errorf("public artifact path '%s' is not unique", a.Path)
		}
		artifactPaths[a.Path] = struct{}{}

		// artifact must be in exactly one part
		enclosingParts := 0
		for _, p := range b.Parts {
			if util.IsPathPrefix(a.Path, p.Name) {
				enclosingParts++
			}
		}
		if enclosingParts == 0 {
			return fmt.Errorf("public artifact path '%s' is not in any part", a.Path)
		}
		if enclosingParts > 1 {
			// This might not be possible given the requirements on parts.Name
			return fmt.Errorf("public artifact path '%s' is in more multiple parts (the parts are specified incorrectly)", a.Path)
		}
	}
	return nil
}

// Validate Bottle using ozzo-validation
func (b Bottle) Validate() error {
	return b.ValidateWithContext(context.Background())
}

// ValidateWithContext Bottle using ozzo-validation
// If a "manifest" is provided in the context that is used for further validation
func (b Bottle) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &b,
		validation.Field(&b.APIVersion, validation.Required, validation.In(GroupVersion.String())),
		validation.Field(&b.Kind, validation.Required, validation.In("Bottle")),
		validation.Field(&b.Labels, val.KubernetesLabels),
		validation.Field(&b.Annotations, val.KubernetesAnnotations),
		validation.Field(&b.Sources),
		validation.Field(&b.Authors),
		validation.Field(&b.Metrics, validation.By(func(value any) error {
			return validateMetrics(value.([]Metric))
		})),
		validation.Field(&b.PublicArtifacts, validation.By(func(value any) error {
			return validatePublicArtifacts(b)
		})),
		validation.Field(&b.Parts, validation.WithContext(func(ctx context.Context, value any) error {
			return validatePartsWithContext(ctx, value.([]Part))
		})),
	)
}
