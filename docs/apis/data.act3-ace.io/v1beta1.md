# API Reference

## Packages
- [bottle.data.act3-ace.io/v1beta1](#bottledataact3-aceiov1beta1)


## bottle.data.act3-ace.io/v1beta1

Package v1beta1 provides the Bottle types used ACE Data Bottles

### Resource Types
- [Bottle](#bottle)



#### Author



Author is a collection of information about a author, including name, email, and a URL link

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `email` _string_ |  |
| `url` _string_ |  |


#### Bottle



Bottle represents the overall structure of a data set entry.json or entry.yaml



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `bottle.data.act3-ace.io/v1beta1`
| `kind` _string_ | `Bottle`
| `labels` _object (keys:string, values:string)_ |  |
| `annotations` _object (keys:string, values:string)_ |  |
| `description` _string_ |  |
| `sources` _[Source](#source) array_ |  |
| `authors` _[Author](#author) array_ |  |
| `metrics` _[Metric](#metric) array_ |  |
| `publicArtifacts` _[PublicArtifact](#publicartifact) array_ |  |
| `parts` _[Part](#part) array_ |  |


#### Metric



Metric is a collection of data about an experiment. Used to document results in metadata

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `description` _string_ |  |
| `value` _string_ |  |


#### Part



Part represents the layout of individual file records in a bottle metadata json file

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `size` _integer_ |  |
| `digest` _Digest_ |  |
| `labels` _object (keys:string, values:string)_ |  |


#### PublicArtifact



PublicArtifact is a collection of information about files included in the bottle that should be treated specially. There can be multiple entries, each with a type, and referring to a file in the bottle by path. The path provided can be within a directory that is archived (and thus does not correspond to a bottle part directly) These files will be exposed to the telemetry server/catalog explicitly, thus should not contain sensitive information

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `path` _string_ |  |
| `mediaType` _string_ |  |
| `digest` _Digest_ |  |


#### Source



Source is a definition of a data source, containing a name and a url

_Appears in:_
- [Bottle](#bottle)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `uri` _string_ |  |


