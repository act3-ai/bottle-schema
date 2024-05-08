package mediatype

// MediaTypeBottle is the MediaType for the bottle as a whole
const MediaTypeBottle = "application/vnd.act3-ace.bottle"

const (
	// MediaTypeBottleConfig is the media type string for bottle configuration json
	MediaTypeBottleConfig = "application/vnd.act3-ace.bottle.config.v1+json"

	// MediaTypeLayerTarZstd is the media type string for tar+zstd layers
	MediaTypeLayerTarZstd = "application/vnd.act3-ace.bottle.layer.v1.tar+zstd"

	// MediaTypeLayerTarGzip is the media type string for tar+gzip layers
	MediaTypeLayerTarGzip = "application/vnd.act3-ace.bottle.layer.v1.tar+gzip"

	// MediaTypeLayerTar is the media type string for tar layers.  This is often used when a directory's tar archive does not compress well.
	MediaTypeLayerTar = "application/vnd.act3-ace.bottle.layer.v1.tar"

	// MediaTypeLayerZstd is the media type string for zstd compressed files
	MediaTypeLayerZstd = "application/vnd.act3-ace.bottle.layer.v1+zstd"

	// MediaTypeLayer is the media type string for general binary data (i.e., raw files).
	MediaTypeLayer = "application/vnd.act3-ace.bottle.layer.v1"
)

// Still in use but should not be used to create new bottles
const (
	// MediaTypeLayerTarZstdOld is the media type string for tar+zstd archives.  This is older format where the part name was in the archive.
	MediaTypeLayerTarZstdOld = "application/vnd.act3-ace.bottle.layer.v1+tar+zstd"

	// TarGzipMediaType is the media type string for tar+gzip archives.  This is older format where the part name was in the archive.
	MediaTypeLayerTarGzipOld = "application/vnd.act3-ace.bottle.layer.v1+tar+gzip"

	// TarMediaType is the media type string for tar archives.  This is older format where the part name was in the archive.
	MediaTypeLayerTarOld = "application/vnd.act3-ace.bottle.layer.v1+tar"

	// MediaTypeLayerRaw is a raw binary media type that (no tar or compression).  This is often used when a file cannot be compressed.
	MediaTypeLayerRawOld = "application/vnd.act3-ace.bottle.layer.v1+raw"
)

// Obsolete media types
// TODO should we drop support for these in v1.0
const (
	// MediaTypeBottleConfigLegacy is the media type string for the bottle configuration json
	MediaTypeBottleConfigLegacy = "application/vnd.act3-ace.dataset.config.v1+json"

	// MediaTypeLayerTarZstdLegacy is the media type string for tar+zstd archives
	MediaTypeLayerTarZstdLegacy = "application/vnd.act3-ace.dataset.layer.v1+tar+zstd"

	// MediaTypeLayerTarGzipLegacy is the media type string for tar+gzip archives
	MediaTypeLayerTarGzipLegacy = "application/vnd.act3-ace.dataset.layer.v1+tar+gzip"

	// MediaTypeLayerTarLegacy is the media type string for tar archives
	MediaTypeLayerTarLegacy = "application/vnd.act3-ace.dataset.layer.v1+tar"

	// MediaTypeLayerZstdLegacy is the media type string for zstd compressed files
	MediaTypeLayerZstdLegacy = "application/vnd.act3-ace.dataset.layer.v1+zstd"

	// MediaTypeLayerRawLegacy is a binary media type that indicates a compressed or noisy format that can't be compressed
	MediaTypeLayerRawLegacy = "application/vnd.act3-ace.dataset.layer.v1+raw"
)

// IsLayer returns true for valid media types for layers
func IsLayer(mediaType string) bool {
	switch mediaType {
	case MediaTypeLayerTarZstd, MediaTypeLayerTarGzip, MediaTypeLayerTar, MediaTypeLayerZstd, MediaTypeLayer:
		return true
	case MediaTypeLayerTarZstdOld, MediaTypeLayerTarGzipOld, MediaTypeLayerTarOld, MediaTypeLayerRawOld:
		return true
	case MediaTypeLayerTarZstdLegacy, MediaTypeLayerTarGzipLegacy, MediaTypeLayerTarLegacy, MediaTypeLayerZstdLegacy, MediaTypeLayerRawLegacy:
		return true
	default:
		return false
	}
}

// IsArchived returns true if the format is a known archive format
func IsArchived(mediaType string) bool {
	switch mediaType {
	case "":
		panic("media type must be non-empty")
	case MediaTypeLayerTarZstd, MediaTypeLayerTarGzip, MediaTypeLayerTar:
		return true
	case MediaTypeLayerTarZstdOld, MediaTypeLayerTarGzipOld, MediaTypeLayerTarOld:
		return true
	case MediaTypeLayerTarZstdLegacy, MediaTypeLayerTarGzipLegacy, MediaTypeLayerTarLegacy:
		return true
	default:
		return false
	}
}

// IsRaw returns true if the format indicates that no processing should be done on the file
func IsRaw(mediaType string) bool {
	switch mediaType {
	case "":
		panic("media type must be non-empty")
	case MediaTypeLayer, MediaTypeLayerRawOld, MediaTypeLayerRawLegacy:
		return true
	default:
		return false
	}
}

// IsCompressed returns true if the format is a known compressed format, or incompressible data
func IsCompressed(mediaType string) bool {
	switch mediaType {
	case "":
		panic("media type must be non-empty")
	case MediaTypeLayerTarZstd, MediaTypeLayerTarGzip, MediaTypeLayerZstd:
		return true
	case MediaTypeLayerTarZstdOld, MediaTypeLayerTarGzipOld:
		return true
	case MediaTypeLayerTarZstdLegacy, MediaTypeLayerTarGzipLegacy, MediaTypeLayerZstdLegacy:
		return true
	default:
		return false
	}
}

// IsBottleConfig returns true if the provided media type matches a known bottle config media type
func IsBottleConfig(mediaType string) bool {
	return mediaType == MediaTypeBottleConfig || mediaType == MediaTypeBottleConfigLegacy
}
