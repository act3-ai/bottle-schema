package util

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/opencontainers/go-digest"

	"gitlab.com/act3-ai/asce/data/schema/pkg/mediatype"
	"gitlab.com/act3-ai/asce/data/schema/pkg/selectors"
)

// IsPathPrefix returns true iff prefixPath is a path prefix for path
func IsPathPrefix(path, prefixPath string) bool {
	// We have to be careful matching because
	// foo.txt (part.Name is foo.txt)
	// foo/bar.txt (part.Name is foo)

	// we really mean "/" here (it does not change on windows to "\")
	return path == prefixPath || strings.HasPrefix(path, prefixPath+"/")
}

// ParseSourceURI will parse a URI for a source.  If it contains a bottle reference it will extract it along with any part selectors.
func ParseSourceURI(uri string) (digest.Digest, selectors.LabelSelectorSet, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", nil, err
	}

	if !u.IsAbs() {
		return "", nil, errors.New("sources URI must not be relative")
	}

	// check if the URL references a bottle, we support a couple types
	var bottleDigest digest.Digest
	var partSelectors selectors.LabelSelectorSet
	switch u.Scheme {
	case "bottle":
		// bottle:sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?selector=partkey!=value1,mykey=value2&selector=partkey2=45
		bottleDigest, err = digest.Parse(u.Opaque)
		if err != nil {
			return "", nil, fmt.Errorf("invalid URI digest in sources for scheme \"bottle\": %w", err)
		}
		partSelectors, err = selectors.Parse(u.Query()["selector"])
		if err != nil {
			return "", nil, fmt.Errorf("parsing selectors: %w", err)
		}
	case "hash":
		// handle the case hash://sha256/05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?type=application/vnd.act3-ace.bottle.config.v1+json&selector=partkey!=value1,mykey=value2&selector=partkey2=45
		qs := u.Query()
		typeValue := qs.Get("type")
		// better test here...
		if typeValue == mediatype.MediaTypeBottleConfig {
			digestStr := u.Host + ":" + strings.TrimPrefix(u.Path, "/")
			// validate the digest format
			bottleDigest, err = digest.Parse(digestStr)
			if err != nil {
				return "", nil, fmt.Errorf("invalid URI digest in sources for scheme \"hash\": %w", err)
			}
			partSelectors, err = selectors.Parse(qs["selector"])
			if err != nil {
				return "", nil, fmt.Errorf("parsing selectors: %w", err)
			}
		} // else it is not a bottle reference so we ignore it
	}

	return bottleDigest, partSelectors, nil
}
