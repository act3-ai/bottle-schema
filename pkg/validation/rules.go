package validation

import (
	"context"
	"errors"
	"mime"
	"path"
	"strings"

	"git.act3-ace.com/ace/data/schema/pkg/mediatype"
	"git.act3-ace.com/ace/data/schema/pkg/util"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/opencontainers/go-digest"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// checkIsMediaType parses the media type just to check for errors in formatting
// Any media type is valid
func checkIsMediaType(value any) error {
	_, _, err := mime.ParseMediaType(value.(string))
	return err
}

// IsMediaType checkes to make sure it in the correct format of a media type
var IsMediaType = validation.By(checkIsMediaType)

func checkIsRelativePath(value any) error {
	if path.IsAbs(value.(string)) {
		return errors.New("path must be relative")
	}
	return nil
}

// IsRelativePath makes sure it is a relative path
var IsRelativePath = validation.By(checkIsRelativePath)

func checkIsDigest(value any) error {
	return value.(digest.Digest).Validate()
}

// IsDigest makes sure the digest.Digest is a validly formatted digest
var IsDigest = validation.By(checkIsDigest)

func checkIsPortablePath(value any) error {
	if !util.IsPortablePath(value.(string)) {
		return errors.New("path contains invalid (non-portable) characters")
	}
	return nil
}

// IsPortablePath makes sure it is a portable path
var IsPortablePath = validation.By(checkIsPortablePath)

func checkTrailingSlash(value any) error {
	if !strings.HasSuffix(value.(string), "/") {
		return errors.New("must have a trailing slash")
	}
	return nil
}

// TrailingSlash makes sure the has a trailing slash
var TrailingSlash = validation.By(checkTrailingSlash)

func checkNoTrailingSlash(value any) error {
	if strings.HasSuffix(value.(string), "/") {
		return errors.New("must not have a trailing slash")
	}
	return nil
}

// NoTrailingSlash makes sure the has no trailing slash
var NoTrailingSlash = validation.By(checkNoTrailingSlash)

type layerIndexKey struct{}

func checkTrailingSlashWithContext(ctx context.Context, value any) error {
	if manifest := ManifestFromContext(ctx); manifest != nil {
		// find our layer
		index := ctx.Value(layerIndexKey{})
		i := index.(*int)
		desc := manifest.Layers[*i]

		if mediatype.IsArchived(desc.MediaType) {
			// require a trailing slash
			return checkTrailingSlash(value)
		}
		// make sure one it not there
		return checkNoTrailingSlash(value)
	}
	return nil
}

// TrailingSlashInPart ensure the trailing slash is there iff the part is a archive part based on the manifest in the context
var TrailingSlashInPart = validation.WithContext(checkTrailingSlashWithContext)

// KubernetesLabels ensures that the labels keys and value conform to the Kubernetes rules for labels
var KubernetesLabels = validation.By(func(value any) error {
	return v1validation.ValidateLabels(value.(map[string]string), field.NewPath("")).ToAggregate()
})

// KubernetesAnnotations ensures that the keys and values conform to the Kubernetes rules for annotations
var KubernetesAnnotations = validation.By(func(value any) error {
	return apivalidation.ValidateAnnotations(value.(map[string]string), field.NewPath("")).ToAggregate()
})
