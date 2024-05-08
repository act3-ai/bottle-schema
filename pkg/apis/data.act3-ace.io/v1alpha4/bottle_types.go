package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RELEASED VERSION:  Do not modify the values in this file.
// SCHEMA UPDATES:
//  - version update v1alpha3 to v1alpha4
//  - Rename File to Part (also change Bottle.Files to Bottle.Parts)
//  - Remove Uncompressed size (USize) from Part config data (Size becomes uncompressed size)
//  - Recontextualize Digest to Content Digest (uncompressed digest)
//  - Remove Format from config data, this data is tracked using the media type in the manifest
//  - Add LayerSize to Part structure, not tracked in json, refers to blob size in manifest
//  - Add LayerDigest to Part structure, not tracked in json, refers to blob digest in manifest
//  - Add Usage for recording usage file paths and topics, removes usage flag from parts
//
// MIGRATION NOTES:
//  dataset.RecalcCompressedSizes is used to calculate both the uncompressed size and content digest
//  when retrieving v1alpha2 and v1alpha3 schema versions.  For performance, v1alpha2 to v1alpha3
//  migration is removed (to avoid performing the required tar operation more than once)

// DigestMap represents the map of digest values available in a dataset metadata
// record
type DigestMap struct {
	Sha256 string `json:"sha256"`
}

// Part represents the layout of individual file records in a dataset
// metadata json file
type Part struct {
	Name        string            `json:"name"`
	Size        int64             `json:"size"`
	LayerSize   int64             `json:"layerSize"`
	Format      string            `json:"format"`
	Digest      DigestMap         `json:"digest"`
	LayerDigest DigestMap         `json:"layerDigest"`
	Modified    metav1.Time       `json:"modified"`
	Labels      map[string]string `json:"labels,omitempty"`
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

// Usage is a collection of information about usage documentation included in the bottle.  There can be multiple
// entries, each with a unique topic name, and referring to a file in the bottle by path.  The path provided can be
// within a directory that is archived (and thus does not correspond to a bottle part directly)
type Usage struct {
	Topic string `json:"topic"`
	Name  string `json:"name"`
	File  string `json:"file"`
}

// +kubebuilder:object:root=true

// Bottle represents the overall structure of a data set entry.json
// or entry.yaml
type Bottle struct {
	metav1.TypeMeta `json:",inline"`
	Catalog         bool         `json:"catalog"`
	Description     string       `json:"description"`
	Sources         []Source     `json:"sources"`
	Maintainers     []Maintainer `json:"maintainers"`
	Usage           []Usage      `json:"usage"`
	Keywords        []string     `json:"keywords"`
	Expiration      string       `json:"expiration"`
	Parts           []Part       `json:"parts"`
}

// DocDefaults maintains default section documentation for the items in the bottle definition. These values can be
// used to display information about bottle metadata fields
var DocDefaults = map[string]string{
	"apiVersion":  "Bottle definition document containing metadata and a file list",
	"catalog":     "\nA flag indicating whether or not to include this bottle in the bottle catalog",
	"description": "\nThe description field will be indexed and used by researchers to discover your bottle.",
	"sources":     "\nInformation about the bottle sources including names and urls.\nsources:\n  - name: MyDataSource\n    url: my.example.com",
	"maintainers": "\nContact information for bottle maintainers. \nmaintainers:\n  - name: My Name\n    email: my.email@example.com\n    url: my.example.url",
	"usage":       "\nDocuments or files that provide usage instructions or information for the data.  The name provides a tag or name that can be used to refer to the usage document on a pull, and topic is a short description of what the document covers.\nusage:\n  - topic: Usage Topic Description\n    name: usage-name\n    file: path/to/usage-file",
	"keywords":    "\nA list of keywords to associate with the bottle.\nkeywords:\n  - images\n  - faces",
	"expiration":  "\nA date and time when the bottle data should be considered expired, in YYYY-MM-DD HH:MM:SS format (UTC).  Leave empty for no expiration.",
	"parts":       "\nA list of parts associated with the bottle.\nThis information should be added automatically from files and directories discovered in the bottle directory.",
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
