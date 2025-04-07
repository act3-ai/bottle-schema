package v1

import (
	"fmt"
	"strings"

	"github.com/opencontainers/go-digest"
	ocispecv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"k8s.io/apimachinery/pkg/conversion"

	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha2"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha3"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha4"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha5"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1beta1"
	"github.com/act3-ai/bottle-schema/pkg/mediatype"
)

// SetDefault_Bottle sets the fields not already set to default values
func SetDefault_Bottle(in *Bottle) { //revive:disable-line:var-naming
	in.SetGroupVersionKind(GroupVersion.WithKind("Bottle"))

	// TODO do we want to default the slices and the maps?
}

// Convert_v1beta1_Bottle_To_v1_Bottle converts Bottle from v1beta1 to v1
func Convert_v1beta1_Bottle_To_v1_Bottle(in *v1beta1.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	out.APIVersion = GroupVersion.String()
	out.Kind = "Bottle"

	// No root level migrations
	out.Annotations = in.Annotations
	out.Labels = in.Labels
	out.Description = in.Description

	// migrate sources, URL => URI
	out.Sources = make([]Source, len(in.Sources))
	for i, s := range in.Sources {
		out.Sources[i] = Source(s)
	}

	// migrate authors -> stays the same
	out.Authors = make([]Author, len(in.Authors))
	for i, a := range in.Authors {
		out.Authors[i] = Author(a)
	}

	// migrate metrics -> stays the same
	out.Metrics = make([]Metric, len(in.Metrics))
	for i, m := range in.Metrics {
		out.Metrics[i] = Metric(m)
	}

	// migrate public artifact -> stays te same
	out.PublicArtifacts = make([]PublicArtifact, len(in.PublicArtifacts))
	for i, art := range in.PublicArtifacts {
		out.PublicArtifacts[i] = PublicArtifact(art)
	}

	// migrate deprecates move from annotation to the new field
	if deprecatesAnnotation := in.Annotations[v1beta1.AnnotationDeprecates]; deprecatesAnnotation != "" {
		deprecatedBottleDigests := strings.Split(deprecatesAnnotation, ",")
		out.Deprecates = make([]digest.Digest, len(deprecatedBottleDigests))
		for i, d := range deprecatedBottleDigests {
			dgst, err := digest.Parse(d)
			if err != nil {
				return fmt.Errorf("could not parse deprecated bottle digest from string %s", d)
			}
			out.Deprecates[i] = dgst
		}
		delete(out.Annotations, v1beta1.AnnotationDeprecates)
	}

	manifest, haveManifest := scope.Meta().Context.(*ocispecv1.Manifest)

	// migrate parts -> stays the same
	out.Parts = make([]Part, len(in.Parts))
	for i, f := range in.Parts {
		out.Parts[i] = Part(f)
		if haveManifest {
			// use the manifest to handle the conversion for ensuring directory parts have a trailing slash
			desc := manifest.Layers[i]
			if mediatype.IsArchived(desc.MediaType) {
				// require a trailing slash
				if !strings.HasSuffix(out.Parts[i].Name, "/") {
					// add a trailing slash
					out.Parts[i].Name += "/"
				}
			} else {
				// make sure one it not there
				out.Parts[i].Name = strings.TrimSuffix(out.Parts[i].Name, "/")
			}
		}
	}

	return nil
}

// Convert_v1alpha5_Bottle_To_v1_Bottle converts Bottle from v1alpha5 to v1
func Convert_v1alpha5_Bottle_To_v1_Bottle(in *v1alpha5.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	b1 := &v1beta1.Bottle{}
	if err := scope.Convert(in, b1); err != nil {
		return err
	}
	return scope.Convert(b1, out)
}

// Convert_v1alpha4_Bottle_To_v1_Bottle converts Bottle from v1alpha4 to v1
func Convert_v1alpha4_Bottle_To_v1_Bottle(in *v1alpha4.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a5 := &v1alpha5.Bottle{}
	if err := scope.Convert(in, a5); err != nil {
		return err
	}
	return scope.Convert(a5, out)
}

// Convert_v1alpha3_Bottle_To_v1_Bottle converts Bottle from v1alpha3 to v1
func Convert_v1alpha3_Bottle_To_v1_Bottle(in *v1alpha3.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a4 := &v1alpha4.Bottle{}
	if err := scope.Convert(in, a4); err != nil {
		return err
	}
	return scope.Convert(a4, out)
}

// Convert_v1alpha2_Bottle_To_v1_Bottle converts Bottle from v1alpha2 to v1
func Convert_v1alpha2_Bottle_To_v1_Bottle(in *v1alpha2.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a4 := &v1alpha4.Bottle{}
	if err := scope.Convert(in, a4); err != nil {
		return err
	}
	return scope.Convert(a4, out)
}
