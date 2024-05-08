package bottle

import (
	"context"
	"testing"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go"
	ocispecv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	v1 "gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/v1"
	"gitlab.com/act3-ai/asce/data/schema/pkg/mediatype"
	val "gitlab.com/act3-ai/asce/data/schema/pkg/validation"
)

type ConversionTestSuite struct {
	suite.Suite
	scheme *runtime.Scheme
	codecs serializer.CodecFactory
}

func (suite *ConversionTestSuite) SetupSuite() {
	suite.scheme = runtime.NewScheme()
	suite.NoError(AddToScheme(suite.scheme))

	suite.codecs = serializer.NewCodecFactory(suite.scheme, serializer.EnableStrict)
}

func (suite *ConversionTestSuite) TestLoad_NoMigration() {
	jsonData := `
	{
		"apiVersion": "data.act3-ace.io/v1",
		"kind": "Bottle",
		"description": "This is a bottle.",
		"sources": [
			{
				"name": "Training",
				"uri": "bottle:sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"
			}
		],
		"authors": [
			{
				"name": "Jane Smith",
				"email": "jane.smith@example.com"
			},
			{
				"name": "Bob Dillon",
				"email": "bdill@example.com"
			}
		]
	}
`

	bottleOriginal, err := runtime.Decode(suite.codecs.UniversalDeserializer(), []byte(jsonData))
	suite.NoError(err)
	suite.T().Log(bottleOriginal)

	// NoOp conversion
	bottle := &v1.Bottle{}
	suite.NoError(suite.scheme.Convert(bottleOriginal, bottle, nil))
	// Defaulting adds back in the apiVersion and kind
	suite.scheme.Default(bottle)
	suite.Equal(v1.GroupVersion.WithKind("Bottle"), bottle.GroupVersionKind())

	// clear the bottle and try the old way
	bottle = &v1.Bottle{}
	// Regardless of if the bytes are of any external version,
	// it will be read successfully and converted into the internal version
	suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(jsonData), bottle))
	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("This is a bottle.", bottle.Description)
	suite.Len(bottle.Authors, 2)
	suite.NoError(bottle.Validate())
}

func (suite *ConversionTestSuite) TestLoad_MigrateWithManifest_v1beta1() {
	jsonData := `
	{
		"apiVersion": "data.act3-ace.io/v1beta1",
		"kind": "Bottle",
		"description": "This is a bottle.",
		"sources": [
			{
				"name": "Training",
				"uri": "bottle:sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"
			}
		],
		"authors": [
			{
				"name": "Jane Smith",
				"email": "jane.smith@example.com"
			},
			{
				"name": "Bob Dillon",
				"email": "bdill@example.com"
			}
		],
		"parts": [
			{
				"name": "dog",
				"digest": "sha256:9fdb955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0",
				"size": 150
			}
		]
	}
`

	manifest := &ocispecv1.Manifest{
		Versioned: ocispec.Versioned{SchemaVersion: 2},
		MediaType: ocispecv1.MediaTypeImageManifest,
		Config: ocispecv1.Descriptor{
			MediaType: mediatype.MediaTypeBottleConfig,
			Digest:    digest.FromString(jsonData),
			Size:      int64(len(jsonData)),
		},
		Layers: []ocispecv1.Descriptor{
			{
				MediaType: mediatype.MediaTypeLayerTarGzip,
				Digest:    digest.Digest("sha256:deedbeef282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"),
				Size:      50,
			},
		},
	}

	bottleOriginal, err := runtime.Decode(suite.codecs.UniversalDeserializer(), []byte(jsonData))
	suite.NoError(err)

	bottle := &v1.Bottle{}
	suite.NoError(suite.scheme.Convert(bottleOriginal, bottle, manifest))

	// Instead of doing a one-stop decode (desearialize, default, and convert) we do the above to allow us to inject a context
	// suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(jsonData), bottle))

	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("This is a bottle.", bottle.Description)
	suite.Len(bottle.Authors, 2)
	ctxManifest := val.ContextWithManifest(context.Background(), manifest)
	suite.NoError(bottle.ValidateWithContext(ctxManifest))
}

func (suite *ConversionTestSuite) TestLoad_Migrate_v1beta1() {
	jsonData := `
	{
		"apiVersion": "data.act3-ace.io/v1beta1",
		"kind": "Bottle",
		"description": "This is a bottle.",
		"sources": [
			{
				"name": "Training",
				"uri": "bottle:sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"
			}
		],
		"authors": [
			{
				"name": "Jane Smith",
				"email": "jane.smith@example.com"
			},
			{
				"name": "Bob Dillon",
				"email": "bdill@example.com"
			}
		]
	}
`

	bottleOriginal, err := runtime.Decode(suite.codecs.UniversalDeserializer(), []byte(jsonData))
	suite.NoError(err)

	bottle := &v1.Bottle{}
	suite.NoError(suite.scheme.Convert(bottleOriginal, bottle, nil))

	// Instead of doing a one-stop decode (desearialize, default, and convert) we do the above to allow us to inject a context
	// suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(jsonData), bottle))

	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("This is a bottle.", bottle.Description)
	suite.Len(bottle.Authors, 2)
	suite.NoError(bottle.Validate())
}

func (suite *ConversionTestSuite) TestLoad_Migrate_v1alpha4() {
	jsonData := `
	{
		"apiVersion": "data.act3-ace.io/v1alpha4",
		"kind": "Bottle",
		"description": "This is bottle 4 folks!\nIt is an older bottle (v1alpha4).",
		"sources": [
			{
				"name": "Training",
				"url": "bottle:sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"
			}
		],
		"maintainers": [
			{
				"name": "Jane Smith",
				"email": "jane.smith@example.com"
			},
			{
				"name": "Bob Dillon",
				"email": "bdill@example.com"
			}
		],
		"parts": [
			{
				"name": "foo/bar",
				"size": 45,
				"digest": {"sha256": "9a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae"},
				"labels": {
					"label1": "myfirstlabel",
					"label2": "mysecondlabel"
				}
			},
			{
				"name": "someusage",
				"size": 450,
				"digest": {"sha256": "3a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae"},
				"labels": {
					"label1": "myfirstlabel",
					"mykey": "something-else"
				}
			}
		]
	}
`
	bottle := &v1.Bottle{}
	suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(jsonData), bottle))
	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("This is bottle 4 folks!\nIt is an older bottle (v1alpha4).", bottle.Description)
	suite.Len(bottle.Authors, 2)
	suite.NoError(bottle.Validate())
}

func (suite *ConversionTestSuite) TestLoad_Migrate_v1alpha3() {
	jsonData := `
	{
		"apiVersion": "data.act3-ace.io/v1alpha3",
		"kind": "Bottle",
		"description": "This is a v1alpha3.",
		"sources": [
			{
				"name": "Training",
				"url": "bottle:sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"
			}
		],
		"maintainers": [
			{
				"name": "Jane Smith",
				"email": "jane.smith@example.com"
			},
			{
				"name": "Bob Dillon",
				"email": "bdill@example.com"
			}
		],
		"keywords": ["dog", "cat"],
		"files": [
			{
				"name": "foo/bar",
				"size": 45,
				"digest": {"sha256": "9a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae"},
				"labels": {
					"label1": "myfirstlabel",
					"label2": "mysecondlabel"
				}
			},
			{
				"name": "someusage",
				"size": 450,
				"digest": {"sha256": "3a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae"},
				"labels": {
					"label1": "myfirstlabel",
					"mykey": "something-else"
				}
			}
		]
	}
`
	bottle := &v1.Bottle{}
	suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(jsonData), bottle))
	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("This is a v1alpha3.", bottle.Description)
	suite.Len(bottle.Authors, 2)
	// v1beta1.Bottle.Parts.Digest is the content digest. In v1alpha3.Bottle.Files.Digest is the layer (compressed archive) digest.
	// Therefore we do not expect a digest and validation will fail until it is populated.
	suite.Equal("", bottle.Parts[0].Digest.String())
	// suite.NoError(bottle.Validate())
}

func (suite *ConversionTestSuite) TestLoad_Migrate_v1alpha2() {
	jsonData := `
	{
		"apiVersion": "data.act3-ace.io/v1alpha2",
		"kind": "Bottle",
		"description": "This is a v1alpha2.",
		"sources": [
			{
				"name": "Training",
				"url": "bottle:sha256:9dab955c282ecaacf81b1e1eda09300d42dfebf148583eef2b38ddd342da77c0"
			}
		],
		"maintainers": [
			{
				"name": "Jane Smith",
				"email": "jane.smith@example.com"
			},
			{
				"name": "Bob Dillon",
				"email": "bdill@example.com"
			}
		],
		"keywords": ["dog", "cat"],
		"files": [
			{
				"name": "foo/bar",
				"size": 45,
				"digest": {"sha256": "9a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae"},
				"labels": {
					"label1": "myfirstlabel",
					"label2": "mysecondlabel"
				}
			},
			{
				"name": "someusage",
				"size": 450,
				"digest": {"sha256": "3a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae"},
				"labels": {
					"label1": "myfirstlabel",
					"mykey": "something-else"
				}
			}
		]
	}
`
	bottle := &v1.Bottle{}
	suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(jsonData), bottle))
	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("This is a v1alpha2.", bottle.Description)
	suite.Len(bottle.Authors, 2)
	// v1beta1.Bottle.Parts.Digest is the content digest. In v1alpha3.Bottle.Files.Digest is the layer (compressed archive) digest.
	// Therefore we do not expect a digest and validation will fail until it is populated.
	suite.Equal("", bottle.Parts[0].Digest.String())
	// suite.NoError(bottle.Validate())
}

func (suite *ConversionTestSuite) TestLoad_WithYAML() {
	yamlData := `
apiVersion: data.act3-ace.io/v1
kind: Bottle
labels:
  type: testing
  group: testset
  epoch: "13"
  learning-rate: "0.001"
  refname: bottle1
annotations:
  viewer.data.act3-ace.io/Jupyter-Base: '{"accept":"application/x.jupyter.notebook+json, */*;q=0.8","acehub":{"image":"docker.io/jupyter/base-notebook","jupyter":true,"proxyType":"straight","resources":{"cpu":"1","memory":"1Gi"}}}'
description: |-
  MNIST Dataset
  Next Line
sources:
  - name: Data page
    uri: http://data.example.com
  - name: Hash type reference - NOT KNOWN to this telemetry server
    uri: hash://sha256/42a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?type=application/vnd.act3-ace.bottle.config.v1%2Bjson#partkey!=value1,mykey=value2|partkey2=45
authors:
  - name: John Smith
    url: https://john.example.com
    email: john.smith@example.com
metrics:
  - name: training loss
    value: "3.141592654"
  - name: AUC
    value: "0.985"
    description: Area under the curve
publicArtifacts:
  - name: Some text
    mediaType: text/plain
    path: sample.txt
    digest: sha256:eab4fe92c4c81e25676d91b3dac3191fe3d0a22e2a6644b76726a7683862a339
parts:
  - name: foo/bar
    size: 45
    digest: sha256:0b1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae
    labels:
      label1: myfirstlabel
      label2: mysecondlabel
  - name: sample.txt
    size: 450
    digest: sha256:0a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae
    labels:
      label1: myfirstlabel
      label2: otherlabel
`

	bottle := &v1.Bottle{}
	suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(yamlData), bottle))
	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("MNIST Dataset\nNext Line", bottle.Description)
	suite.Len(bottle.Authors, 1)
	suite.Len(bottle.Metrics, 2)
	suite.Equal(bottle.Metrics[0].Name, "training loss")
	suite.NoError(bottle.Validate())
}

func (suite *ConversionTestSuite) TestLoad_Migrate_V1Alpha5() {
	yamlData := `
apiVersion: data.act3-ace.io/v1alpha5
kind: Bottle
labels:
  type: testing
  group: testset
  epoch: "13"
  learning-rate: "0.001"
  refname: bottle1
annotations:
  viewer.data.act3-ace.io/Jupyter-Base: '{"accept":"application/x.jupyter.notebook+json, */*;q=0.8","acehub":{"image":"docker.io/jupyter/base-notebook","jupyter":true,"proxyType":"straight","resources":{"cpu":"1","memory":"1Gi"}}}'
description: |-
  MNIST Dataset
  Next Line
sources:
  - name: Data page
    url: http://data.example.com
  - name: Hash type reference - NOT KNOWN to this telemetry server
    url: hash://sha256/42a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?type=application/vnd.act3-ace.bottle.config.v1%2Bjson#partkey!=value1,mykey=value2|partkey2=45
authors:
  - name: John Smith
    url: https://john.example.com
    email: john.smith@example.com
metrics:
  - name: training loss
    value: "3.141592654"
  - name: AUC
    value: "0.985"
    description: Area under the curve
publicArtifacts:
  - name: Some text
    type: text/plain
    path: sample.txt
    digest: sha256:eab4fe92c4c81e25676d91b3dac3191fe3d0a22e2a6644b76726a7683862a339
parts:
  - name: foo/bar
    size: 45
    digest: sha256:0b1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae
    labels:
      label1: myfirstlabel
      label2: mysecondlabel
  - name: sample.txt
    size: 450
    digest: sha256:0a1de4364cfd94d75e7bda5d0583bcb136d6437c88a36dc06bcd64566a3530ae
    labels:
      label1: myfirstlabel
      label2: otherlabel
`

	bottle := &v1.Bottle{}
	suite.NoError(runtime.DecodeInto(suite.codecs.UniversalDecoder(), []byte(yamlData), bottle))
	suite.Equal(v1.GroupVersion.Version, bottle.GroupVersionKind().Version)
	suite.Equal("Bottle", bottle.GroupVersionKind().Kind)
	suite.Equal("MNIST Dataset\nNext Line", bottle.Description)
	suite.Len(bottle.Authors, 1)
	suite.Len(bottle.Metrics, 2)
	suite.Equal(bottle.Metrics[0].Name, "training loss")
	suite.NoError(bottle.Validate())
}

func TestConversionTestSuite(t *testing.T) {
	suite.Run(t, new(ConversionTestSuite))
}
