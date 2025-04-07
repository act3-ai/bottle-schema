// Package main defines a single command, genschema, to output the bottle schema
package main

import (
	"context"
	"os"

	"github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/jsonschema"

	"github.com/act3-ai/go-common/pkg/cmd"
	"github.com/act3-ai/go-common/pkg/runner"
)

func main() {
	root := cmd.NewGenschemaCmd(jsonschema.FS, []cmd.SchemaAssociation{})
	ctx := context.Background()
	root.SetContext(ctx)
	if err := runner.Run(root.Context(), root, "GENSCHEMA_VERBOSITY"); err != nil {
		os.Exit(1)
	}
}
