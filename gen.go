// Package gen is a fake top level package to house generate commands
package gen

//go:generate go run pkg/apis/data.act3-ace.io/jsonschema/gen/main.go pkg/apis/data.act3-ace.io/jsonschema
//go:generate tool/controller-gen object paths=./...
