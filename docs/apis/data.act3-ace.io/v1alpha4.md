# API Reference

## Packages
- [bottle.data.act3-ace.io/v1alpha4](#bottledataact3-aceiov1alpha4)


## bottle.data.act3-ace.io/v1alpha4

Package v1alpha4 provides the Bottle types used ACE Data Bottles

### Resource Types
- [Bottle](#bottle)



#### Bottle



Bottle represents the overall structure of a data set entry.json or entry.yaml



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `bottle.data.act3-ace.io/v1alpha4`
| `kind` _string_ | `Bottle`
| `catalog` _boolean_ |  |
| `description` _string_ |  |
| `sources` _[Source](#source) array_ |  |
| `maintainers` _[Maintainer](#maintainer) array_ |  |
| `usage` _[Usage](#usage) array_ |  |
| `keywords` _string array_ |  |
| `expiration` _string_ |  |
| `parts` _[Part](#part) array_ |  |


#### DigestMap



DigestMap represents the map of digest values available in a dataset metadata record

_Appears in:_
- [Part](#part)

| Field | Description |
| --- | --- |
| `sha256` _string_ |  |


#### Maintainer



Maintainer is a collection of information about a maintainer, including name, email, and a URL link

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `email` _string_ |  |
| `url` _string_ |  |


#### Part



Part represents the layout of individual file records in a dataset metadata json file

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `size` _integer_ |  |
| `layerSize` _integer_ |  |
| `format` _string_ |  |
| `digest` _[DigestMap](#digestmap)_ |  |
| `layerDigest` _[DigestMap](#digestmap)_ |  |
| `modified` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ |  |
| `labels` _object (keys:string, values:string)_ |  |


#### Source



Source is a definition of a dataset source, containing a name and a url

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `url` _string_ |  |


#### Usage



Usage is a collection of information about usage documentation included in the bottle.  There can be multiple entries, each with a unique topic name, and referring to a file in the bottle by path.  The path provided can be within a directory that is archived (and thus does not correspond to a bottle part directly)

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `topic` _string_ |  |
| `name` _string_ |  |
| `file` _string_ |  |


