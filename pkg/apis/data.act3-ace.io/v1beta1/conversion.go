package v1beta1

import (
	"github.com/opencontainers/go-digest"
	"k8s.io/apimachinery/pkg/conversion"

	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha2"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha3"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha4"
	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha5"
	"github.com/act3-ai/bottle-schema/pkg/mediatype"
)

// Convert_v1alpha5_Bottle_To_v1beta1_Bottle converts Bottle from v1alpha5 to v1beta1
func Convert_v1alpha5_Bottle_To_v1beta1_Bottle(in *v1alpha5.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	out.APIVersion = GroupVersion.String()
	out.Kind = "Bottle"

	// No root level migrations
	out.Annotations = in.Annotations
	out.Labels = in.Labels
	out.Description = in.Description

	// migrate sources, URL => URI
	out.Sources = make([]Source, len(in.Sources))
	for i, s := range in.Sources {
		out.Sources[i] = Source{
			Name: s.Name,
			URI:  s.URL,
		}
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

	// migrate public artifact to use mediaType instead of type
	out.PublicArtifacts = make([]PublicArtifact, len(in.PublicArtifacts))
	for i, artV5 := range in.PublicArtifacts {
		mediaType := mediatype.DetermineType(artV5.Path) // determine the mediaType of the public artifact
		out.PublicArtifacts[i] = PublicArtifact{
			MediaType: mediaType,
			Name:      artV5.Name,
			Path:      artV5.Path,
			Digest:    digest.Digest(artV5.Digest),
		}
	}

	// migrate parts, new part structure with layer info removed from bottle definition core schema
	out.Parts = make([]Part, len(in.Parts))
	for i, f := range in.Parts {
		out.Parts[i] = Part{
			Name:   f.Name,
			Size:   f.Size,
			Digest: digest.Digest(f.Digest),
			Labels: f.Labels,
		}
	}

	// a local interface for providing local part data stripped from the schema to the controller
	// type LocalPartDataConsumer interface {
	// 	SetLocalPartInfo(name string, layerSize int64, format string, layerDigest string, modified time.Time)
	// }

	// capture the local part data that was iniously stored in the config, by offering it to a part update consumer
	// in the controller
	// if pmc, ok := m.Controller.(types.PostMigrateConsumer); ok {
	// 	pmc.RegisterPostMigrateOp("record-local-part-info", func(ctl any) error {
	// 		auc, ok := ctl.(LocalPartDataConsumer)
	// 		if !ok {
	// 			return nil
	// 		}
	// 		for _, e := range in.Parts {
	// 			auc.SetLocalPartInfo(e.Name, e.LayerSize, e.Format, e.LayerDigest, e.Modified.Time)
	// 		}
	// 		return nil
	// 	}, false)
	// }

	return nil
}

// Convert_v1alpha4_Bottle_To_v1beta1_Bottle converts Bottle from v1alpha4 to v1beta1
func Convert_v1alpha4_Bottle_To_v1beta1_Bottle(in *v1alpha4.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a5 := &v1alpha5.Bottle{}
	if err := scope.Convert(in, a5); err != nil {
		return err
	}
	return scope.Convert(a5, out)
}

// Convert_v1alpha3_Bottle_To_v1beta1_Bottle converts Bottle from v1alpha3 to v1beta1
func Convert_v1alpha3_Bottle_To_v1beta1_Bottle(in *v1alpha3.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a4 := &v1alpha4.Bottle{}
	if err := scope.Convert(in, a4); err != nil {
		return err
	}
	return scope.Convert(a4, out)
}

// Convert_v1alpha2_Bottle_To_v1beta1_Bottle converts Bottle from v1alpha2 to v1beta1
func Convert_v1alpha2_Bottle_To_v1beta1_Bottle(in *v1alpha2.Bottle, out *Bottle, scope conversion.Scope) error { //revive:disable-line:var-naming
	a4 := &v1alpha4.Bottle{}
	if err := scope.Convert(in, a4); err != nil {
		return err
	}
	return scope.Convert(a4, out)
}
