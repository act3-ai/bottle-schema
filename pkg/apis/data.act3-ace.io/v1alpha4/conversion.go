package v1alpha4

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"

	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/v1alpha2"
	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/v1alpha3"
)

// Convert_v1alpha2_Bottle_To_v1alpha4_Bottle is a bottle converter
func Convert_v1alpha2_Bottle_To_v1alpha4_Bottle(in *v1alpha2.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	// Migrate Files -> Parts
	// New Fields: usage, Expiration

	out.Catalog = in.Catalog
	out.Description = in.Description
	out.Usage = []Usage{}
	out.Keywords = in.Keywords
	out.Expiration = ""

	// migrate sources -> no changes
	out.Sources = make([]Source, len(in.Sources))
	for i, s := range in.Sources {
		out.Sources[i] = Source(s)
	}

	// migrate maintainers -> no changes
	out.Maintainers = make([]Maintainer, len(in.Maintainers))
	for i, mt := range in.Maintainers {
		out.Maintainers[i] = Maintainer(mt)
	}
	// Migrate size -> layersize. digest -> layerdigest.
	// Add fields: size, digest  (uncompressed data, will need recalculating)
	out.Parts = make([]Part, len(in.Files))
	for i, f := range in.Files {
		part := Part{
			Name:        f.Name,
			Size:        0,
			LayerSize:   f.Size,
			Format:      f.Format,
			Digest:      DigestMap{},
			LayerDigest: DigestMap{Sha256: f.Digest.Sha256},
			Modified:    metav1.Time{Time: time.Now()},
			Labels:      f.Labels,
		}
		if part.Format == "" {
			part.Format = "raw"
		}
		out.Parts[i] = part
	}

	// If the controller object supports post migrate operations, this migration will use
	//  RecalcUncompressedSizes() error   to inform the controller that uncompressed digest and size information needs to
	//  be recalculated

	// Register uncompressed size/digest calculator post migrate operation.  This data was not included in the v1alpha2
	// if pmc, ok := m.Controller.(types.PostMigrateConsumer); ok {
	// 	pmc.RegisterPostMigrateOp("recalc-uncompressed", recalcOp, false)
	// }

	return nil
}

// Convert_v1alpha3_Bottle_To_v1alpha4_Bottle is a bottle converter
func Convert_v1alpha3_Bottle_To_v1alpha4_Bottle(in *v1alpha3.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	out.APIVersion = GroupVersion.String()
	out.Kind = "Bottle"

	// Migrate: Files -> Parts
	// New fields: Usage, expiration

	out.Catalog = in.Catalog
	out.Description = in.Description
	out.Usage = []Usage{}
	out.Keywords = in.Keywords
	out.Expiration = ""

	out.Sources = make([]Source, len(in.Sources))
	for i, s := range in.Sources {
		out.Sources[i] = Source(s)
	}

	out.Maintainers = make([]Maintainer, len(in.Maintainers))
	for i, mt := range in.Maintainers {
		out.Maintainers[i] = Maintainer(mt)
	}

	// Migrate: size -> layersize. digest -> layerdigest. usize -> size
	// Add fields: digest  (uncompressed data, will need recalculating)
	out.Parts = make([]Part, len(in.Files))
	for i, f := range in.Files {
		part := Part{
			Name:        f.Name,
			Size:        f.USize,
			LayerSize:   f.Size,
			Format:      f.Format,
			Digest:      DigestMap{},
			LayerDigest: DigestMap{Sha256: f.Digest.Sha256},
			Modified:    metav1.Time{Time: time.Now()},
			Labels:      f.Labels,
		}
		if part.Format == "" {
			part.Format = "raw" // oci.FormatFromMediaType(oci.RawMediaType)
		}
		out.Parts[i] = part
	}

	// If the controller object supports post migrate operations, this migration will use
	//  RecalcUncompressedSizes() error   to inform the controller that uncompressed digest and size information needs to
	//  be recalculated

	// Register uncompressed size/digest calculator post migrate operation.  This data was not included in the v1alpha2
	// if pmc, ok := m.Controller.(types.PostMigrateConsumer); ok {
	// 	pmc.RegisterPostMigrateOp("recalc-uncompressed", recalcOp, false)
	// }

	return nil
}

// recalcOp is a post migration operation that triggers a "RecalcUncompressedSizes()" function on the controller if it
// exists.
// func recalcOp(set any) error {
// 	type recalculator interface {
// 		RecalcUncompressedSizes() error
// 	}
// 	if r, ok := set.(recalculator); ok {
// 		err := r.RecalcUncompressedSizes()
// 		if os.IsNotExist(err) {
// 			// discard post migration errors where files aren't found, which occurs if a part selector is used
// 			return nil
// 		}
// 		return err
// 	}
// 	return nil
// }
