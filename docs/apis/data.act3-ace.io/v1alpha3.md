# API Reference

## Packages
- [bottle.data.act3-ace.io/v1alpha3](#bottledataact3-aceiov1alpha3)


## bottle.data.act3-ace.io/v1alpha3

Package v1alpha3 provides the Bottle types used ACE Data Bottles

### Resource Types
- [Bottle](#bottle)



#### Bottle



Bottle represents the overall structure of a data set entry.json or entry.yaml



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `bottle.data.act3-ace.io/v1alpha3`
| `kind` _string_ | `Bottle`
| `catalog` _boolean_ |  |
| `description` _string_ |  |
| `sources` _[Source](#source) array_ |  |
| `maintainers` _[Maintainer](#maintainer) array_ |  |
| `keywords` _string array_ |  |
| `files` _[File](#file) array_ |  |


#### DigestMap



DigestMap represents the map of digest values available in a dataset metadata record

_Appears in:_
- [File](#file)

| Field | Description |
| --- | --- |
| `sha256` _string_ |  |


#### File



File represents the layout of individual file records in a dataset metadata json file

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `size` _integer_ |  |
| `usize` _integer_ |  |
| `format` _string_ |  |
| `digest` _[DigestMap](#digestmap)_ |  |
| `modified` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ |  |
| `labels` _object (keys:string, values:string)_ |  |


#### Maintainer



Maintainer is a collection of information about a maintainer, including name, email, and a URL link

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `email` _string_ |  |
| `url` _string_ |  |


#### Source



Source is a definition of a dataset source, containing a name and a url

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `url` _string_ |  |


