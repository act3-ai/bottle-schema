package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RELEASED VERSION:  Do not modify the values in this file.
// SCHEMA UPDATES:
//  - None, initial version

// DigestMap represents the map of digest values available in a dataset metadata
// record
type DigestMap struct {
	Sha256 string `json:"sha256"`
}

// File represents the layout of individual file records in a dataset
// metadata json file
type File struct {
	Name     string            `json:"name"`
	Size     int64             `json:"size"`
	Format   string            `json:"format"`
	Digest   DigestMap         `json:"digest"`
	Modified metav1.Time       `json:"modified"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// Source is a definition of a dataset source, containing a name
// and a url
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Maintainer is a collection of information about a maintainer,
// including name, email, and a URL link
type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// +kubebuilder:object:root=true

// Bottle represents the overall structure of a data set entry.json
// or entry.yaml
type Bottle struct {
	metav1.TypeMeta `json:",inline"`

	Catalog     bool         `json:"catalog"`
	Description string       `json:"description"`
	Sources     []Source     `json:"sources"`
	Maintainers []Maintainer `json:"maintainers"`
	Keywords    []string     `json:"keywords"`
	Files       []File       `json:"files"`
}

// DocDefaults maintains default section documentation for the items in the bottle definition. These values can be
// used to display information about bottle metadata fields
var DocDefaults = map[string]string{
	"apiVersion":  "Bottle definition document containing metadata and a file list",
	"Catalog":     "\nA flag indicating whether or not to include this bottle in the bottle catalog",
	"description": "\nThe description field will be indexed and used by researchers to discover your bottle.",
	"sources":     "\nInformation about the bottle sources including names and urls.\nsources:\n  - name: MyDataSource\n    url: my.example.com",
	"maintainers": "\nContact information for bottle maintainers. \nmaintainers:\n  - name: My Name\n    email: my.email@example.com\n    url: my.example.url",
	"keywords":    "\nA list of keywords to associate with the bottle.\nkeywords:\n  - images\n  - faces",
	"files":       "\nA list of files associated with the bottle.\nThis information should be added automatically from files discovered in the bottle directory.\nYou may modify the labels dict.",
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
