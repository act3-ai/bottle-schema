package v1

import (
	"encoding/json"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	yamljson "sigs.k8s.io/yaml"
)

func testBottle() *Bottle {
	bottle := Bottle{}
	bottle.APIVersion = GroupVersion.String()
	bottle.Kind = "Bottle"
	bottle.Labels = map[string]string{
		"mykey": "myvalue",
		"a":     "b",
	}
	bottle.Description = "My bottle name\nMy cool bottle is so neat!"
	bottle.Sources = []Source{
		{
			Name: "Original data",
			URI:  "https://mydataset.example.com",
		},
	}
	bottle.PublicArtifacts = []PublicArtifact{
		{
			MediaType: "text/plain",
			Name:      "some file",
			Path:      "file.txt",
			Digest:    "sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0",
		},
	}
	bottle.Deprecates = []digest.Digest{
		digest.Digest("sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9"),
		digest.Digest("sha256:2dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9"),
	}
	bottle.Parts = []Part{
		{
			Name:   "file.txt",
			Size:   45,
			Digest: "sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0",
			Labels: map[string]string{
				"key": "value",
			},
		},
	}
	return &bottle
}

func TestBottle_ToDocumentedYAML(t *testing.T) {
	assert := assert.New(t)

	bottle := testBottle()

	out, err := bottle.ToDocumentedYAML()
	assert.NoError(err)

	expected :=
		`# ACE Data Bottle definition document containing the metadata

apiVersion: data.act3-ace.io/v1
kind: Bottle

# Labels are used to classify a bottle.  Selectors can later be used on these labels to select a subset of bottles.
# Follows Kubernetes conventions for labels.
labels:
    a: b
    mykey: myvalue

# Arbitrary user-defined content. Useful for storing non-standard metadata.
# Follows Kubernetes conventions for annotations.
annotations: {}
# key: "some value that is allowed to contain spaces and other character!"


# A human readable description of this Bottle.
# This field will be searched by researchers to discover this bottle.
description: |-
    My bottle name
    My cool bottle is so neat!

# Information about the bottle sources (where this bottle came from)
sources:
    - name: Original data
      uri: https://mydataset.example.com

# Contact information for bottle authors
authors: []
# - name: Your full name
#   email: someone@example.com
#   url: https://myhomepage.example.com # optional


# Contains metric data for a given experiment
metrics: []
# - name: log loss
#   description: natural log of the loss function
#   value: "45.2" # must be a numeric string (the quotes are required)


# Files intended to be exposed to the telemetry server for easy viewing
publicArtifacts:
    - name: some file
      path: file.txt
      mediaType: text/plain
      digest: sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0

# Bottle ID(s) to be deprecated by this bottle
deprecates:
    - sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9
    - sha256:2dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9

# Each bottle part may also have "part labels".
`

	assert.Equal(expected, string(out))
}

func TestBottle_YAML(t *testing.T) {
	assert := assert.New(t)
	bottle := testBottle()

	out, err := yamljson.Marshal(bottle)
	assert.NoError(err)

	// the keys are correct (match JSON) but the keys are sorted
	expected :=
		`apiVersion: data.act3-ace.io/v1
deprecates:
- sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9
- sha256:2dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9
description: |-
  My bottle name
  My cool bottle is so neat!
kind: Bottle
labels:
  a: b
  mykey: myvalue
parts:
- digest: sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0
  labels:
    key: value
  name: file.txt
  size: 45
publicArtifacts:
- digest: sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0
  mediaType: text/plain
  name: some file
  path: file.txt
sources:
- name: Original data
  uri: https://mydataset.example.com
`
	assert.Equal(expected, string(out))
}

func TestBottle_JSON(t *testing.T) {
	assert := assert.New(t)
	bottle := testBottle()

	out, err := json.Marshal(bottle)
	assert.NoError(err)

	expected := `{"kind":"Bottle","apiVersion":"data.act3-ace.io/v1","labels":{"a":"b","mykey":"myvalue"},"description":"My bottle name\nMy cool bottle is so neat!","sources":[{"name":"Original data","uri":"https://mydataset.example.com"}],"publicArtifacts":[{"name":"some file","path":"file.txt","mediaType":"text/plain","digest":"sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"}],"deprecates":["sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9","sha256:2dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c9"],"parts":[{"name":"file.txt","size":45,"digest":"sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0","labels":{"key":"value"}}]}`
	assert.Equal(expected, string(out))
}
