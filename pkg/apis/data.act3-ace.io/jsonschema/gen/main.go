// Package main is a fake package for generating code.
package main

import (
	"fmt"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/runtime"

	bottle "git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io"
	"git.act3-ace.com/ace/go-common/pkg/genschema"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify a target directory for schema generation.")
	}

	dir := os.Args[1]

	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatal(fmt.Errorf("failed to create schema directory: %w", err))
	}

	scheme := runtime.NewScheme()
	if err := bottle.AddToScheme(scheme); err != nil {
		log.Fatal(fmt.Errorf("error adding type data to conversion scheme: %w", err))
	}

	if err := genschema.GenerateGroupSchemas(dir, scheme, []string{"data.act3-ace.io"}, "git.act3-ace.com/ace/data/schema"); err != nil {
		log.Fatal(fmt.Errorf("error generating schema: %w", err))
	}
}
