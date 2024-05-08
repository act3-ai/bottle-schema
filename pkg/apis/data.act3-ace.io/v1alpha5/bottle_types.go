package v1alpha5

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RELEASED VERSION:  Do not modify the values in this file.
// SCHEMA UPDATES:
//  - version update v1alpha4 to v1alpha5
//  - Added labels(annotations), and arbitrary mapping of string to string
//  - Added metrics, used to document experiment results
//  - deleted catalog, catalog should be based off of exported files + registry location
//  - Changed maintainers to authors
//  - Changed 'usage' to PublicArtifacts, serves as a pointer to various important files, no longer just usage
//  - removed digestMap for a string

// Part represents the layout of individual file records in a bottle
// metadata json file
type Part struct {
	Name        string            `json:"name,omitempty"`
	Size        int64             `json:"size,omitempty"`
	LayerSize   int64             `json:"layerSize"`
	Format      string            `json:"format"`
	Digest      string            `json:"digest"`
	LayerDigest string            `json:"layerDigest"`
	Modified    metav1.Time       `json:"modified"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Source is a definition of a data source, containing a name
// and a url
type Source struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
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
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Path   string `json:"path,omitempty"`
	Digest string `json:"digest,omitempty"`
}

// Metric is a collection of data about an experiment. Used to document results in metadata
type Metric struct {
	Name        string `json:"name,omitempty"` // TODO how do metrics match? Name, unit, ...
	Description string `json:"description,omitempty"`
	Value       string `json:"value"`
}

// +kubebuilder:object:root=true

// Bottle represents the overall structure of a data set entry.json
// or entry.yaml
type Bottle struct {
	metav1.TypeMeta `json:",inline"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Description     string            `json:"description,omitempty"`
	Sources         []Source          `json:"sources"`
	Authors         []Author          `json:"authors"`
	Metrics         []Metric          `json:"metrics"`
	PublicArtifacts []PublicArtifact  `json:"publicArtifacts"`
	Parts           []Part            `json:"parts"`
}

// DocDefaults maintains default section documentation for the items in the bottle definition. These values can be
// used to display information about bottle metadata fields
var DocDefaults = map[string]string{
	"apiVersion":      "Bottle definition document containing metadata and a file list",
	"annotations":     "\nArbitrary user-defined content. String to string mapping. Use to hold arbitrary meta data.",
	"labels":          "\nFollows Kubernetes Label syntax:\n  63 characters or less (or empty).\n  Must begin and end with an alphanumeric character, unless empty\n  Contains dashes, underscores, dots, and alphanumerics between",
	"description":     "\nThe description field will be indexed and used by researchers to discover your bottle.",
	"sources":         "\nInformation about the bottle sources including names and urls.\nsources:\n  - name: MyDataSource\n    url: my.example.com",
	"authors":         "\nContact information for bottle authors. \nauthors:\n  - name: My Name\n    email: my.email@example.com\n    url: my.example.url",
	"metrics":         "\nContains metric data for a given experiment. \nmetrics:\n  - name: log loss\n    description: natural log of the loss function\n    value: 45.2",
	"publicArtifacts": "\nDocuments or files intended to be exposed to the telemetry server.\nThe name provides a label that can be used to refer to the file on a pull,\nand type is a single word description of what the document contains.\npublicArtifacts:\n  - type: Description\n    name: file-name\n    path: path/to/file\n    digest: sha256:digest",
	"parts":           "\nA list of parts associated with the bottle.\nThis information should be added automatically from files and directories discovered in the bottle directory.",
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
