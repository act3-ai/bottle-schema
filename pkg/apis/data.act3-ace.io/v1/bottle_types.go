package v1

import (
	"github.com/opencontainers/go-digest"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	util "git.act3-ace.com/ace/data/schema/pkg/apis/internal/yaml"
)

// RELEASED VERSION:  Do make any schema breaking changes in this file.
// SCHEMA UPDATES:
//  - version update v1beta1 to v1
//  - require Part.Name to have a trailing slash for directory parts (no slash for file parts)
//  - add deprecates as a field

// Part represents the layout of individual file records in a bottle
// metadata json file
type Part struct {
	// Name is the path to the part in the bottle.
	// File parts have no trailing slash.
	// Directory parts have a trailing slash.
	Name string `json:"name,omitempty"`

	// Size is the number of bytes in the raw/uncompressed part.
	// For files this is simply the size of the original file.
	// For directories this is the size of the archive.
	Size int64 `json:"size,omitempty"`

	// Digest is the content digest.
	// For files this is the digest of the file.
	// For directories this is the digest of the archive.
	Digest digest.Digest `json:"digest,omitempty"`

	// Labels to apply to the part (useful for use with part selectors to refer to partial bottles).
	Labels map[string]string `json:"labels,omitempty"`
}

// Source is a definition of a data source used to track data lineage.
// A source is another URI (e.g., website, bottle) that this bottle was derived from.
// For example a bottle containing a ML model should include a source for the training set.
type Source struct {
	// Name is the human understandable name of the source
	Name string `json:"name,omitempty"`

	// URI points to the source.
	// TODO document all the ways we support (the docs are in telemetry/conventions.md right now and need to move over here).
	URI string `json:"uri,omitempty"`
}

// Author is a collection of information about a author.
type Author struct {
	// Name of the author.
	Name string `json:"name,omitempty"`

	// Email of the author.
	Email string `json:"email,omitempty"`

	// URL of the author's homepage.
	URL string `json:"url,omitempty"`
}

// PublicArtifact is a collection of information about files included in the bottle that should be treated specially.
// The path provided can be within a directory that is archived (and thus does not correspond to a bottle part directly).
// These files will be exposed to the telemetry server/catalog explicitly, thus should not contain sensitive information.
// Often artifacts are figures of merit or key evaluation/performance results outlining what is in the bottle.
// They must be relatively small (< 1MiB) in size for compatibility with the Telemetry server.
// Public artifacts are just files.  They are not allowed to be directories.
type PublicArtifact struct {
	// Name is the human understandable name of the artifact.
	Name string `json:"name,omitempty"`

	// Path is the path to the file in this bottle (this can drill down into a directory part).
	Path string `json:"path,omitempty"`

	// MediaType is the an RFC 2045 compliant media type for use in determining how to display this artifact.
	// For ipynb files use "application/x.jupyter.notebook+json".
	MediaType string `json:"mediaType,omitempty" yaml:"mediaType"` // yaml tag is needed because encoded field name (mediaType) has a capital letter.  This is needed for the HACK ToYamlNodes() to function properly.

	// Digest of the file.
	Digest digest.Digest `json:"digest,omitempty"`
}

// Metric is a collection of data about an experiment. Used to document quantifiable results in metadata
type Metric struct {
	// Name is the name for this metric.
	// Try to be consistent in naming of metrics.
	Name string `json:"name,omitempty"`

	// Description is the detailed description of what this metric represents.
	Description string `json:"description,omitempty"`

	// Value is the floating point value (stored as a string) for this metric.
	Value string `json:"value,omitempty"`
}

// +kubebuilder:object:root=true

// Bottle represents the overall structure of a data set entry.json
// or entry.yaml
type Bottle struct {
	metav1.TypeMeta `json:",inline"`

	// Labels are the bottle.
	// The allowable grammar for the keys and values matches kubernetes.
	// These are bottle wide labels (not to be confused with labels on individual parts).
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are the bottle.
	// The allowable grammar for the keys and values matches kubernetes.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Description is a detailed description of the bottle contents.
	Description string `json:"description,omitempty"`

	// Sources is the list of sources.
	Sources []Source `json:"sources,omitempty"`

	// Authors is the list of authors.
	Authors []Author `json:"authors,omitempty"`

	// Metrics is the list of metrics.
	Metrics []Metric `json:"metrics,omitempty"`

	// PublicArtifacts is the list of artifacts.
	PublicArtifacts []PublicArtifact `json:"publicArtifacts,omitempty"`

	// Deprecates is an array of bottle IDs that this bottle deprecates (a.k.a. supersedes).
	// Deprecated bottles should not be used for new work.
	// The deprecating bottle often fixes a typo or some other mistake in the deprecated bottle.
	Deprecates []digest.Digest `json:"deprecates,omitempty"`

	// Parts is a list of parts (the actual data of the bottle is referred to in the parts).
	Parts []Part `json:"parts,omitempty"`
}

// ToDocumentedYAML converts the bottle into YAML with comments explaining each field.  It omits the parts field.
func (b Bottle) ToDocumentedYAML() ([]byte, error) {
	// create a top level yaml.Node.  Note, the document node is already created by
	//  this point, so the top level is a key value mapping node.
	nodes := []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "apiVersion"},
		{Kind: yaml.ScalarNode, Value: b.APIVersion},

		{Kind: yaml.ScalarNode, Value: "kind"},
		{Kind: yaml.ScalarNode, Value: b.Kind},
	}

	addField := func(name, header, footer string, subNodes []*yaml.Node, empty bool) {
		if !empty {
			footer = ""
		}
		nodes = append(nodes, &yaml.Node{
			Kind:        yaml.ScalarNode,
			Value:       name,
			HeadComment: "\n" + header,
			FootComment: footer,
		})
		nodes = append(nodes, subNodes...)
	}

	// Labels
	subNodes, err := util.ToYamlNodes(b.Labels)
	if err != nil {
		return nil, err
	}
	addField("labels", commentLabelsHead, commentLabelsFoot, subNodes, len(b.Labels) == 0)

	// Annotations
	subNodes, err = util.ToYamlNodes(b.Annotations)
	if err != nil {
		return nil, err
	}
	addField("annotations", commentAnnotationsHead, commentAnnotationsFoot, subNodes, len(b.Annotations) == 0)

	// Description
	nodes = append(nodes,
		&yaml.Node{
			Kind:        yaml.ScalarNode,
			Value:       "description",
			HeadComment: "\n" + commentDescription,
		},
		&yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.LiteralStyle,
			Value: b.Description,
		},
	)

	// Sources
	subNodes, err = util.ToYamlNodes(b.Sources)
	if err != nil {
		return nil, err
	}
	addField("sources", commentSourcesHead, commentSourcesFoot, subNodes, len(b.Sources) == 0)

	// Authors
	subNodes, err = util.ToYamlNodes(b.Authors)
	if err != nil {
		return nil, err
	}
	addField("authors", commentAuthorsHead, commentAuthorsFoot, subNodes, len(b.Authors) == 0)

	// Metrics
	subNodes, err = util.ToYamlNodes(b.Metrics)
	if err != nil {
		return nil, err
	}
	addField("metrics", commentMetricsHead, commentMetricsFoot, subNodes, len(b.Metrics) == 0)

	// Public Artifacts
	subNodes, err = util.ToYamlNodes(b.PublicArtifacts)
	if err != nil {
		return nil, err
	}
	addField("publicArtifacts", commentPublicArtifactsHead, commentPublicArtifactsFoot, subNodes, len(b.PublicArtifacts) == 0)

	// Deprecates
	subNodes, err = util.ToYamlNodes(b.Deprecates)
	if err != nil {
		return nil, err
	}
	addField("deprecates", commentDeprecatesHead, commentDeprecatesFoot, subNodes, len(b.Deprecates) == 0)

	// We do not output Parts in this documented YAML view

	doc := &yaml.Node{
		Kind:        yaml.DocumentNode,
		HeadComment: commentTopHead,
		FootComment: commentTopFoot,
		Content: []*yaml.Node{
			{
				Kind:    yaml.MappingNode,
				Content: nodes,
			},
		},
	}

	return yaml.Marshal(doc)
}

// NewBottle returns a definition containing data initialized to default values.
func NewBottle() Bottle {
	return Bottle{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       "Bottle",
		},
	}
}

// DocDefaults maintains default section documentation for the items in the bottle definition. These values can be
// used to display information about bottle metadata fields
const (
	commentTopHead = "ACE Data Bottle definition document containing the metadata"
	commentTopFoot = `Each bottle part may also have "part labels".`

	commentLabelsHead = `Labels are used to classify a bottle.  Selectors can later be used on these labels to select a subset of bottles.
Follows Kubernetes conventions for labels.`
	commentLabelsFoot = "key: value"

	commentAnnotationsHead = `Arbitrary user-defined content. Useful for storing non-standard metadata.
Follows Kubernetes conventions for annotations.`
	commentAnnotationsFoot = `key: "some value that is allowed to contain spaces and other character!"`

	commentDescription = `A human readable description of this Bottle.
This field will be searched by researchers to discover this bottle.`

	commentSourcesHead = "Information about the bottle sources (where this bottle came from)"
	commentSourcesFoot = `- name: Name of source
  uri: https://my-source.example.com
- name: Bottle reference name
  uri: bottle://sha256:deedbeef`

	commentAuthorsHead = "Contact information for bottle authors"
	commentAuthorsFoot = `- name: Your full name
  email: someone@example.com
  url: https://myhomepage.example.com # optional`

	commentMetricsHead = "Contains metric data for a given experiment"
	commentMetricsFoot = `- name: log loss
  description: natural log of the loss function
  value: "45.2" # must be a numeric string (the quotes are required)`

	commentPublicArtifactsHead = "Files intended to be exposed to the telemetry server for easy viewing"
	commentPublicArtifactsFoot = `- name: name of artifact
  path: path/to/file/in/bottle
  mediaType: application/file-media-type # e.g., image/png
  digest: sha256:deedbeef # digest of file contents`

	commentDeprecatesHead = "Bottle ID(s) to be deprecated by this bottle"
	commentDeprecatesFoot = `- sha256:deedbeef # bottle ID`
)
