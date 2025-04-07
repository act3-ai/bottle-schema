# Conventions

BottleID is the digest of the bottle config (e.g., `sha256:beefefd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9`).  It often used to identify the bottle created.

BottleRef is a URI that references a bottle or a part of a bottle via the part selectors.  The optional part selector is in the fragment of the URI. The selectors are separated by "|".  BottleRefs come in many forms shown below:

- OCI repository with a tag (e.g., `registry.example.com/repo/name:v1`, `registry.example.com/repo/name:v1#partkey!=value1,mykey=value2|partkey2=45`)
- OCI repository with a manifest digest (e.g., `registry.example.com/repo/name@sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9#partkey!=value1,mykey=value2|partkey2=45`)
- OCI repository with a tag (ignored) and a manifest digest (e.g., `registry.example.com/repo/name:v1@sha256:05a8efd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9#partkey!=value1,mykey=value2|partkey2=45`)
- URI with `bottle` scheme  (e.g., URL encoded version of `bottle:sha256:beefefd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?selector=partkey!=value1,mykey=value2&selector=partkey2=45`)
- URI with `hash` scheme from [hash-uri](https://github.com/hash-uri/hash-uri) (e.g., URL encoded version of `hash://sha256/beefefd3483c60a4364d3f6f328ee1897facdbffb043b51941424a34121bbbe9?type=application/vnd.act3-ace.bottle.config.v1+json&selector=partkey!=value1,mykey=value2&selector=partkey2=45`)

All of the above BottleRefs could point to the same bottle.  The bottle may be stored in many different places and maybe be referenced by different digests (different algorithms).  It may also have different manifests because compression, encryption, embedded signatures, different compression level and different compression algorithms all could change manifest without changing the BottleID.

Note that the manifest ID (a.k.a., manifest digest) is not the same as the bottle ID (a.k.a., bottle digest).
