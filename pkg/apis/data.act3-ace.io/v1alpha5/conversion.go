package v1alpha5

import (
	"strings"

	"k8s.io/apimachinery/pkg/conversion"

	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/v1alpha2"
	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/v1alpha3"
	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/v1alpha4"
)

// Convert_v1alpha4_Bottle_To_v1alpha5_Bottle converts Bottle from v1alpha4 to v1alpha5
func Convert_v1alpha4_Bottle_To_v1alpha5_Bottle(in *v1alpha4.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	out.APIVersion = GroupVersion.String()
	out.Kind = "Bottle"

	// Migrate: Maintainers -> Authors, Usage -> PublicArtifacts
	// New Fields: Annotations, Labels, Metrics
	out.Annotations = map[string]string{}
	out.Labels = map[string]string{}
	out.Metrics = []Metric{}

	out.Description = in.Description

	out.Sources = make([]Source, len(in.Sources))
	for i, s := range in.Sources {
		out.Sources[i] = Source(s)
	}

	out.Authors = make([]Author, len(in.Maintainers))
	for i, mt := range in.Maintainers {
		out.Authors[i] = Author(mt)
	}

	// Migrate usage info to public artifact.  Digest information missing, need to recalculate in post
	out.PublicArtifacts = make([]PublicArtifact, len(in.Usage))
	for i, u := range in.Usage {
		out.PublicArtifacts[i].Name = u.Name
		out.PublicArtifacts[i].Path = u.File
		out.PublicArtifacts[i].Type = u.Topic
	}

	// If the controller object supports post migrate operations, this migration will use
	//  CalculatePublicArtifactDigest(index int) error   to inform the controller that public artifact information needs to
	//  be recalculated

	// a local interface for accessing artifact digest calculation functionality in the controller
	// type ArtifactUpdateConsumer interface {
	// 	CalculatePublicArtifactDigest(index int) error
	// }

	// capture the local part data that was iniously stored in the config, by offering it to a part update consumer
	// in the controller
	// if pmc, ok := m.Controller.(types.PostMigrateConsumer); ok {
	// 	pmc.RegisterPostMigrateOp("calc-artifact-digest", func(ctl any) error {
	// 		auc, ok := ctl.(ArtifactUpdateConsumer)
	// 		if !ok {
	// 			return nil
	// 		}
	// 		for i := range out.PublicArtifacts {
	// 			// ignore errors here?  This is currently done again during the bottle load process, which offers better
	// 			// logging opportunities
	// 			_ = auc.CalculatePublicArtifactDigest(i)
	// 		}
	// 		return nil
	// 	}, true)
	// }

	// Migrate digest and layer digest digestMap structure to string representation of digest
	out.Parts = make([]Part, len(in.Parts))
	for i, f := range in.Parts {
		part := Part{
			Name:        f.Name,
			Size:        f.Size,
			LayerSize:   f.LayerSize,
			Format:      f.Format,
			Digest:      safeDigest(f.Digest.Sha256),
			LayerDigest: safeDigest(f.LayerDigest.Sha256),
			Modified:    f.Modified,
			Labels:      f.Labels,
		}

		if part.Format == "" {
			part.Format = "raw"
		}
		out.Parts[i] = part
	}

	return nil
}

// safeDigest returns a digest string that can be parsed with digest.Parse, or an empty string if the input is empty.
// In v1alpha4, digests were represented in a digestMap structure that stripped off the algorithm and included only the
// hash. This was later updated to the standard string representation without the structure
func safeDigest(digestStr string) string {
	if digestStr == "" || strings.HasPrefix(digestStr, "sha256:") {
		return digestStr
	}
	return strings.Join([]string{"sha256", digestStr}, ":")
}

// Convert_v1alpha3_Bottle_To_v1alpha5_Bottle converts Bottle from v1alpha3 to v1alpha5
func Convert_v1alpha3_Bottle_To_v1alpha5_Bottle(in *v1alpha3.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a4 := &v1alpha4.Bottle{}
	if err := scope.Convert(in, a4); err != nil {
		return err
	}
	return scope.Convert(a4, out)
}

// Convert_v1alpha2_Bottle_To_v1alpha5_Bottle converts Bottle from v1alpha3 to v1alpha5
func Convert_v1alpha2_Bottle_To_v1alpha5_Bottle(in *v1alpha2.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a3 := &v1alpha3.Bottle{}
	if err := scope.Convert(in, a3); err != nil {
		return err
	}
	return scope.Convert(a3, out)
}
