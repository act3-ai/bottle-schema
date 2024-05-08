// Package jsonschema is a single purpose package to store the jsonschema file for the data.act3-ace.io API
package jsonschema

import (
	"embed"
	"fmt"

	"github.com/invopop/jsonschema"
)

// FS stores the JSONSchema file as an embed.FS
//
//go:embed data.act3-ace.io.schema.json
var FS embed.FS

// Filename is the name of the embedded schema file
const Filename = "data.act3-ace.io.schema.json"

// Data returns the raw data of the JSONSchema file
func Data() ([]byte, error) {
	data, err := FS.ReadFile(Filename)
	if err != nil {
		return data, fmt.Errorf("could not read file %s: %w", Filename, err)
	}

	return data, nil
}

// Schema returns the schema as a invopo/jsonschema.Schema struct
func Schema() (*jsonschema.Schema, error) {
	data, err := Data()
	if err != nil {
		return nil, err
	}

	schema := &jsonschema.Schema{}
	err = schema.UnmarshalJSON(data)
	if err != nil {
		return schema, fmt.Errorf("could not unmarshal into schema struct: %w", err)
	}

	return schema, nil
}
