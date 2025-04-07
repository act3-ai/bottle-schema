package v1beta1

import (
	"github.com/opencontainers/go-digest"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	util "github.com/act3-ai/bottle-schema/pkg/apis/internal/yaml"
)

// RELEASED VERSION:  Do make any schema breaking changes in this file.
// SCHEMA UPDATES:
//  - version update v1alpha5 to v1beta1
//  - remove local yaml-only values from part entry
//  - Changed Source URL field to URI
//

// AnnotationDeprecates is the key for the deprecation annotation.  The value is a comma separated list of bottle IDs (digests)
const AnnotationDeprecates = "bottle.data.act3-ace.io/deprecates"

// Part represents the layout of individual file records in a bottle
// metadata json file
type Part struct {
	Name   string            `json:"name,omitempty"`
	Size   int64             `json:"size,omitempty"`
	Digest digest.Digest     `json:"digest,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Source is a definition of a data source, containing a name
// and a url
type Source struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// Author is a collection of information about a author,
// including name, email, and a URL link
type Author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// PublicArtifact is a collection of information about files included in the bottle that should be treated specially.
// There can be multiple entries, each with a type, and referring to a file in the bottle by path.
// The path provided can be within a directory that is archived (and thus does not correspond to a bottle part directly)
// These files will be exposed to the telemetry server/catalog explicitly, thus should not contain sensitive information
type PublicArtifact struct { // Name change : publicArtifacts, publicFiles, extracted...
	Name      string        `json:"name,omitempty"`
	Path      string        `json:"path,omitempty"`
	MediaType string        `json:"mediaType,omitempty" yaml:"mediaType"` // yaml tag is needed because encoded field name (mediaType) has a capital letter.  This is needed for the HACK ToYamlNodes() to function properly.
	Digest    digest.Digest `json:"digest,omitempty"`
}

// Metric is a collection of data about an experiment. Used to document results in metadata
type Metric struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
}

// +kubebuilder:object:root=true

// Bottle represents the overall structure of a data set entry.json
// or entry.yaml
type Bottle struct {
	metav1.TypeMeta `json:",inline"`

	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	Description     string            `json:"description,omitempty"`
	Sources         []Source          `json:"sources,omitempty"`
	Authors         []Author          `json:"authors,omitempty"`
	Metrics         []Metric          `json:"metrics,omitempty"`
	PublicArtifacts []PublicArtifact  `json:"publicArtifacts,omitempty"`
	Parts           []Part            `json:"parts,omitempty"`
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
	commentTopFoot = `Each bottle part may also have "part labels".  Those can be added with "ace-dt bottle commit --label" or directly in the .labels.yaml files.`

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
  value: 45.2 # must be numeric`

	commentPublicArtifactsHead = "Files intended to be exposed to the telemetry server for easy viewing"
	commentPublicArtifactsFoot = `- name: name of artifact
  path: path/to/file/in/bottle
  mediaType: application/file-media-type # e.g., image/png
  digest: sha256:deedbeef # digest of file contents`
)
