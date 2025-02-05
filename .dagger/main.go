// A generated module for Schema functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/schema/internal/dagger"
)

type Schema struct {
	// source code directory
	Source *dagger.Directory

	// +private
	RegistryConfig *dagger.RegistryConfig
}

func New(
	// top level source code directory
	// +defaultPath="/"
	src *dagger.Directory,
) *Schema {
	return &Schema{
		Source:         src,
		RegistryConfig: dag.RegistryConfig(),
	}
}

// Add credentials for a registry.
func (s *Schema) WithRegistryAuth(
	// registry's hostname
	address string,
	// username in registry
	username string,
	// password or token for registry
	secret *dagger.Secret,
) *Schema {
	s.RegistryConfig = s.RegistryConfig.WithRegistryAuth(address, username, secret)
	return s
}

// Removes credentials for a registry.
func (s *Schema) WithoutRegistryAuth(
	// registry's hostname
	address string,
) *Schema {
	s.RegistryConfig = s.RegistryConfig.WithoutRegistryAuth(address)
	return s
}

// Run all tests
func (s *Schema) Test(
	ctx context.Context,
) (string, error) {
	return dag.Go().
		WithSource(s.Source).
		WithCgoDisabled().
		Exec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}

// Build the genschema executable
func (s *Schema) Build(
	ctx context.Context,
) *dagger.File {
	return dag.Go().
		WithSource(s.Source).
		WithCgoDisabled().
		Build(dagger.GoWithSourceBuildOpts{
			Pkg: "./cmd/genschema",
		})
}

func (s *Schema) Release(
	// top level source code directory
	// +defaultPath="/"
	src *dagger.Directory,

	// gitlab token
	// +optional
	token *dagger.Secret,
) *Release {
	return &Release{
		Source: src,
		Token:  token,
	}
}
