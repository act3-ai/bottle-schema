// Package main defines a single command, genschema, to output the bottle schema
package main

import (
	"os"

	"git.act3-ace.com/ace/data/schema/pkg/apis/data.act3-ace.io/jsonschema"
	"git.act3-ace.com/ace/go-common/pkg/cmd"
	"git.act3-ace.com/ace/go-common/pkg/runner"
)

func main() {
	root := cmd.NewGenschemaCmd(jsonschema.FS, []cmd.SchemaAssociation{})

	if err := runner.Run(root, "GENSCHEMA_VERBOSITY"); err != nil {
		os.Exit(1)
	}
}
