package v1beta1

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"gitlab.com/act3-ai/asce/data/schema/pkg/util"
	val "gitlab.com/act3-ai/asce/data/schema/pkg/validation"
)

// Validate Part using ozzo-validation and k8s label validation
func (p Part) Validate() error {
	errs := v1validation.ValidateLabels(p.Labels, field.NewPath("labels", p.Name))
	if len(errs) > 0 {
		return errs.ToAggregate()
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, val.IsRelativePath, val.IsPortablePath),
		// zero is a valid value so we cannot use the validation.Required test
		validation.Field(&p.Size, validation.Min(0)),
		validation.Field(&p.Digest, validation.Required, val.IsDigest),
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

// ValidateMetrics ensures that the metrics names are unique
func ValidateMetrics(metrics []Metric) error {
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

// ValidateParts ensures that the part name is unique and that no part is a prefix of any other part
func ValidateParts(parts []Part) error {
	// Parts.Name is unique
	partNames := make(map[string]bool, len(parts))
	for i, p := range parts {
		if _, exists := partNames[p.Name]; exists {
			return fmt.Errorf("part name '%s' is not unique", p.Name)
		}
		partNames[p.Name] = true

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

// ValidatePublicArtifacts validates public artifacts path is unique and that each artifact belongs to a single part
func ValidatePublicArtifacts(b Bottle) error {
	// PublicArtifacts.Path is unique
	artifactPaths := make(map[string]bool, len(b.Parts))
	for _, a := range b.PublicArtifacts {
		if _, exists := artifactPaths[a.Path]; exists {
			return fmt.Errorf("public artifact path '%s' is not unique", a.Path)
		}
		artifactPaths[a.Path] = true

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

// Validate Bottle using ozzo-validation. Returns a list of errors
func (b Bottle) Validate() error {
	var allErrs field.ErrorList
	allErrs = append(allErrs, v1validation.ValidateLabels(b.Labels, field.NewPath("labels"))...)
	allErrs = append(allErrs, apivalidation.ValidateAnnotations(b.Annotations, field.NewPath("annotations"))...)

	if len(allErrs) > 0 {
		return allErrs.ToAggregate()
	}

	if err := ValidateMetrics(b.Metrics); err != nil {
		return err
	}

	if err := ValidateParts(b.Parts); err != nil {
		return err
	}

	if err := ValidatePublicArtifacts(b); err != nil {
		return err
	}

	return validation.ValidateStruct(&b,
		validation.Field(&b.APIVersion, validation.Required, validation.In(GroupVersion.String())),
		validation.Field(&b.Kind, validation.Required, validation.In("Bottle")),
		validation.Field(&b.Sources),
		validation.Field(&b.Authors),
		validation.Field(&b.Metrics),
		validation.Field(&b.PublicArtifacts),
		validation.Field(&b.Parts),
	)
}
