// Package main defines a single command, genschema, to output the bottle schema
package main

import (
	"os"

	"gitlab.com/act3-ai/asce/data/schema/pkg/apis/data.act3-ace.io/jsonschema"

	"gitlab.com/act3-ai/asce/go-common/pkg/cmd"
	"gitlab.com/act3-ai/asce/go-common/pkg/runner"
)

func main() {
	root := cmd.NewGenschemaCmd(jsonschema.FS, []cmd.SchemaAssociation{})

	if err := runner.Run(root, "GENSCHEMA_VERBOSITY"); err != nil {
		os.Exit(1)
	}
}
